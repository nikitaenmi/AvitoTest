.PHONY: build run test clean dev migrate start stop restart logs lint lint-fix e2e e2e-local e2e-clean reset-db pgadmin app loadtest env

ENV_SAMPLE = samples/.env.example

env:
	@if [ ! -f .env ]; then \
		echo "Copying $(ENV_SAMPLE) to .env..."; \
		cp $(ENV_SAMPLE) .env; \
		echo "Please update .env file with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi


build: env
	docker-compose build

run: env
	docker-compose up

start: build run

stop:
	docker-compose down

restart: stop start

clean:
	docker-compose down -v
	docker system prune -f

dev:
	go run ./cmd

logs:
	docker-compose logs -f

migrate:
	docker-compose exec postgres psql -U root -d pr_review -f /docker-entrypoint-initdb.d/001_init.sql

reset-db:
	docker-compose down -v
	docker-compose up -d postgres
	sleep 5
	docker-compose up -d app pgadmin

pgadmin:
	open http://localhost:5051

app:
	open http://localhost:8080


lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix


e2e:
	docker-compose -f docker-compose.e2e.yml down
	docker-compose -f docker-compose.e2e.yml up --build --abort-on-container-exit --exit-code-from e2e-tests

e2e-local:
	go test -v ./e2e/... -timeout=5m

e2e-clean:
	docker-compose -f docker-compose.e2e.yml down -v


loadtest:
	go run loadtest/main.go loadtest/config.go loadtest/test_helpers.go


test: e2e-local

all: lint test build