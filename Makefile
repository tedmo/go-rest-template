include .env

### Variables
SQLC_VERSION = 1.25.0
MIGRATE_VERSION = 4

### App
.PHONY: run
run:
	DB_USER=${DB_USER} \
	DB_PASSWORD=${DB_PASSWORD} \
	DB_HOST=${DB_HOST} \
	DB_PORT=${DB_PORT} \
	DB_NAME=${DB_NAME} \
	DB_SSL_MODE=${DB_SSL_MODE} \
	go run ./cmd/app/main.go

### Go
.PHONY: go/update
go/update:
	go get -u -t ./...
	go mod tidy

### SQLC
.PHONY: sqlc/gen
sqlc/gen:
	#sqlc generate
	docker run --rm -v ./:/src -w /src sqlc/sqlc:${SQLC_VERSION} generate

### Database
.PHONY: db/up
db/up:
	docker compose up --detach

.PHONY: db/down
db/down:
	docker compose down

### Schema
.PHONY: db/migrations/create
db/migration/new:
	docker run --rm -v ./migrations:/migrations --network host migrate/migrate:${MIGRATE_VERSION} \
        create -seq -ext=.sql -dir=/migrations ${name}

.PHONY: db/migrations/up
db/migrations/up:
	docker run --rm -v ./migrations:/migrations --network host migrate/migrate:${MIGRATE_VERSION} \
		-path=/migrations \
		-database=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE} \
		up

.PHONY: db/migrations/down
db/migrations/down:
	docker run --rm -v ./migrations:/migrations --network host migrate/migrate:${MIGRATE_VERSION} \
		-path=/migrations \
		-database=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE} \
		down
