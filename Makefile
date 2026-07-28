server:
	go run ./cmd

send:
	wscat -c ws://localhost:8080/ws

.PHONY: server send
