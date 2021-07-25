#!/bin/bash

go env -w GOPROXY=https://goproxy.cn,direct

if [ "$1" = "all" ]; then
  GIN_MODE=release go test -cover -covermode=atomic -parallel 2 -race -v ./...
else
  list="$(go list ./... | grep -v review)"
  old=$IFS IFS=$'\n'
  for item in $list; do
    GIN_MODE=release go test -cover -covermode=atomic -parallel 2 -race -v "$item"
  done
  IFS=$old
fi
