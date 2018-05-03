#/bin/bash
set -e 
set -o pipefail

./codechecks.sh 

echo "Building...."
rm -r -f ./bin
mkdir -p ./bin

GOOS=windows go build -o traefik-appinsights-watchdog -o ./bin/traefik-appinsights-watchdog.exe
GOOS=linux CGO_ENABLED=0 go build -o traefik-appinsights-watchdog -o ./bin/traefik-appinsights-watchdog
