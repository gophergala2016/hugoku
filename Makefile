build:
	go build

install:
	go install

test:
	go test -v -race ./...

lint:
	go vet ./...
	golint ./...

clean:
	go clean

deps: dev-deps

dev-deps:
	go get github.com/golang/lint/golint
