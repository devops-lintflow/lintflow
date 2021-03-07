#!/bin/bash

src=lint

#protoc --go_out=./${src} --go_opt=paths=source_relative --go-grpc_out=./${src} --go-grpc_opt=paths=source_relative ./${src}/lint.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./lint/lint.proto
