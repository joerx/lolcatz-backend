NAME := lolcatzd
VERSION ?= $(shell git rev-parse --short HEAD)

IMAGE_REGISTRY ?= ghcr.io
IMAGE_REPO_NAME ?= $(IMAGE_REGISTRY)/$(shell git remote get-url origin | cut -d':' -f2 | sed 's/.git//')
IMAGE_TAG ?= latest
IMAGE_NAME ?= $(IMAGE_REPO_NAME):$(IMAGE_TAG)

validate: clean build

build: bin/$(NAME)-linux-x64 bin/$(NAME)-darwin-x64

test: 
	go test ./test

clean:
	rm -rf bin

bin/$(NAME)-linux-x64:
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-linux-amd64 ./cmd/server.go

bin/$(NAME)-darwin-x64:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-darwin-amd64 ./cmd/server.go

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-push: docker-build
	docker push $(IMAGE_NAME)

.PHONY: validate clean build test docker-build docker-push
