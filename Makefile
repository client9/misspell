
all: lint install test

install:
	(cd generators; go run wikipedia.go > ../lib/wikipedia.go)
	go install .

lint:
	golint ./...
	go vet ./...
	gofmt -w -s *.go */*.go

test:
	go test ./lib/...
	misspell README.md main.go lib/replace.go

clean:
	rm *~
	go clean ./...
