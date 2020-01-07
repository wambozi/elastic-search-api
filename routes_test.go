package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func TestHealthcheck(t *testing.T) {
	s := Server{}
	s.logger = logrus.New()
	s.router = chi.NewRouter()
	s.routes()

	endpoint := os.Getenv("ELASTICSEARCH_ENDPOINT")
	mockCfg := generateElasticConfig(endpoint)
	mockClient, err := CreateElasticClient(mockCfg)
	if err != nil {
		t.Fatal("Can't connect to elasticsearch endpoint.")
	}

	s.elasticClient = mockClient

	req := httptest.NewRequest("GET", "/healthcheck", nil)

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Healthcheck returned %d but expected %d", w.Result().StatusCode, http.StatusOK)
	}
}

func TestSearchRoute(t *testing.T) {
	s := Server{}
	s.logger = logrus.New()
	s.router = chi.NewRouter()
	s.routes()

	endpoint := os.Getenv("ELASTICSEARCH_ENDPOINT")
	mockCfg := generateElasticConfig(endpoint)
	mockClient, err := CreateElasticClient(mockCfg)
	if err != nil {
		t.Fatal("Can't connect to elasticsearch endpoint.")
	}

	s.elasticClient = mockClient

	req := httptest.NewRequest("GET", "/search?q=test&i=test", nil)

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	newStr := buf.String()

	fmt.Printf(newStr)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Search route returned %d but expected %d", w.Result().StatusCode, http.StatusOK)
	}
}
