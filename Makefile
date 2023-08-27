.PHONY: build
build: 
	go build -o bin/app -v ./cmd/app

.PHONY: test
test: 
	go test -v -timeout 30s ./...

.PHONY: run 
run: 
	./bin/app

.DEFAULT_GOAL := build