# build stage
FROM golang:1.9.2-alpine AS build-env
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

RUN chmod +x codechecks.sh; sync; ./codechecks.sh

RUN go build -o traefik-appinsights-watchdog -v

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog/traefik-appinsights-watchdog .
ENTRYPOINT ["./traefik-appinsights-watchdog"]