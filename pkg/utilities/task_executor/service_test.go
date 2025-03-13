package task_executor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para Tasker
type MockTask struct {
	mock.Mock
}

func (m *MockTask) Execute(ctx context.Context) (interface{}, int, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Int(1), args.Error(2)
}

func TestWorkerPool_Success(t *testing.T) {
	ctx := context.Background()

	task1 := new(MockTask)
	task1.On("Execute", ctx).Return("result1", 100, nil)

	task2 := new(MockTask)
	task2.On("Execute", ctx).Return("result2", 200, nil)

	tasks := map[string]Tasker{
		"task1": task1,
		"task2": task2,
	}

	results := WorkerPool(ctx, tasks, 2)

	assert.Len(t, results, 2)
	assert.Equal(t, "result1", results["task1"].Res)
	assert.Equal(t, 100, results["task1"].Time)
	assert.NoError(t, results["task1"].Err)

	assert.Equal(t, "result2", results["task2"].Res)
	assert.Equal(t, 200, results["task2"].Time)
	assert.NoError(t, results["task2"].Err)

	task1.AssertExpectations(t)
	task2.AssertExpectations(t)
}

func TestWorkerPool_WithErrors(t *testing.T) {
	ctx := context.Background()

	task1 := new(MockTask)
	task1.On("Execute", ctx).Return(nil, 50, errors.New("error in task1"))

	task2 := new(MockTask)
	task2.On("Execute", ctx).Return("result2", 100, nil)

	tasks := map[string]Tasker{
		"task1": task1,
		"task2": task2,
	}

	results := WorkerPool(ctx, tasks, 2)

	assert.Len(t, results, 2)

	assert.Nil(t, results["task1"].Res)
	assert.Equal(t, 50, results["task1"].Time)
	assert.EqualError(t, results["task1"].Err, "error in task1")

	assert.Equal(t, "result2", results["task2"].Res)
	assert.Equal(t, 100, results["task2"].Time)
	assert.NoError(t, results["task2"].Err)

	task1.AssertExpectations(t)
	task2.AssertExpectations(t)
}

func TestTask_Execute(t *testing.T) {
	mockFunc := func(ctx context.Context, input int) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "Processed", nil
	}

	task := Task[int, string]{
		Func: mockFunc,
		Args: 42,
	}

	start := time.Now()
	result, duration, err := task.Execute(context.Background())
	elapsed := time.Since(start)

	assert.Equal(t, "Processed", result)
	assert.GreaterOrEqual(t, duration, 100) // Debe estar en el rango esperado
	assert.LessOrEqual(t, duration, int(elapsed.Milliseconds()))
	assert.NoError(t, err)
}

func TestTask_Execute_WithError(t *testing.T) {
	mockFunc := func(ctx context.Context, input int) (string, error) {
		return "", errors.New("execution failed")
	}

	task := Task[int, string]{
		Func: mockFunc,
		Args: 99,
	}

	result, duration, err := task.Execute(context.Background())

	assert.Empty(t, result)
	assert.GreaterOrEqual(t, duration, 0) // No hay espera real, debería ser inmediato
	assert.EqualError(t, err, "execution failed")
}
