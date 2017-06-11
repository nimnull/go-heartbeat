#!/bin/bash

GOOS=linux GOARCH=386 go build -o bin/node_agent.linux.x86
GOOS=linux GOARCH=amd64 go build -o bin/node_agent.linux.amd64

