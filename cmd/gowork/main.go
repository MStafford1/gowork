package main

import (
	"fmt"

	"github.com/MStafford1/gowork/internal/worker"
)

type EmailHandler struct {
}

func (eh *EmailHandler) Process(job *worker.Job) error {
	if job.Type != "Email" {
		return fmt.Errorf("Job must be of type Email for this handler")
	}
	fmt.Printf("Sending an email lol for job %s: %s\n", job.ID, job.Payload)
	return nil
}

type ETLHandler struct {}

func (etlh *ETLHandler) Process(job *worker.Job) error {
	if job.Type != "ETL" {
		return fmt.Errorf("Job must be of type ETL for this handler")
	}
	fmt.Printf("Running ETL pipeline for job %s: %s\n", job.ID, job.Payload)
	return nil
}

func main() {
	var runner = worker.NewRunner()
	runner.Register("Email", &EmailHandler{})
	runner.Register("ETL", &ETLHandler{})

	jobs := make([]*worker.Job, 0, 3)
	jobs = append(jobs, worker.NewJob("Email", []byte("me@gmail.com")))
	jobs = append(jobs, worker.NewJob("ETL", []byte("postgres")))
	jobs = append(jobs, worker.NewJob("Unknown", []byte("huh?")))

	for _, v := range jobs {
		err := runner.Run(v)
		if err != nil {
			fmt.Printf("Error processing job id: %s type: %s\n %w", v.ID, v.Type, err)
		}
	}
}