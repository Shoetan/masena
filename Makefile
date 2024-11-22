GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

.PHONY: setup generate test docker-up docker-down migrate-up migrate-down run init

SQLC_VERSION=v1.27.0
MIGRATE_VERSION=v4.17.0

setup:
	go mod download
	GOBIN=$(GOBIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	GOBIN=$(GOBIN) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

generate:
	$(GOBIN)/sqlc generate -f internal/db/sqlc.yaml

test:
	go test -v ./...

integration-test:
	go test -v -tags=integration ./tests/integration/...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

migrate-up:
	$(GOBIN)/migrate -path="internal/db/migrations" -database "postgres://bookstore:bookstore@localhost:5432/bookstore?sslmode=disable" up

migrate-down:
	$(GOBIN)/migrate -path="internal/db/migrations" -database "postgres://bookstore:bookstore@localhost:5432/bookstore?sslmode=disable" down

# Run the application
run:
	DATABASE_URL="postgres://bookstore:bookstore@localhost:5432/bookstore?sslmode=disable" \
	PORT=8080 \
	go run ./cmd/api

# Initialize everything and run the app
init: setup docker-up migrate-up run