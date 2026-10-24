package worker

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

type Job struct {
	ID      string
	Type    string
	Payload []byte
	Status  Status
}

func NewJob(id string, jobType string, payload []byte) *Job {
	j := Job{
		ID:      id,
		Type:    jobType,
		Payload: payload,
		Status:  StatusPending,
	}
	return &j
}

type Handler interface {
	Process(job *Job) error
}
