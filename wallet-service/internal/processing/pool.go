package processing

import (
	"context"
	"sync"

	"github.com/E-meliss/wallet-service/internal/domain"
)

type Processor interface {
	Process(ctx context.Context, task Task) error
}

type Pool struct {
	wg      sync.WaitGroup
	workers int
	queue   *Queue
	stats   *Stats
	proc    Processor
}

func NewPool(workers int, queue *Queue, proc Processor, stats *Stats) *Pool {
	return &Pool{workers: workers, queue: queue, proc: proc, stats: stats}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.queue.Chan() {
				err := p.proc.Process(ctx, task)
				if err != nil {
					p.stats.IncFailed()
				} else {
					p.stats.IncProcessed()
				}
				if task.ResultCh != nil {
					task.ResultCh <- err
				}
			}
		}()
	}
}

func (p *Pool) Stop() {
	p.queue.Close()
	p.wg.Wait()
	_ = domain.Money(0)
}
