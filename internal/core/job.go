package core

import (
	"errors"
	"fmt"
	"slices"
)

type Status string

var (
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidTransition = errors.New("invalid status transition")
)

const (
	StatusPending   = "pending"
	StatusClaimed   = "claimed"
	StatusRunning   = "running"
	StatusFailed    = "failed"
	StatusSuccessed = "successed"
)

func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusClaimed, StatusRunning, StatusFailed, StatusSuccessed:
		return true
	default:
		return false
	}
}

type Job struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

var transformations = map[Status][]Status{
	StatusPending:   {StatusClaimed},
	StatusClaimed:   {StatusRunning},
	StatusRunning:   {StatusFailed, StatusSuccessed},
	StatusFailed:    {StatusPending},
	StatusSuccessed: {},
}

func (j *Job) UpdateStatus(new Status) error {
	if !new.Valid() {
		return ErrInvalidStatus
	}

	if slices.Contains(transformations[j.Status], new) {
		j.Status = new
		return nil
	}

	return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, j.Status, new)
}
