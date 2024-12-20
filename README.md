# lintflow

[![Actions Status](https://github.com/devops-lintflow/lintflow/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/devops-lintflow/lintflow/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/devops-lintflow/lintflow)](https://goreportcard.com/report/github.com/devops-lintflow/lintflow)
[![License](https://img.shields.io/github/license/devops-lintflow/lintflow.svg?color=brightgreen)](https://github.com/devops-lintflow/lintflow/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/devops-lintflow/lintflow.svg?color=brightgreen)](https://github.com/devops-lintflow/lintflow/tags)



## Introduction

*lintflow* is a server of lint workers written in Go.

- See *[lintwork](https://github.com/devops-lintflow/lintwork/)* as a worker of *lintflow*.



## Prerequisites

- Go >= 1.22.0



## Run

```bash
make build
./bin/lintflow --config-file="tests/config.yml" --code-review="gerrit" --commit-hash="{hash}"
```



## Docker

```bash
docker build --no-cache -f Dockerfile -t devops-lintflow/lintflow:latest .
docker run devops-lintflow/lintflow:latest /lintflow --config-file="config.yml" --code-review="gerrit" --commit-hash="{hash}"
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
  --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
  --[no-]version             Show application version.
  --code-review=CODE-REVIEW  Code review (bitbucket|gerrit|gitee|github|gitlab)
  --commit-hash=COMMIT-HASH  Commit hash (SHA-1)
  --config-file=CONFIG-FILE  Config file (.yml)
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
  flow:
    timeout: 120s
  lint:
    - name: lintai
      host: 127.0.0.1
      port: 9090
      filter:
        include:
          extension:
            - .c
            - .cc
            - .cpp
            - .h
            - .hpp
            - .java
          file:
          repo:
      vote: AI-Verified
    - name: lintcommit
      host: 127.0.0.1
      port: 9090
      filter:
        include:
          extension:
          file:
            - COMMIT_MSG
          repo:
      vote: Verified
    - name: lintcpp
      host: 127.0.0.1
      port: 9090
      filter:
        include:
          extension:
            - .c
            - .cc
            - .cpp
            - .h
            - .hpp
          file:
          repo:
      vote: Lint-Verified
    - name: lintjava
      host: 127.0.0.1
      port: 9090
      filter:
        include:
          extension:
            - .java
            - .xml
          file:
          repo:
      vote: Lint-Verified
    - name: lintpython
      host: 127.0.0.1
      port: 9090
      filter:
        include:
          extension:
            - .py
          file:
          repo:
      vote: Lint-Verified
    - name: lintshell
      host: 127.0.0.1
      port: 9090
      filter:
        include:
          extension:
            - .sh
          file:
          repo:
      vote: Lint-Verified
  review:
    name: gerrit
    url: http://127.0.0.1:8080
    user: user
    pass: pass
    vote:
      - label: AI-Verified
        approval: +1
        disapproval: -1
        message: Voting AI-Verified by lintflow
      - label: Lint-Verified
        approval: +1
        disapproval: -1
        message: Voting Lint-Verified by lintflow
      - label: Verified
        approval: +1
        disapproval: -1
        message: Voting Verified by lintflow
```



## Project

- **Commit Files**

```
lintwork-20240630231055/
├── COMMIT_MSG
├── {change-number}-{commit-id}.meta
├── {change-number}-{commit-id}.patch
└── path/to/file
```

- **Commit Meta**

```json
{
  "branch": "main",
  "owner": {
    "name": "name"
  },
  "project": "name",
  "revisions": {
    "39fe82c424a319e9613126d2ef1c837e114440c5": {
      "_number": 1
    }
  },
  "updated": "2024-09-20T07:15:44+08:00",
  "url": "http://127.0.0.1:8080"
}
```



## Report

- **JSON**

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

- **Text**

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
