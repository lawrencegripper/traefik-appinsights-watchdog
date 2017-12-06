
#!/usr/bin/env bash
set -e
set -o pipefail

PKGS=$(go list ./... | grep -v '/vendor/') 
GOFILES=$(go list -f '{{range $index, $element := .GoFiles}}{{$.Dir}}/{{$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/') 
echo "----------> Running gofmt"
unformatted=$(gofmt -l $GOFILES)
if [ ! -z "$unformatted" ]; then
  echo "needs formatting: $unformatted"
  exit 1
fi
echo "----------> Running golint" 
golint -set_exit_status $PKGS 
echo "----------> Running simple"
gosimple $PKGS 
echo "----------> Running staticcheck" 
staticcheck $PKGS 
echo "----------> Running go test"
go test -v -cover $PKGS 