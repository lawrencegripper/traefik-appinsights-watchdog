# build stage
FROM ataraev/golang-alpine-git AS build-env
ENV GOBIN /go/bin
RUN apk add --update openssl && \ 
    wget -O /go/bin/dep http://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && \
    chmod +x /go/bin/dep

COPY ./Gopkg.* /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog/
WORKDIR /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog
RUN dep ensure --vendor-only -v

COPY . /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog/
RUN go build -o traefik-appinsights-watchdog -v

# final stage
FROM golang:alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/lawrencegripper/traefik-appinsights-watchdog/traefik-appinsights-watchdog .
ENTRYPOINT ["./traefik-appinsights-watchdog"]