package core

import (
	"fmt"
	"time"
)

type Registry interface {
	Exists(name string) bool
}

type StubRegistry struct{}

func NewStubRegistry() Registry {
	return &StubRegistry{}
}

func (r *StubRegistry) Exists(name string) bool {
	for i := range stubRegistry {
		if i == name {
			return true
		}
	}
	return false
}

var stubRegistry = map[string]any{"sleep": sleepHandler}

func sleepHandler() {
	fmt.Print("sleep_handler_started")
	time.Sleep(5 * time.Second)
	fmt.Print("sleep_handler_completed")
}
