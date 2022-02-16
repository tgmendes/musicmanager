package worker

import (
	"context"
	"fmt"
)

type Task struct {
	ReqID string
	F     func() error
}

type Pool struct {
	TasksChan  chan *Task
	NumWorkers int
}

func (p *Pool) Run(ctx context.Context) {
	for w := 0; w <= p.NumWorkers; w++ {
		go p.spawnWorker(ctx, w)
	}
}

func (p *Pool) AddTask(task *Task) {
	p.TasksChan <- task
}

func (p *Pool) spawnWorker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return
		case task := <-p.TasksChan:
			fmt.Printf("worker %d started on %s\n", id, task.ReqID)
			err := task.F()
			if err != nil {
				fmt.Printf("task %s errored: %s", task.ReqID, err)
			}
			fmt.Printf("worker %d finished on %s\n", id, task.ReqID)
		}
	}
}
