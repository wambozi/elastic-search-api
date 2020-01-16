package searching

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/wambozi/elastic-search-api/m/pkg/clients"
)

var (
	ee       = "http://localhost:9200"
	username = "elastic"
	password = "changeme"
)

func TestSearch(t *testing.T) {
	l := logrus.New()
	cfg := clients.GenerateElasticConfig([]string{ee}, username, password)
	ec, err := clients.CreateElasticClient(cfg)
	if err != nil {
		t.Errorf("Unexpected error creating Elasticsearch client: %s", err)
	}

	searchReq := SearchRequest{
		SearchTerm: "test",
		Index:      "test",
		Fields:     []string{"text"},
	}
	bodyJSON, err := json.Marshal(searchReq)
	encodedBody := bytes.NewReader(bodyJSON)

	req, err := http.NewRequest("POST", "/search", encodedBody)
	actual := Search(ec, req, searchReq, l)

	print(actual)
}
