FROM golang:1.13.5-alpine3.11

RUN apk --update add bash wget dpkg-dev

RUN addgroup -S elastic && adduser -S elastic -G elastic

COPY --chown=elastic:elastic ./bin/elastic-search-api /opt/bin/elastic-search-api

USER elastic

WORKDIR /opt/bin

CMD [ "./elastic-search-api" ]
