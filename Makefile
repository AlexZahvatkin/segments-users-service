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

.PHONY: creatandmigrate
creatandmigrate:
	psql -U $(DB_USER) -w -c 'create database $(DB_DATABASE);'
	migrate -source file://internal/sql/postgresql/schema -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=$(DB_SSLMODE) up

.PHONY: migrate
migrate:
	psql -U $(DB_USER) -w -c 'create database $(DB_DATABASE);'
	migrate -source file://internal/sql/postgresql/schema -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=$(DB_SSLMODE) up

.DEFAULT_GOAL := build