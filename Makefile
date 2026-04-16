.PHONY: build test install clean all dev

all: build test

build:
	cd code-vault && go build ./...
	cd sonic-screwdriver && mkdir -p bin && go build -o bin/sonic cmd/sonic/main.go

test:
	cd code-vault && go test ./...
	cd sonic-screwdriver && go test ./...

install:
	cp sonic-screwdriver/bin/sonic /usr/local/bin/
	sonic library update

clean:
	rm -f sonic-screwdriver/bin/sonic
	rm -rf sonic-screwdriver/tmp/

dev:
	cd sonic-screwdriver && go run cmd/sonic/main.go
