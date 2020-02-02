NAME := $(shell basename ${PWD})
VERSION ?= $(shell git rev-parse --short HEAD)

build: darwin-x64 linux-x64

linux-x64: bin/$(NAME)-$(VERSION)-linux-x64

darwin-x64: bin/$(NAME)-$(VERSION)-darwin-x64

test: build

clean:
	rm -rf bin

bin/$(NAME)-$(VERSION)-linux-x64:
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-$(VERSION)-linux-amd64 .

bin/$(NAME)-$(VERSION)-darwin-x64:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-$(VERSION)-darwin-amd64 .

.PHONY: clean build linux-x64 darwin-x64
