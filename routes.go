package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) routes() {
	s.router.Use(middleware.Heartbeat("/healthcheck"))

	s.router.Group(func(r chi.Router) {
		r.Get("/search", s.search)
	})

}
