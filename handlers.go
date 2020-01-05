package main

import (
	"encoding/json"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

func (s *server) search(w http.ResponseWriter, r *http.Request) {
	index := "demo"
	keys, ok := r.URL.Query()["q"]
	if !ok || len(keys[0]) < 1 {
		s.logger.Error("Url Param 'q' is missing")
		return
	}

	go func(es *elasticsearch.Client, logger *logrus.Logger, i string, q string) {
		err := indexQuery(es, i, q)
		if err != nil {
			logger.Error(err)
		}
	}(s.elasticClient, s.logger, index, keys[0])

	res, err := searchQuery(s.elasticClient, index, keys[0])
	if err != nil {
		s.logger.Error(err)
	}

	json.NewEncoder(w).Encode(res)
}
