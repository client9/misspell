#!/bin/sh
set -ex

# DEP
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# golangci-lint
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

# remove the default misspell to make sure
rm -f `which misspell`
