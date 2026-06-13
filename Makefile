APP_NAME := lucene-to-dsl
BUILD_DIR := bin
MAIN_FILE := ./cmd/main.go

ifeq ($(OS),Windows_NT)
    EXT := .exe
    MKDIR_CMD := if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
    RM_CMD := if exist "$(BUILD_DIR)" rmdir /s /q "$(BUILD_DIR)"
else
    EXT :=
    MKDIR_CMD := mkdir -p $(BUILD_DIR)
    RM_CMD := rm -rf $(BUILD_DIR)
endif

.PHONY: build clean

build:
	$(MKDIR_CMD)
	go build -o $(BUILD_DIR)/$(APP_NAME)$(EXT) $(MAIN_FILE)

clean:
	$(RM_CMD)
