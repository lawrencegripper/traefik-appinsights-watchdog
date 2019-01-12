# build stage
FROM golang:1.11.4-alpine
ENV GOBIN /go/bin
RUN apk add --update --no-progress openssl git wget bash gcc musl-dev && \ 
    rm -rf /var/cache/apk/* && \
    wget --quiet -O /go/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && \
    chmod +x /go/bin/dep && \
    go get github.com/golang/lint/golint && \
    go get -u honnef.co/go/tools/cmd/...


COPY ./Gopkg.* /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog/
WORKDIR /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog
RUN dep ensure --vendor-only -v

COPY . /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog/

RUN chmod +x ./codechecks.sh;
RUN chmod +x ./build.sh
