package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MStafford1/gowork/internal/worker"
)

type EmailHandler struct {
}

func (eh *EmailHandler) Process(ctx context.Context, job *worker.Job) error {
	if job.Type != "Email" {
		return fmt.Errorf("Job must be of type Email for this handler")
	}
	select { 
	case <- time.After(500 * time.Millisecond):
		log.Printf("The generic timeout won!")
		fmt.Printf("Sending an email lol for job %s: %s\n", job.ID, job.Payload)
		return nil
	case <- ctx.Done():
		log.Printf("The context was cancelled/timed out mid-job")
		return ctx.Err()
	}
}

type ETLHandler struct {}

func (etlh *ETLHandler) Process(ctx context.Context, job *worker.Job) error {
	if job.Type != "ETL" {
		return fmt.Errorf("Job must be of type ETL for this handler")
	}
	select { 
	case <- time.After(500 * time.Millisecond):
		log.Printf("The generic timeout won!")
		fmt.Printf("Running ETL pipeline for job %s: %s\n", job.ID, job.Payload)
		return nil
	case <- ctx.Done():
		log.Printf("The context was cancelled/timed out mid-job")
		return ctx.Err()
	}
}

func main() {
	var runner = worker.NewRunner()
	runner.Register("Email", &EmailHandler{})
	runner.Register("ETL", &ETLHandler{})

	pool := worker.NewPool(runner, 3, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 800 * time.Millisecond)
	defer cancel()
	pool.Start(ctx)
	
	for v := range 10 {
		if v % 2 == 0 {
			pool.Submit(worker.NewJob("Email", []byte(fmt.Sprintf("%d@gmail.com", v))))
			log.Printf("Submitted Email job %d", v)
		} else {
			pool.Submit(worker.NewJob("ETL", []byte(fmt.Sprintf("Processing data from DB number %d", v))))
			log.Printf("Submitted ETL job %d", v)
		}
	}
	pool.Shutdown()
}