package main

import (
	"github.com/elastic/go-elasticsearch/v8"
)

// CreateElasticClient returns the Elasticsearch client used by the function to connect to
// Elasticsearch, using the config provided
func CreateElasticClient(cfg elasticsearch.Config) (client *elasticsearch.Client, err error) {
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return client, err
	}

	return client, nil
}
