#/bin/sh
docker build . -t traefik-appinsights-watchdog 
docker run -it -v $PWD/bin:/go/src/github.com/lawrencegripper/traefik-appinsights-watchdog/bin traefik-appinsights-watchdog bash -f build.sh