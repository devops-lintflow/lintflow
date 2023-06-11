# lintflow

[![Actions Status](https://github.com/devops-lintflow/lintflow/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/devops-lintflow/lintflow/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/devops-lintflow/lintflow)](https://goreportcard.com/report/github.com/devops-lintflow/lintflow)
[![License](https://img.shields.io/github/license/devops-lintflow/lintflow.svg?color=brightgreen)](https://github.com/devops-lintflow/lintflow/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/devops-lintflow/lintflow.svg?color=brightgreen)](https://github.com/devops-lintflow/lintflow/tags)



## Introduction

*lintflow* is a server of lint workers written in Go.

- See *[lintwork](https://github.com/devops-lintflow/lintwork/)* as a worker of *lintflow*.



## Prerequisites

- Go >= 1.18.0
- gRPC >= 1.36.0



## Run

```bash
make build
./bin/lintflow --config-file="config.yml" --code-review="gerrit" --commit-hash="{hash}" --output-file="output.json"
```



## Docker

```bash
docker build --no-cache -f Dockerfile -t devops-lintflow/lintflow:latest .
docker run devops-lintflow/lintflow:latest /lintflow --config-file="/config.yml" --code-review="gerrit" --commit-hash="{hash}" --output-file="/output.json"
```



## Compose

```bash
# Run workers
docker-compose -f docker-compose.yml pull
docker-compose -f docker-compose.yml up -d

# Stop workers
docker-compose -f docker-compose.yml stop
docker-compose -f docker-compose.yml rm -f
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
  --output-file=OUTPUT-FILE  Output file (.json|.txt)
```



## Settings

*lintflow* parameters can be set in the directory [config](https://github.com/devops-lintflow/lintflow/blob/main/config).

An example of configuration in [config.yml](https://github.com/devops-lintflow/lintflow/blob/main/config/config.yml):

```yaml
apiVersion: v1
kind: server
metadata:
  name: lintflow
spec:
  lint:
    - name: lintcpp
      host: 127.0.0.1
      port: 9090
      timeout: 300
      filter:
        include:
          extension:
            - .c
            - .cc
            - .cpp
            - .h
            - .hpp
          file:
            - message
          repo:
            - foo
    - name: lintjava
      host: 127.0.0.1
      port: 9091
      timeout: 300
      filter:
        include:
          extension:
            - .java
            - .xml
          file:
            - message
          repo:
            - foo
    - name: lintpython
      host: 127.0.0.1
      port: 9092
      timeout: 300
      filter:
        include:
          extension:
            - .py
          file:
            - message
          repo:
            - foo
    - name: lintshell
      host: 127.0.0.1
      port: 9093
      timeout: 300
      filter:
        include:
          extension:
            - .sh
          file:
            - message
          repo:
            - foo
  review:
    - name: gerrit
      host: http://127.0.0.1/
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



## Errorformat

- **JSON format**

```json
{
  "lint": [
    {
      "file": "name",
      "line": 1,
      "type": "Error",
      "details": "text"
    }
  ]
}
```

- **Text format**

```text
{lint}:{file}:{line}:{type}:{details}
```



## Issues

- Fix comments issue with [change.maxComments](https://gerrit-documentation.storage.googleapis.com/Documentation/3.3.3/config-gerrit.html#change.maxComments).

```
One or more comments were rejected in validation: Exceeding maximum number of comments: 5001 (existing) + 1 (new) > 5000
```



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
