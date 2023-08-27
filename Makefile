.PHONY: build
build: 
	go build -o bin/app -v ./cmd/segments-users-service

.PHONY: test
test: 
	go test -v -timeout 30s ./...

.PHONY: run 
run: build
	./bin/app

.DEFAULT_GOAL := build