package task_executor

import "context"

type Tasker interface {
	Execute(ctx context.Context) (result interface{}, duration int, err error)
}

type Task[I, O any] struct {
	Func func(context.Context, I) (O, error)
	Args I
}

type Result struct {
	ID   string
	Err  error
	Res  interface{}
	Time int
}
