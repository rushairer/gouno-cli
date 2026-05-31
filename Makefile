BINARY_NAME := gouno-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/rushairer/gouno-cli/gouno.Version=$(VERSION)"

.PHONY: build run test clean

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
