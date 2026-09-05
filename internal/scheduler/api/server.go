package api

import (
	"context"
	"net/http"

	v1 "github.com/mochivi/ctljob/internal/scheduler/api/v1"
)

type Config struct {
	Port int
}

type Container struct {
	Job *v1.JobHandler
}

type Server struct {
	http.Server
	cfg *Config
}

func NewServer(cfg *Config, container Container) *Server {
	mux := http.NewServeMux()
	setup_router(mux, container)

	return &Server{
		Server: http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func setup_router(mux *http.ServeMux, container Container) {
	mux.HandleFunc("/health", v1.Health)
	mux.HandleFunc("/ready", v1.Ready)

	mux.HandleFunc("POST /jobs", container.Job.Submit)
	mux.HandleFunc("GET /jobs/{id}", container.Job.Get)
	mux.HandleFunc("PATCH /jobs/{id}", container.Job.UpdateStatus)
}
