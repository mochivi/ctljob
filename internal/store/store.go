package store

import (
	"errors"

	"github.com/mochivi/ctljob/internal/core"
)

var (
	ErrNotFound = errors.New("job not found")
)

type JobStore interface {
	Create(job *core.Job) error
	Get(id string) (*core.Job, error)
	ClaimNext(workerID string) (*core.Job, error)
	UpdateStatus(id string, status core.Status)
	List(filter Filter) ([]*core.Job, error)
}

type Filter struct {
}
