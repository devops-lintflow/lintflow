# lintflow

[![Actions Status](https://github.com/craftslab/lintflow/workflows/CI/badge.svg?branch=master&event=push)](https://github.com/craftslab/lintflow/actions?query=workflow%3ACI)
[![Docker](https://img.shields.io/docker/pulls/craftslab/lintflow)](https://hub.docker.com/r/craftslab/lintflow)
[![Go Report Card](https://goreportcard.com/badge/github.com/craftslab/lintflow)](https://goreportcard.com/report/github.com/craftslab/lintflow)
[![License](https://img.shields.io/github/license/craftslab/lintflow.svg?color=brightgreen)](https://github.com/craftslab/lintflow/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/tag/craftslab/lintflow.svg?color=brightgreen)](https://github.com/craftslab/lintflow/tags)



## Introduction

*lintflow* is a master of lint workers written in Go.

- See *[lintaosp](https://github.com/craftslab/lintaosp/)* as a worker of *lintflow*.



## Prerequisites

- Go >= 1.15.0
- gRPC == 1.26.0



## Build

```bash
git clone https://github.com/craftslab/lintflow.git

cd lintflow
make build
```



## Run

```bash
./lintflow --config-file="config.yml" --code-review="gerrit" --commit-hash="{hash}" --output-file="output.json"
```



## Docker

```bash
git clone https://github.com/craftslab/lintflow.git

cd lintflow
docker build --no-cache -f Dockerfile -t craftslab/lintflow:latest .
docker run -it -v /tmp:/tmp craftslab/lintflow:latest ./bin/lintflow --config-file="./etc/config.yml" --code-review="gerrit" --commit-hash="{hash}" --output-file="/tmp/output.json"
```



## Usage

```
usage: lintflow --code-review=CODE-REVIEW --commit-hash=COMMIT-HASH --config-file=CONFIG-FILE [<flags>]

Lint Flow

Flags:
  --help                     Show context-sensitive help (also try --help-long
                             and --help-man).
  --version                  Show application version.
  --code-review=CODE-REVIEW  Code review (bitbucket|gerrit|gitee|github|gitlab)
  --commit-hash=COMMIT-HASH  Commit hash (SHA-1)
  --config-file=CONFIG-FILE  Config file (.yml)
  --output-file=OUTPUT-FILE  Output file (.json|.txt|.xlsx)
```



## Settings

*lintflow* parameters can be set in the directory [config](https://github.com/craftslab/lintflow/blob/master/config).

An example of configuration in [config.yml](https://github.com/craftslab/lintflow/blob/master/config/config.yml):

```yaml
apiVersion: v1
kind: master
metadata:
  name: lintflow
spec:
  lint:
    - name: lintaosp
      host: 127.0.0.1
      port: 9090
      filter:
        include:
          extension:
            - .java
            - .xml
    - name: lintkernel
      host: 127.0.0.1
      port: 9091
      filter:
        include:
          extension:
            - .c
            - .diff
            - .h
    - name: lintlang
      host: 127.0.0.1
      port: 9092
      filter:
        include:
          extension:
            - .c
            - .cpp
            - .go
            - .h
            - .hpp
            - .java
            - .py
            - .sh
  review:
    - name: gerrit
      host: http://127.0.0.1
      port: 8080
      user: user
      pass: pass
      vote:
        approval: +1
        disapproval: -1
        label: Code-Review
        message: Voting Code-Review by lintflow
```



## Design

![design](design.png)



## License

Project License can be found [here](LICENSE).



## Reference

### Gerrit

- [get-change-detail](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change-detail)
- [get-content](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-content)
- [get-patch](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-patch)
- [query-changes](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#query-changes)
- [set-review](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-review)



### Misc

- [gRPC](https://grpc.io/docs/languages/go/)
- [protocol-buffers](https://developers.google.com/protocol-buffers/docs/proto3)
- [reviewdog](https://github.com/reviewdog/reviewdog)
