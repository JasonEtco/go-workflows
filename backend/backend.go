package backend

import (
	"context"

	"github.com/cschleiden/go-workflows/core"
	"github.com/cschleiden/go-workflows/core/task"
	"github.com/cschleiden/go-workflows/internal/history"
)

//go:generate mockery --name=Backend --inpackage
type Backend interface {
	// CreateWorkflowInstance creates a new workflow instance
	CreateWorkflowInstance(ctx context.Context, event core.WorkflowEvent) error

	// CancelWorkflowInstance cancels a running workflow instance
	CancelWorkflowInstance(ctx context.Context, instance core.WorkflowInstance) error

	// SignalWorkflow signals a running workflow instance
	SignalWorkflow(ctx context.Context, instanceID string, event history.Event) error

	// GetWorkflowInstance returns a pending workflow task or nil if there are no pending worflow executions
	GetWorkflowTask(ctx context.Context) (*task.Workflow, error)

	// ExtendWorkflowTask extends the lock of a workflow task
	ExtendWorkflowTask(ctx context.Context, instance core.WorkflowInstance) error

	// CompleteWorkflowTask checkpoints a workflow task retrieved using GetWorkflowTask
	//
	// This checkpoints the execution. events are new events from the last workflow execution
	// which will be added to the workflow instance history. workflowEvents are new events for the
	// completed or other workflow instances.
	CompleteWorkflowTask(ctx context.Context, instance core.WorkflowInstance, executedEvents []history.Event, workflowEvents []core.WorkflowEvent) error

	// GetActivityTask returns a pending activity task or nil if there are no pending activities
	GetActivityTask(ctx context.Context) (*task.Activity, error)

	// CompleteActivityTask completes a activity task retrieved using GetActivityTask
	CompleteActivityTask(ctx context.Context, instance core.WorkflowInstance, activityID string, event history.Event) error

	// ExtendActivityTask extends the lock of an activity task
	ExtendActivityTask(ctx context.Context, activityID string) error
}