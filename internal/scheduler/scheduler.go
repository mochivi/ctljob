package scheduler

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mochivi/ctljob/internal/core"
	"github.com/mochivi/ctljob/internal/store"
)

type Scheduler struct {
	store    store.JobStore
	registry core.Registry
}

func NewScheduler(store store.JobStore, registry core.Registry) *Scheduler {
	return &Scheduler{store: store, registry: registry}
}

type JobSubmission struct {
	Name string `json:"name"`
}

func (j JobSubmission) into() *core.Job {
	return &core.Job{
		ID:     uuid.NewString(),
		Name:   j.Name,
		Status: core.StatusPending,
	}
}

func (s *Scheduler) Get(id string) (*core.Job, error) {
	return s.store.Get(id)
}

func (s *Scheduler) Submit(job JobSubmission) error {
	if !s.registry.Exists(job.Name) {
		return errors.New("function not found")
	}
	if err := s.store.Create(job.into()); err != nil {
		return fmt.Errorf("failed to submit job: %s", err)
	}
	return nil
}

type StatusUpdate struct {
	Status core.Status `json:"status"`
}

func (s *Scheduler) UpdateStatus(id string, update StatusUpdate) (bool, error) {
	return s.store.UpdateStatus(id, update.Status)
}
