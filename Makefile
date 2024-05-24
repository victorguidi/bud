run: build
	@./budengine

build: cli
	@go build -v -o budengine ./src/

cli:
	@go build -v -o bud ./cmd/

sql:
	@sqlite3 ./data/userdata.db "VACUUM;"

.PHONY: run build
