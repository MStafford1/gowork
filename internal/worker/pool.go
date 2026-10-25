package worker

import (
	"context"
	"log"
	"sync"
)

type Pool struct {
	runner *Runner
	numWorkers int
	jobQueueChannel chan *Job
	wg sync.WaitGroup
}

func NewPool(runner *Runner, numWorkers int, queueSize int) *Pool {
	return &Pool{
		runner: runner,
		numWorkers: numWorkers,
		jobQueueChannel: make(chan*Job, queueSize),
		wg: sync.WaitGroup{},
	}
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	for job := range p.jobQueueChannel {
		log.Printf("Worker id: %d starting job %s\n", id, job.ID)
		err := p.runner.Run(ctx, job)
		if err != nil {
			log.Printf("Error processing job: %s\n  worker id: %d\n %v", job.ID, id, err)
		}
	}
}

func (p *Pool) Start(ctx context.Context) {
	for v := range p.numWorkers {
		p.wg.Add(1)
		go p.worker(ctx, v)
	}
}

func (p *Pool) Submit(job *Job) {
	p.jobQueueChannel <- job
}

func (p *Pool) Shutdown() {
	close(p.jobQueueChannel)
	p.wg.Wait()
}