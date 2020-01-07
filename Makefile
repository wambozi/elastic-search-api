OUT := ./bin/elastic-search-api
PKG := github.com/wambozi/elastic-search-api
ELASTIC_VERSION := 7.5.1

VERSION := $(shell git describe --always --long)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . --name '*.go'  | grep -v /vendor/)

.PHONY: clean
clean:
	-@rm -rf ${OUT} ${OUT}-v*
	for elasticRunner in $$(docker ps -a --filter=name=elastic-test -q); do \
		docker stop $$elasticRunner; \
		docker rm -f $$elasticRunner; \
	done
	for network in $$(docker network ls | grep testing | awk '{print $$1}'); do \
		docker network rm $$network; \
	done

.PHONY: compile
compile:
	go env -w GOPRIVATE=github.com/wambozi/*
	export GOFLAGS="-mod=vendor"
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -o ${OUT}-${VERSION} -ldflags="-extldflags \"-static\" -w -s -X main.version=${VERSION}"

.PHONY: build
build: compile
	docker build --build-arg VERSION="${VERSION}" -t wambozi/elastic-search-api:${VERSION} .

.PHONY: publish
publish:
	docker login --username wambozi --password ${DOCKER_TOKEN}
	docker push wambozi/elastic-search-api:${VERSION}

.PHONY: format
format:
	@gofmt -w *.go $$(ls -d */ | grep -v /vendor/)

.PHONY: test
test: clean
	[ -d reports ] || mkdir reports
	docker network create testing --subnet=172.18.0.0/16 --gateway=172.18.0.1
	docker run -it --network testing --ip 172.18.0.2 -d --name elastic-test -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:${ELASTIC_VERSION}
	sleep 30
	go test --coverprofile=reports/cov.out $$(go list ./... | grep -v /vendor/)
	go tool cover -func=reports/cov.out

.PHONY: vet
vet:
	@go vet .

.PHONY: lint
lint:
	@for file in ${GO_FILES}; do \
		golint $$file; \
	done

.PHONY: terraform-deploy
terraform-deploy:
	terraform init \
	terraform get \
	terraform plan \
	terraform apply

.PHONY: serverless-deploy
serverless-deploy:
	@export VERSION=${VERSION}
	serverless deploy --stage="dev"