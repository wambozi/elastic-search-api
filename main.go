package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Server represents this APIs HTTP server and corresponding methods
type Server struct {
	logger        *logrus.Logger
	router        *chi.Mux
	elasticClient *elasticsearch.Client
}

type config struct {
	Region          string
	ElasticEndpoint string
}

func main() {
	handler()
}

func handler() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load .env if its present (used for local dev)
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found.")
	}

	config := config{
		Region:          os.Getenv("REGION"),
		ElasticEndpoint: os.Getenv("ELASTICSEARCH_ENDPOINT"),
	}

	cfg := generateElasticConfig(config.ElasticEndpoint)

	s := Server{
		logger: logger,
	}
	s.router = chi.NewRouter()
	s.routes()

	elasticClient, err := CreateElasticClient(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	s.elasticClient = elasticClient

	const port = ":8080"
	server := http.Server{
		Addr:    port,
		Handler: s.router,
	}

	go func(server *http.Server) {
		logger.Info("Server listening on", port)
		if err := server.ListenAndServe(); err != nil {
			s.logger.Error(err.Error())
		}
	}(&server)

	// capture interrupt (ctrl-c)
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)

	// wait indefinitely until interrupt signal
	<-ctrlC

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		s.logger.Fatal(err.Error())
	}

}
