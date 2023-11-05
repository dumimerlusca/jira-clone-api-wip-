run:
	go run ./...

docker/build:
	docker build -t dumimerlusca/jira-clone-api .
docker/push:
	docker push dumimerlusca/jira-clone-api
startdb:
	docker compose -f ./.docker/docker-compose.yaml up -d
stopdb:
	docker compose -f ./.docker/docker-compose.yaml down
start-test-db:
	docker compose -f ./.docker/docker-compose-test.yaml up -d
stop-test-db:
	docker compose -f ./.docker/docker-compose-test.yaml down
migrations/up :
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db/migrations -database="${DATABASE_URL}" up
migrations/down :
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db/migrations -database="${DATABASE_URL}" down

migrations/goto:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db/migrations -database="${DATABASE_URL}" goto ${version}

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./db/migrations ${name}


## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db/migrations -database="${DATABASE_URL}" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./db/migrations -database="${DATABASE_URL}" version

test:
	go clean -testcache
	go test ./...

.PHONY: migrate/up migrate/down stopdb startdb test docker/build