#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $( find . -maxdepth 1 -regex "./[a-z_]*" ); do
  go test -v -race -coverprofile=profile.out -covermode=atomic $d/...
  if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
  fi
done
