package redis

import (
	"context"
	"fmt"

	"github.com/cschleiden/go-workflows/internal/history"
	"github.com/cschleiden/go-workflows/internal/tracing"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (rb *redisBackend) SignalWorkflow(ctx context.Context, instanceID string, event history.Event) error {
	instanceState, err := readInstance(ctx, rb.rdb, instanceID)
	if err != nil {
		return err
	}

	ctx = tracing.UnmarshalSpan(ctx, instanceState.Metadata)
	a := event.Attributes.(*history.SignalReceivedAttributes)
	_, span := rb.Tracer().Start(ctx, fmt.Sprintf("SignalWorkflow: %s", a.Name), trace.WithAttributes(
		attribute.String(tracing.WorkflowInstanceID, instanceID),
		attribute.String("signal.name", event.Attributes.(*history.SignalReceivedAttributes).Name),
	))
	defer span.End()

	if _, err = rb.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := addEventToStreamP(ctx, p, pendingEventsKey(instanceID), &event); err != nil {
			return fmt.Errorf("adding event to stream: %w", err)
		}

		if err := rb.workflowQueue.Enqueue(ctx, p, instanceID, nil); err != nil {
			if err != errTaskAlreadyInQueue {
				return fmt.Errorf("queueing workflow task: %w", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
