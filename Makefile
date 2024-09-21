# Dev
dev: up-infra
	go run cmd/api/main.go
	


# Docker Compose

up:
	docker compose up -d

up-infra:
	docker compose up -d --scale api=0

up-build:
	docker compose up -d --build

down:
	docker compose down --volumes --remove-orphans 

restart:
	docker compose restart

# Logging

api-logs:
	docker compose logs -f api

mongo-logs:
	docker compose logs -f mongo

# Tests

test-e2e-single:
	(cd test/e2e && npm run test-single)