APP_NAME=miniredis
TARGET_DIR=bin

default: run

build:
	@go build -o $(TARGET_DIR)/$(APP_NAME) cmd/miniredis/main.go

clean:
	@rm -rf $(TARGET_DIR)

run: build
	@./$(TARGET_DIR)/$(APP_NAME)

watch:
	@air

.PHONY: default build clean run watch
