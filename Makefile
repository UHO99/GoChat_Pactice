DB_URL=postgres://gochat:gochat@localhost:5432/gochat?sslmode=disable

server:
	go run ./step$(SERVER)_*/cmd

QUERY :=
ifdef NICKNAME
QUERY := nickname=$(NICKNAME)
endif
ifdef ROOM
ifdef NICKNAME
QUERY := $(QUERY)&room=$(ROOM)
else
QUERY := room=$(ROOM)
endif
endif

send:
ifdef QUERY
	wscat -c "ws://localhost:808$(SERVER)/ws?$(QUERY)"
else
	wscat -c "ws://localhost:808$(SERVER)/ws"
endif

sqlc:
	sqlc generate

postgres:
	docker compose up -d

createdb:
	docker compose exec postgres createdb --username=gochat --owner=gochat gochat

dropdb:
	docker compose exec postgres dropdb --username=gochat gochat

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

.PHONY: server send sqlc createdb dropdb migrateup migratedo postgres
