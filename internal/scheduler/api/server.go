package api

import (
	"context"
	"net/http"

	v1 "github.com/mochivi/ctljob/internal/scheduler/api/v1"
)

type Config struct {
	Port int
}

type Server struct {
	http.Server
	cfg *Config
}

func NewServer(cfg *Config) *Server {
	mux := http.NewServeMux()
	setup_router(mux)

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

func setup_router(mux *http.ServeMux) {
	mux.HandleFunc("/health", v1.Health)
	mux.HandleFunc("/ready", v1.Ready)
}
