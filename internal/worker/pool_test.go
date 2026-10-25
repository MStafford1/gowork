package worker

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

type CountHandler struct {
	lock sync.Mutex
	count int
}

func (c *CountHandler) Process(ctx context.Context, job *Job) error {
	c.lock.Lock()
	defer c.lock.Unlock() 
	c.count++
	return nil
}

func TestPool(t *testing.T) {
	runner := NewRunner()
	ch := &CountHandler{}
	runner.Register("test", ch)

	pool := NewPool(runner, 5, 20)
	pool.Start(context.Background())

	for v := range 100 {
		pool.Submit(NewJob("test", []byte(fmt.Sprintf("test %d", v))))
	}

	pool.Shutdown()

	if ch.count != 100 {
		t.Errorf("Count did not perfectly match")
	}
}