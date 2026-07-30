listroom:
	curl -X GET http://localhost:8086/rooms

createroom:
	curl -X POST http://localhost:8086/rooms -d '{"name":"$(ROOMNAME)"}'

createuser:
	curl -X POST http://localhost:8086/users -d '{"name":"$(USERNAME)"}'

.PHONY: listroom createroom createuser