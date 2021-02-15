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
./lintflow --config-file="config.yml"
```



## Docker

```bash
git clone https://github.com/craftslab/lintflow.git

cd lintflow
docker build --no-cache -f Dockerfile -t craftslab/lintflow:latest .
docker run -it craftslab/lintflow:latest ./bin/lintflow --config-file="./etc/config.yml"
```



## Usage

```
TBD
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
    lintaosp:
      host: 127.0.0.1
      port: 9090
    lintkernel:
      host: 127.0.0.1
      port: 9091
    lintlang:
      host: 127.0.0.1
      port: 9092
  review:
    bitbucket:
    gerrit:
      host: 127.0.0.1
      port: 8080
      user:
      pass:
    gitee:
    github:
    gitlab:
```



## Design

![design](design.png)



## License

Project License can be found [here](LICENSE).



## Reference

- [gRPC](https://grpc.io/docs/languages/go/)
- [protocol-buffers](https://developers.google.com/protocol-buffers/docs/proto3)
