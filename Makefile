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

.PHONY: server send sqlc
