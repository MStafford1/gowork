package worker

import (
	"fmt"
	"context"
)

type Runner struct {
	handlers map[string]Handler
}

func NewRunner() *Runner {
	return &Runner{
		handlers: make(map[string]Handler),
	}
}

func (r *Runner) Register(jobType string, h Handler) error {
	_, ok := r.handlers[jobType]
	if ok {
		return fmt.Errorf("job type %s already registered on %s. Please deregister first", jobType, h)
	}
	r.handlers[jobType] = h
	return nil
}

func (r *Runner) Deregister(jobType string) error {
	_, ok := r.handlers[jobType]
	if !ok {
		return fmt.Errorf("no job type %s registered.", jobType)
	}
	delete(r.handlers, jobType)
	return nil
}

func (r *Runner) Run(ctx context.Context, job *Job) error {
	handler, ok := r.handlers[job.Type]

	if !ok {
		job.Status = StatusFailed
		return fmt.Errorf("no handler registered for job type: %s", job.Type)
	}

	job.Status = StatusRunning
	err := handler.Process(ctx, job)
	if err != nil {
		job.Status = StatusFailed
		return fmt.Errorf("error processing job: %s, %w", job.ID, err)
	}
	job.Status = StatusCompleted
	return nil
}
