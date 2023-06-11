#!/bin/bash

go env -w GOPROXY=https://goproxy.cn,direct

if [ "$1" = "report" ]; then
  CGO_ENABLED=1 go test -cover -covermode=atomic -coverprofile=coverage.txt -parallel 2 -race -v ./...
else
  list="$(go list ./... | grep -v review)"
  old=$IFS IFS=$'\n'
  for item in $list; do
    CGO_ENABLED=0 go test -cover -covermode=atomic -parallel 2 -v "$item"
  done
  IFS=$old
fi
