.PHONY: default build test verify vet lint check run clean

default:
	echo "Opinionated Go CLI for generating backend service boilerplate"

build:
	go build -o bin/goforge main.go
	chmod +x bin/goforge

test:
	go test ./...

verify:
	go mod verify

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: verify test vet lint build

run: build
	./bin/goforge

clean:
	rm -rf ./bin/**
