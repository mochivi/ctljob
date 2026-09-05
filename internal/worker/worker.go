package worker

import (
	"context"

	"github.com/google/uuid"
	"github.com/mochivi/ctljob/internal/core"
)

type worker struct {
	id string
}

type Result struct {
	id  string
	err error
}

func (w *worker) run(ctx context.Context, ch <-chan core.Workload, resultsCh chan<- Result) {
	for {
		select {
		case workload := <-ch:
			if workload == nil {
				return
			}
			if err := workload.Execute(ctx); err != nil {
				resultsCh <- Result{
					id:  workload.ID(),
					err: err,
				}
				continue
			}
			resultsCh <- Result{
				id:  workload.ID(),
				err: nil,
			}
		case <-ctx.Done():
			return
		}
	}
}

type Worker struct {
	pool []*worker
	cfg  *Config

	workCh   chan core.Workload
	resultCh chan Result

	schedulerClient *schedulerClient
}

type Config struct {
	poolSize     int
	schedulerURL string
}

func NewWorker(cfg *Config) *Worker {
	pool := make([]*worker, 0, cfg.poolSize)

	for range cfg.poolSize {
		pool = append(pool, &worker{id: uuid.NewString()})
	}

	// how to find URL from worker eventually? config or service discovery?
	schedulerClient := NewSchedulerClient(cfg.schedulerURL)

	return &Worker{
		pool:            pool,
		cfg:             cfg,
		workCh:          make(chan core.Workload, max(cfg.poolSize, 100)),
		resultCh:        make(chan Result, max(cfg.poolSize, 100)),
		schedulerClient: schedulerClient,
	}
}

func (w *Worker) Run(ctx context.Context) {
	for _, worker := range w.pool {
		go worker.run(ctx, w.workCh, w.resultCh)
	}

	for result := range w.resultCh {
		if result.err != nil {
			continue
		}
		if err := w.schedulerClient.ReportSuccess(ctx); err != nil {
			continue // implement retries? work must be at-least-once?
			// need to think on how we're going to ensure idempotency
			// or if we need to!
		}
	}
}
