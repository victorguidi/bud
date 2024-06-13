# Define the paths to the whisper include and library directories
C_INCLUDE_PATH=$(PWD)/src/include
LIBRARY_PATH=$(PWD)/src/lib

run: build
	@./budengine

# Define the build target
build: cli
	@cd ./src && tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
	@templ generate
	@C_INCLUDE_PATH=$(C_INCLUDE_PATH) LIBRARY_PATH=$(LIBRARY_PATH) go build -v -o budengine ./src/

cli:
	@go build -v -o bud ./cmd/

sql:
	@sqlite3 ./data/userdata.db "VACUUM;"

.PHONY: run build
