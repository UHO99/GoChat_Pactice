server:
	go run ./step$(SERVER)_*/cmd

send:
ifdef NICKNAME
	wscat -c "ws://localhost:808$(SERVER)/ws?nickname=$(NICKNAME)"
else
	wscat -c ws://localhost:808$(SERVER)/ws
endif

.PHONY: server send
