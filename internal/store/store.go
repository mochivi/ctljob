package store

import (
	"errors"
	"sync"

	"github.com/mochivi/ctljob/internal/core"
)

var (
	ErrNotFound = errors.New("job not found")
)

type JobStore interface {
	Create(job *core.Job) error
	Get(id string) (*core.Job, error)
	ClaimNext(workerID string) (*core.Job, error)
	UpdateStatus(id string, status core.Status) (bool, error)
	List(filter Filter) ([]*core.Job, error)
}

type Filter struct{}

type InMemoryStore struct {
	mu    sync.Mutex
	jobs  map[string]*core.Job
	order []string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{jobs: make(map[string]*core.Job)}
}

func (s *InMemoryStore) Create(job *core.Job) error {
	if job == nil || job.ID == "" {
		return errors.New("job must have an id")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.jobs[job.ID]; ok {
		return errors.New("job already exists")
	}

	clone := *job
	s.jobs[job.ID] = &clone
	s.order = append(s.order, job.ID)
	return nil
}

func (s *InMemoryStore) Get(id string) (*core.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return nil, ErrNotFound
	}

	clone := *job
	return &clone, nil
}

func (s *InMemoryStore) UpdateStatus(id string, status core.Status) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return false, ErrNotFound
	}

	if job.Status == status {
		return false, nil
	}

	if err := job.UpdateStatus(status); err != nil {
		return false, err
	}
	return true, nil
}

func (s *InMemoryStore) ClaimNext(workerID string) (*core.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range s.order {
		job, ok := s.jobs[id]
		if !ok || job.Status != core.StatusPending {
			continue
		}
		if err := job.UpdateStatus(core.StatusClaimed); err != nil {
			return nil, err
		}
		clone := *job
		return &clone, nil
	}
	return nil, ErrNotFound
}

func (s *InMemoryStore) List(filter Filter) ([]*core.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobs := make([]*core.Job, 0, len(s.jobs))
	for _, id := range s.order {
		clone := *s.jobs[id]
		jobs = append(jobs, &clone)
	}
	return jobs, nil
}