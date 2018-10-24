package acd

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/waltton/talktask/manager"
)

// Job ...
type Job struct{}

type pool struct {
	ctx      context.Context
	poolSize int
	jobs     chan Job
}

// New returns a manager.ServiceFunc
func New(ctx context.Context, poolSize int, jobs chan Job) manager.ServiceFunc {
	p := pool{ctx, poolSize, jobs}
	return p.run
}

func (p *pool) drain() {
	log.Printf("%d jobs were marked for recovery", len(p.jobs))

	for range p.jobs {
		// database ?!
	}

}

func (p *pool) run() error {
	var wg sync.WaitGroup

	wg.Add(p.poolSize)
	for i := 0; i < p.poolSize; i++ {
		go p.worker(&wg)
	}

	wg.Wait()

	close(p.jobs)

	p.drain()

	return nil
}

func (p *pool) worker(wg *sync.WaitGroup) {
	var job Job

loop:
	for {
		select {
		case job = <-p.jobs:
			p.do(job)
		case <-p.ctx.Done():
			break loop
		}
	}

	wg.Done()
}

func (p *pool) do(job Job) {
	time.Sleep(time.Second)
}
