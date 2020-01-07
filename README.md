![Elasticsearch Search API](docs/img/elastic-search-api.png)

[![Sonarcloud Status](https://sonarcloud.io/api/project_badges/measure?project=wambozi_elastic-search-api&metric=coverage)](https://sonarcloud.io/dashboard?id=wambozi_elastic-search-api)

[![Release](https://github.com/wambozi/elastic-search-api/workflows/Release/badge.svg)](https://github.com/wambozi/elastic-search-api/)

## Description

Golang API that returns search results from Elasticsearch.

## Dependencies

- `go 1.13.5^`

## Usage

To run locally: `go run $(go list github.com/wambozi/elastic-search-api/... | grep -v /vendor/)`

To test locally: `make test-local`

## Routes

### `GET /healthcheck`

Example: http://localhost:8080/healthcheck

- returns the HTTP statusCode of the API. i.e. `200, 403, 502`

### `GET /search?q=${search_term}&i=${index}`

Example: http://localhost:8080/search?q=r2d2&i=droids

```JSON
{
    "Total": {
        "value": 1,
        "relation": "eq"
    },
    "max_score": 0.2876821,
    "hits": [
        {
            "_index": "droids",
            "_type": "_doc",
            "_id": "1234",
            "_score": 0.2876821,
            "Source": {
                "Name": "R2D2",
                "Species": "Robot"
            }
        }
    ]
}

```

## Docker Container

Docker Hub: https://hub.docker.com/repository/docker/wambozi/elastic-search-api

To run:

```shell
docker run --rm -it -p 8080:8080 wambozi/elastic-search-api:latest
```
