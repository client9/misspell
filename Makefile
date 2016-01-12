
all: install lint test

install:
	go run cmd/genwords/*.go > ./words.go
	go install .

lint:
	golint ./...
	go vet ./...
	find . -name '*.go' | xargs gofmt -w -s

test:
	go test ./...
	misspell *.md replace.go cmd/misspell/*.go

clean:
	rm -f *~
	go clean ./...
