.PHONY: build test clean all dev

all: build test

build:
	go build ./cmd/...
	go build ./pkg/...

test:
	go test ./cmd/...
	go test ./pkg/...

clean:
	rm -f sonic
	rm -rf tmp/

dev:
	go run ./cmd/sonic/main.go
