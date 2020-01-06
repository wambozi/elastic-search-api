package main

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func TestRoutes(t *testing.T) {
	s := Server{}
	s.logger = logrus.New()
	s.router = chi.NewRouter()
	s.routes()
}
