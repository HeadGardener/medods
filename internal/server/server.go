package server

import (
	"context"
	"net/http"
	"time"

	"github.com/HeadGardener/medods/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(conf config.ServerConfig, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + conf.Port,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
