BINARY_NAME=envpop

VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -s -w
GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

build:
	GOARCH=amd64 GOOS=darwin go build $(GOFLAGS) -o ${BINARY_NAME}-darwin *.go
	GOARCH=amd64 GOOS=linux go build $(GOFLAGS) -o ${BINARY_NAME}-linux *.go

optimize:
	if [ -x /usr/bin/upx ] || [ -x /usr/local/bin/upx ]; then upx --brute ${BINARY_NAME}-*; fi

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux

dep:
	go mod download

test:
	go test ./...

test_with_coverage:
	go install github.com/jstemmer/go-junit-report@latest
	go test -v 2>&1 | go-junit-report > report.xml

