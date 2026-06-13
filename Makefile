APP_NAME := lucene-to-dsl
BUILD_DIR := dist

# Detect OS
UNAME_S := $(shell uname -s 2>/dev/null || echo Windows)
IS_WINDOWS := $(findstring MINGW,$(UNAME_S))$(findstring MSYS,$(UNAME_S))$(findstring CYGWIN,$(UNAME_S))$(findstring Windows,$(UNAME_S))

ifneq ($(IS_WINDOWS),)
    APP_EXT := .exe
else
    APP_EXT :=
endif

APP_NAME_WITH_EXT := $(APP_NAME)$(APP_EXT)

.PHONY: build run test clean

build:
	go build -o $(BUILD_DIR)/$(APP_NAME_WITH_EXT) ./cmd

run: build
	./$(BUILD_DIR)/$(APP_NAME_WITH_EXT)

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)
