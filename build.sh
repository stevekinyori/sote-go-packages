#!/bin/bash
#
# This will build for the target operating systems
#
env GOOS=linux GOARCH=amd64 go build -o build/linux/packages -ldflags="-s -w"
env GOOS=windows GOARCH=amd64 go build -o build/windows/packages.exe -ldflags="-s -w"
env GOOS=darwin go build -o build/mac/packages -ldflags="-s -w"