BINARY_NAME=youtube-search-cli
GO=go

all: build

build:
	$(GO) build -o $(BINARY_NAME) .

run:
	$(GO) run .

example:
	$(GO) run ./example/example.go

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

test:
	$(GO) test ./...

.PHONY: all build run example clean test