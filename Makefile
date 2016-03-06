CONTAINER=nickg/misspell

all: install lint native-test

install:
	go get -t ./...
	go build ./...
	go run cmd/genwords/*.go > ./words.go
	go install ./cmd/misspell

lint:
	golint ./...
	go vet ./...
	find . -name '*.go' | xargs gofmt -w -s

native-test: install
	go test .
	misspell *.md replace.go cmd/misspell/*.go
	[[ -f /scowl-wl/words.txt ]] && misspell /scowl-wl/words.txt

clean:
	rm -f *~
	go clean ./...
	git gc

native-ci: install lint native-test

test:
	docker run --rm \
		-e COVERALLS_REPO_TOKEN=$COVERALLS_REPO_TOKEN \
		-v $(PWD):/go/src/github.com/client9/misspell \
		-w /go/src/github.com/client9/misspell \
		${CONTAINER} \
		make native-ci

docker-build:
	docker build -t ${CONTAINER} .

console: docker-build
	docker run --rm -it ${CONTAINER} sh

.PHONY: ci docker-ci
