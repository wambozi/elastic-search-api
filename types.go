package main

// Elasticsearch Types

// DocumentObj represents the document source in the results coming from Elasticsearch
type DocumentObj struct{}

// ResultObj represents the results returned in "hits" from Elasticsearch
type ResultObj struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source DocumentObj
}

type hits struct {
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	}
	MaxScore float64     `json:"max_score,omitempty"`
	Results  []ResultObj `json:"hits"`
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
	Hits hits
}
