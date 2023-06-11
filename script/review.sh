#!/bin/bash

docker pull devops-lintflow/gerritdocker:test-lintflow
docker run --rm -d -p 8080:8080 -p 29418:29418 devops-lintflow/gerritdocker:test-lintflow run
