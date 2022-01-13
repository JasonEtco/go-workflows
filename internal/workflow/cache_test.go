package workflow

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/cschleiden/go-dt/internal/payload"
	"github.com/cschleiden/go-dt/pkg/core"
	"github.com/cschleiden/go-dt/pkg/core/task"
	"github.com/cschleiden/go-dt/pkg/history"
	"github.com/stretchr/testify/require"
)

func Test_Cache_StoreAndGet(t *testing.T) {
	c := NewWorkflowExecutorCache(context.Background(), DefaultWorkflowExecutorCacheOptions)

	i := core.NewWorkflowInstance("instanceID", "executionID")
	task := &task.Workflow{
		WorkflowInstance: i,
		History: []history.Event{
			history.NewHistoryEvent(
				history.EventType_WorkflowExecutionStarted,
				-1,
				&history.ExecutionStartedAttributes{
					Name:   "WorkflowWithActivity",
					Inputs: []payload.Payload{},
				},
			),
		},
	}

	e := NewExecutor(NewRegistry(), task)

	err := c.Store(context.Background(), i, e)
	require.NoError(t, err)

	e2, ok, err := c.Get(context.Background(), i)
	require.NoError(t, err)
	require.True(t, ok)

	require.Equal(t, e, e2)
}

func Test_Cache_Evic(t *testing.T) {
	c := NewWorkflowExecutorCache(context.Background(), WorkflowExecutorCacheOptions{
		CacheDuration: 1, // Should evict immediately
	})

	i := core.NewWorkflowInstance("instanceID", "executionID")
	task := &task.Workflow{
		WorkflowInstance: i,
		History: []history.Event{
			history.NewHistoryEvent(
				history.EventType_WorkflowExecutionStarted,
				-1,
				&history.ExecutionStartedAttributes{
					Name:   "WorkflowWithActivity",
					Inputs: []payload.Payload{},
				},
			),
		},
	}

	e := NewExecutor(NewRegistry(), task)

	err := c.Store(context.Background(), i, e)
	require.NoError(t, err)

	runtime.Gosched()
	time.Sleep(1 * time.Millisecond)

	e2, ok, err := c.Get(context.Background(), i)
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, e2)
}
