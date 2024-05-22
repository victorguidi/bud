run: build
	@./bud

build:
	@go build -v -o bud ./src/

.PHONY: run build
