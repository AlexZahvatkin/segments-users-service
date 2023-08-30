run: ### Run docker-compose
	docker-compose up --build -d server && docker-compose logs -f
.PHONY: run

down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: down

.PHONY: build
build: 
	go build -o bin/app -v ./cmd/segments-users-service

.PHONY: test
test: 
	go test -v -timeout 30s ./...

.PHONY: exec 
exec: build
	./bin/app

.DEFAULT_GOAL := build