package worker

import (
	"context"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

type Job struct {
	ID      uuid.UUID
	Type    string
	Payload []byte
	Status  Status
}

func NewJob(jobType string, payload []byte) *Job {
	j := Job{
		ID:      uuid.New(),
		Type:    jobType,
		Payload: payload,
		Status:  StatusPending,
	}
	return &j
}

type Handler interface {
	Process(ctx context.Context, job *Job) error
}
