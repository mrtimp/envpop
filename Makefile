BINARY_NAME=envpop

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin *.go
	#GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux *.go

optimize:
	if [ -x /usr/local/bin/upx ]; then upx --brute ${BINARY_NAME}-*; fi

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	#rm ${BINARY_NAME}-linux

test:
	go test ./...

dep:
	go mod download

