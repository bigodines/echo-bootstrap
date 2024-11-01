BINARY_NAME=echo-bootstrap

.PHONY: all test build clean run tidy

all: test build

test:
	go test -v -race -cover ./...

build:
	go build -o bin/${BINARY_NAME} server.go

clean:
	go clean
	rm -f bin/${BINARY_NAME}

run: build
	./bin/${BINARY_NAME}

tidy:
	go mod tidy