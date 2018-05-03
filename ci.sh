#/bin/sh
docker build . -t traefik-appinsights-watchdog 
docker run -it -v $PWD:/go/src/github.com/lawrencegripper/traefik-appinsights-watchdog traefik-appinsights-watchdog bash -f build.sh
