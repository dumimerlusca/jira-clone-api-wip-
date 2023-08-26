
startdb:
	docker compose up -d
stopdb:
	docker compose down
start-test-db:
	docker compose -f ./docker-compose-test.yaml up -d
stop-test-db:
	docker compose -f ./docker-compose-test.yaml down
migrate/up :
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db/migrations -database="${DATABASE_URL}" up
migrate/down :
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db/migrations -database="${DATABASE_URL}" down

test:
	go clean -testcache
	go test ./...

.PHONY: migrate/up migrate/down stopdb startdb test