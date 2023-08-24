
migrate := /home/dumi/go/bin/migrate
DATABASE_URL := "postgresql://dumi:dumi@localhost:5435/jira-clone-db?sslmode=disable"

startdb:
	docker compose up -d
stopdb:
	docker compose down
migrateup:
	 $(migrate) -path db/migrations -database $(DATABASE_URL) -verbose up
migratedown:
	 $(migrate) -path db/migrations -database $(DATABASE_URL) -verbose down

.PHONY: migrateup migratedown stopdb startdb