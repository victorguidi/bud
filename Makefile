run: build
	@./bud

build:
	@go build -v -o bud ./src/

sql:
	@sqlite3 ./data/userdata.db "VACUUM;"

.PHONY: run build
