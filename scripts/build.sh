#!/bin/sh
set -ex
dep ensure
go install ./cmd/misspell
golangci-lint run
go test .
