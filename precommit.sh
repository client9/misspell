#!/bin/sh
set -x

type dmnt >/dev/null 2>&1 || go get -u github.com/client9/dmnt
docker run --rm  $(dmnt .) \
  -w /go/src/github.com/client9/misspell \
  nickg/misspell make
