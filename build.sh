#!/bin/bash

#GOOS=linux GOARCH=386 go build -o node_agent.linux.x86
#GOOS=linux GOARCH=amd64 go build -o node_agent.linux.amd64
GOOS=darwin GOARCH=amd64 go build -o node_agent.darwin.amd64
