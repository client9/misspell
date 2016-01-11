
all: lint install test

install:
	go run generators/*.go > ./lib/wikipedia.go
	go install .

lint:
	golint ./...
	go vet ./...
	gofmt -w -s *.go */*.go

test:
	go test ./lib/...
	misspell README.md main.go lib/replace.go generators/main.go

clean:
	rm -f *~
	go clean ./...
