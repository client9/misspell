
all: install lint test

install:
	go run cmd/genwords/*.go > ./words.go
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
