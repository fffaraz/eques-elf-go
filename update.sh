#!/bin/bash
set -euxo pipefail

rm -rf bin build-appimage fyne-cross 

go get -u ./...
go mod tidy

~/go/bin/gofumpt -l -w -extra .
