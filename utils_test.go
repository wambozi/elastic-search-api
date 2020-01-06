package main

import (
	"os"
	"testing"
)

func TestIndexQuery(t *testing.T) {
	endpoint := "http://" + os.Getenv("ELASTICSEARCH_IP") + ":9200"
	mockCfg := generateElasticConfig(endpoint)
	mockClient, err := CreateElasticClient(mockCfg)
	if err != nil {
		t.Fatal("Can't connect to elasticsearch endpoint.")
	}

	_, err = indexQuery(mockClient, "test", "test query")

	if err != nil {
		t.Error(err)
	}
}

func TestSearchQuery(t *testing.T) {
	endpoint := "http://" + os.Getenv("ELASTICSEARCH_IP") + ":9200"
	mockCfg := generateElasticConfig(endpoint)
	mockClient, err := CreateElasticClient(mockCfg)
	if err != nil {
		t.Fatal("Can't connect to elasticsearch endpoint.")
	}

	_, err = searchQuery(mockClient, "test-queries", "test")

	if err != nil {
		t.Error(err)
	}
}
