DB_URL=postgres://root:1234@localhost:5433/gochat?sslmode=disable

sqlc:
	sqlc generate

postgres:
	docker compose up -d

createdb:
	docker compose exec postgres createdb --username=root --owner=root gochat

dropdb:
	docker compose exec postgres dropdb --username=root gochat

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

.PHONY: sqlc createdb dropdb migrateup migratedo postgres