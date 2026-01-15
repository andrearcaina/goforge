.PHONY: default build run clean

default:
	echo "Opinionated Go CLI for generating backend service boilerplate"

build:
	go build -o bin/goforge main.go
	chmod +x bin/goforge

run:
	bash -c "./bin/goforge"

clean:
	rm -rf ./bin/**
