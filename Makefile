.PHONY: build test install clean all dev

all: build test

build:
	cd code-vault && go build ./...
	mkdir -p bin && go build -o bin/sonic cmd/sonic/main.go

test:
	cd code-vault && go test ./...
	go test ./...

install:
	cp bin/sonic /usr/local/bin/
	sonic library update

clean:
	rm -f bin/sonic
	rm -rf tmp/

dev:
	go run cmd/sonic/main.go
