CONTAINER=nickg/misspell

all: install lint test

install:
	go version
	go get -t ./...
	go build ./...
	go run cmd/genwords/*.go > ./words.go
	go install ./cmd/misspell

lint:
	golint ./...
	go vet ./...
	find . -name '*.go' | xargs gofmt -w -s

test: install
	go test .
	misspell *.md replace.go cmd/misspell/*.go

# the grep in line 2 is to remove misspellings in the spelling dictionary
# that trigger false positives!!
falsepositives: /scowl-wl/words.txt
	cat /scowl-wl/words.txt | \
		grep -v -E "binominal|causalities|conquerer|withing" | \
		misspell -debug -error 

clean:
	rm -f *~
	go clean ./...
	git gc

ci-native: install lint test falsepositives

# when development is already in a docker contianer
ci:
	docker run --rm \
		--volumes-from=godev \
		-e COVERALLS_REPO_TOKEN=$COVERALLS_REPO_TOKEN \
		-w /go/src/github.com/client9/misspell \
		${CONTAINER} \
		make ci-native

# for on-mac and travis builds
ci-travis:
	docker run --rm \
		-v $(PWD):/go/src/github.com/client9/misspell \
		-e COVERALLS_REPO_TOKEN=$COVERALLS_REPO_TOKEN \
		-w /go/src/github.com/client9/misspell \
		${CONTAINER} \
		make ci-native 

docker-build:
	docker build -t ${CONTAINER} .

console: docker-build
	docker run --rm -it ${CONTAINER} sh

.PHONY: ci ci-travis ci-native
