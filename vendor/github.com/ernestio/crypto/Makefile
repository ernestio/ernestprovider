install:
	go install -v

build:
	go build -v ./...

deps:

dev-deps: deps

test:
	go test -v ./...

lint:
	golint ./...
	go vet ./...
