listroom:
	curl -X GET http://localhost:$(PORT)/rooms

createroom:
	curl -X POST http://localhost:$(PORT)/rooms -d '{"name":"$(ROOMNAME)"}'

createuser:
	curl -X POST http://localhost:$(PORT)/users -d '{"name":"$(USERNAME)"}'

.PHONY: listroom createroom createuser