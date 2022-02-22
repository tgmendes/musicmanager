package worker

import (
	"context"
	"fmt"
)

type Task struct {
	ReqID string
	F     func(context.Context) error
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
			return
		case task := <-p.TasksChan:
			err := task.F(ctx)
			if err != nil {
				fmt.Printf("task %s errored: %s\n", task.ReqID, err)
				return
			}
			fmt.Printf("worker %d completed task %s!\n", id, task.ReqID)
		}
	}
}
