package main

import (
	"encoding/json"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	// Get the query string parameter values
	query, ok := r.URL.Query()["qt"]
	index, ok := r.URL.Query()["i"]
	if !ok {
		s.logger.Error("Url Params missing. Required: 'qt' = 'query term', 'i' = 'index'")
		return
	}

	go func(es *elasticsearch.Client, logger *logrus.Logger, i string, q string) {
		// we don't care about a successful index response, so ignore it
		_, err := indexQuery(es, i, q)
		if err != nil {
			logger.Error(err)
		}
	}(s.elasticClient, s.logger, index[0], query[0])

	res, err := searchQuery(s.elasticClient, index[0], query[0])
	if err != nil {
		s.logger.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.Hits)
}
