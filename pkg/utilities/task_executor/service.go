package task_executor

import (
	"context"
	"sync"
	"time"
)

func WorkerPool(ctx context.Context, tasks map[string]Tasker, numWorkers int) map[string]Result {
	taskChan := make(chan struct {
		id   string
		task Tasker
	}, len(tasks))
	resultChan := make(chan Result, len(tasks))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskChan {
				res, tTime, err := t.task.Execute(ctx)
				resultChan <- Result{ID: t.id, Res: res, Err: err, Time: tTime}
			}
		}()
	}

	go func() {
		for id, task := range tasks {
			taskChan <- struct {
				id   string
				task Tasker
			}{id, task}
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make(map[string]Result)
	for res := range resultChan {
		results[res.ID] = res
	}

	return results
}

func (t Task[I, O]) Execute(ctx context.Context) (interface{}, int, error) {
	start := time.Now()
	out, err := t.Func(ctx, t.Args)
	duration := time.Since(start)
	return out, int(duration.Milliseconds()), err
}
