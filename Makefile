.PHONY: build test test-integration install clean all dev

all: build test

build:
	cd code-vault && go build ./...
	mkdir -p bin && go build -o bin/sonic cmd/sonic/main.go

test:
	cd code-vault && go test ./...
	go test ./...

test-integration:
	@echo "Running integration tests..."
	@echo "Note: These tests require Docker daemon and test data"
	cd test/integration && go test -v ./...

install:
	cp bin/sonic /usr/local/bin/
	sonic library update

clean:
	rm -f bin/sonic
	rm -rf tmp/

dev:
	go run cmd/sonic/main.go
