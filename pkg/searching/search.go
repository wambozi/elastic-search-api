package searching

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/sirupsen/logrus"
)

// SearchRequest represents a search request on the POST /search route
type SearchRequest struct {
	SearchTerm string   `json:"searchTerm"`
	Index      string   `json:"index"`
	Fields     []string `json:"fields"`
}

// Results represents the Results response coming from Elasticsearch when performing a query
type Results struct {
	Took     int  `json:"took,omitempty"`
	TimedOut bool `json:"timed_out,omitempty"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	}
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		}
		MaxScore float64 `json:"max_score,omitempty"`
		Results  []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source struct {
				Source struct {
					H1 []string `json:"h1,omitempty"`
					H2 []string `json:"h2,omitempty"`
					H3 []string `json:"h3,omitempty"`
					H4 []string `json:"h4,omitempty"`
					P  []string `json:"p,omitempty"`
				} `json:"Source"`
				Meta struct {
					Title       string `json:"title"`
					Description string `json:"Description"`
					Keywords    string `json:"Keywords"`
				} `json:"Meta"`
				URI string `json:"URI"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// Search takes an elasticsearch Client and SearchRequest and returns results for that request
func Search(elasticClient *elasticsearch.Client, s SearchRequest, logger *logrus.Logger) *Results {
	go func(es *elasticsearch.Client, logger *logrus.Logger, i string, q string) {
		// we don't care about a successful index response, so ignore it
		_, err := indexQuery(es, i, q)
		if err != nil {
			logger.Error(err)
		}
	}(elasticClient, logger, s.Index, s.SearchTerm)

	res, err := searchQuery(elasticClient, s.Index, s.SearchTerm, s.Fields)
	if err != nil {
		logger.Error(err)
	}
	return res
}

func indexQuery(es *elasticsearch.Client, i string, q string) (response string, err error) {
	var (
		b strings.Builder
		r map[string]interface{}
	)

	b.WriteString(`{"query" : "`)
	b.WriteString(q)
	b.WriteString(`"}`)

	qb := []byte(q)
	tb := []byte(time.Now().Format("2006-01-02 15:04:05"))
	idBytes := md5.Sum(append(qb, tb...))
	idHash := hex.EncodeToString(idBytes[:])
	indexReq := esapi.IndexRequest{
		Index:      i + "-queries",
		DocumentID: idHash,
		Body:       strings.NewReader(b.String()),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := indexReq.Do(context.Background(), es)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		err := fmt.Errorf("[%s] Error indexing document ID=%s", res.Status(), idHash)
		return "", err
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		e := fmt.Errorf("Error parsing the index response body: %s", err)
		return "", e
	}

	// Print the response status and indexed document version.
	responseString := fmt.Sprintf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	return responseString, err
}

// Query represents the query to Elasticsearch
type Query struct {
	Query struct {
		MultiMatch struct {
			Query  string   `json:"query"`
			Fields []string `json:"fields"`
		} `json:"multi_match"`
	} `json:"query"`
}

func searchQuery(es *elasticsearch.Client, i string, q string, f []string) (r *Results, err error) {
	var (
		buf bytes.Buffer
	)

	query := Query{}
	query.Query.MultiMatch.Query = q
	query.Query.MultiMatch.Fields = f

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	searchRes, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(i),
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(searchRes.Body).Decode(&e); err != nil {
			jErr := fmt.Errorf("Error parsing the response body: %s", err)
			return nil, jErr
		}
		// Print the response status and error information.
		err = fmt.Errorf("[%s] %s: %s",
			searchRes.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
		return nil, err
	}

	if err := json.NewDecoder(searchRes.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}
