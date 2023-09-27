#!/bin/bash
cd "$(dirname "$0")"
podman run --rm -v .:/mnt -w /mnt -e GO111MODULE=auto golang:1-alpine sh -c "go fmt main.go; go build -o smtpcli"
