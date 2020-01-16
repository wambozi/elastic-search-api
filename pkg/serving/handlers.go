package serving

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wambozi/elastic-search-api/m/pkg/searching"
)

// Response is a concrete representation of the response to the client calling the crawl
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleCrawl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b searching.SearchRequest
		var results *searching.Results

		w.Header().Set("Content-Type", "application/json")

		if r.Method == "POST" {
			err := json.NewDecoder(r.Body).Decode(&b)
			if err != nil {
				s.Log.Error(err)
				er := errorResponse{Error: err.Error()}
				ers, _ := json.Marshal(er)

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(ers)
				return
			}
			results = searching.Search(s.ElasticClient, r, b, s.Log)
		}

		if r.Method == "GET" {
			q, ok := r.URL.Query()["qt"]
			i, ok := r.URL.Query()["i"]

			if !ok || len(q[0]) < 1 || len(i[0]) < 1 {
				s.Log.Errorf("Bad request: Missing query string parameters")
				er := errorResponse{Error: "Bad request: Missing query string parameters"}
				ers, _ := json.Marshal(er)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(ers)
				return
			}

			req := searching.SearchRequest{
				SearchTerm: q[0],
				Index:      i[0],
				Fields:     []string{"meta.description^2", "meta.title", "source.h1", "source.h2", "source.p"},
			}
			results = searching.Search(s.ElasticClient, r, req, s.Log)
		}

		response, err := json.Marshal(results)
		if err != nil {
			es := fmt.Sprintf("Failed to marshal %+v", results)
			er := errorResponse{Error: es}
			ers, _ := json.Marshal(er)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(ers)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(response)
	}
}
