#!/bin/bash
set -euxo pipefail

rm -rf bin build-appimage fyne-cross 

# export GOTOOLCHAIN=go1.25.10

go get -u ./...
go mod tidy

~/go/bin/gofumpt -l -w -extra .
