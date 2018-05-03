#/bin/bash
set -e 
set -o pipefail

./codechecks.sh 

echo "Building...."
GOOS=windows go build -o traefik-appinsights-watchdog -o ./bin/traefik-appinsgihts-watchdog.exe
GOOS=windows go build -o traefik-appinsights-watchdog -o ./bin/traefik-appinsgihts-watchdog