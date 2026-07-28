server1:
	go run ./step1_simple_echo/cmd

server2:
	go run ./step2_simple_multi_broadcast/cmd

send:
	wscat -c ws://localhost:8080/ws

.PHONY: server1 server2 send
