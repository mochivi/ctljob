package cmd

import (
	"context"

	"github.com/mochivi/ctljob/internal/core"
	"github.com/mochivi/ctljob/internal/scheduler"
	"github.com/mochivi/ctljob/internal/scheduler/api"
	v1 "github.com/mochivi/ctljob/internal/scheduler/api/v1"
	"github.com/mochivi/ctljob/internal/store"
)

func main() {
	// core
	store := store.NewInMemoryStore()
	registry := core.NewStubRegistry()
	sched := scheduler.NewScheduler(store, registry)

	// handlers
	jobHandler := v1.NewJobHandler(sched)

	server := api.NewServer(
		&api.Config{Port: 8080},
		api.Container{
			Job: jobHandler,
		})

	server.Run(context.Background())
}
