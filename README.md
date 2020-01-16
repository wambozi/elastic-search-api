![Elasticsearch Search API](docs/img/elastic-search-api.png)

[![Sonarcloud Status](https://sonarcloud.io/api/project_badges/measure?project=wambozi_elastic-search-api&metric=coverage)](https://sonarcloud.io/dashboard?id=wambozi_elastic-search-api)

[![Release](https://github.com/wambozi/elastic-search-api/workflows/Release/badge.svg)](https://github.com/wambozi/elastic-search-api/)

## Description

Golang API that returns search results from Elasticsearch.

## Dependencies

- `go 1.13.5^`
- `Elasticsearch v7.5.1^`

## Configuration

Requires an config yaml in `conf`.

For instance:

Path: `/conf/local.yml`

```YAML
elasticsearch:
  endpoint: http://localhost:9200
  password: changeme
  username: elastic

appsearch:
  endpoint: http://localhost:3002
  api: /api/as/v1/
  token: private-pq7aaoSDFapSADosdnfns

server:
  port: 8080
  readHeaderTimeoutMillis: 3000
```


## Usage

### Running Local

Steps:

1. Launch Elasticsearch
2. Create a local config file:

```yaml
elasticsearch:
  endpoint: http://localhost:9200
server:
  port: 8080
  readHeaderTimeoutMillis: 3000
```

3. Install vendor dependencies: `go mod vendor`
4. Export env ID: `export ENV_ID=local`
5. Create an env config in `/conf` (example above). The name of this config should match the value of the env ID exported.
6. Compile (required to run the binary locally): `GO_ENABLED=0 go build -mod vendor -o ./bin/elastic-search-api ./cmd/elastic-search-api/main.go`
7. Run the compiled binary: `./bin/elastic-search-api`
8. Launch a crawl:

```shell
curl -XPOST localhost:8080/search -d '{
    "index": "demo",
    "searchTerm": "luke skywalker"
}'
```

### Running with Docker

This project builds and publishes a container with two tags, `latest` and `commit_hash`, to Docker Hub on merge to master.

Docker Hub: [https://hub.docker.com/repository/docker/wambozi/elastic-search-api](https://hub.docker.com/repository/docker/wambozi/elastic-search-api)

Steps:

1. Launch Elasticsearch
2. Create a local config file:

```yaml
elasticsearch:
  endpoint: http://docker.for.mac.localhost:9200
server:
  port: 8080
  readHeaderTimeoutMillis: 3000
```

3. Run the container.

```shell
docker pull wambozi/elastic-search-api:latest
docker run --rm -it -e "ENV_ID=local" -v "/some/path/to/conf:/conf" -p 8080:8080 wambozi/elastic-search-api:latest 
```

- `-v` : Mount the current dir into /conf dir of the container (so it makes local.yml accessible here). [Using bind mounts in docker](https://docs.docker.com/storage/bind-mounts/)
- `-e`: Required to specify the name of the env file. If `ENV_ID=local` isn't passed into the container, the container will exit with: `error: Error reading config file: Config File "no-config-set" Not Found in "[/conf /opt/bin/conf /opt/bin]"`
- `-p` expose the webserver port. This port should correspond to the value for `server.port` in your config.

4. If using App Search, [create the engine](https://swiftype.com/documentation/app-search/getting-started#engine) in App Search (API doesn't create it for you).
5. Launch a crawl:

```shell
curl -XPOST localhost:8080/search -d '{
    "index": "demo",
    "searchTerm": "luke skywalker"
}'
```

## Routes

### `POST /search`

Body:
```json
{
    "index": "demo",
    "searchTerm": "luke skywalker"
}
```

Response:
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

### `GET /search?q=${search_term}&i=${index}`

Response:

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

## Contributors

- [Adam Bemiller](https://github.com/adambemiller)
  - Adam provided most of the high level project and server/routes framework for this project. Huge thanks to him!

## License

MIT License
