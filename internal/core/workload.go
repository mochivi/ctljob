package core

import "context"

type Workload interface {
	ID() string
	Execute(ctx context.Context) error
}

type workload struct {
	id   string
	Name string `json:"name"`
}

func NewWorkload(name string) workload {
	return workload{
		Name: name,
	}
}

func (w *workload) ID() string                        { return w.id }
func (w *workload) Execute(ctx context.Context) error { return nil }
