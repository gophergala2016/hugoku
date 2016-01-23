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
	go get github.com/julienschmidt/httprouter
	go get golang.org/x/oauth2
	go get golang.org/x/oauth2/github
	go get github.com/google/go-github/github
	go get github.com/codegangsta/negroni
	go get github.com/phyber/negroni-gzip/gzip
	go get github.com/goincremental/negroni-session

dev-deps:
	go get github.com/golang/lint/golint
