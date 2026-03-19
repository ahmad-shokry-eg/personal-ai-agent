.PHONY: install build clean

BIN_NAME := help
INSTALL_PATH := /usr/local/bin
HISTORY_FILE := $(HOME)/.help_history

build:
	go build -o $(BIN_NAME) .

install: build
	@echo "Installing $(BIN_NAME) to $(INSTALL_PATH)... (may require sudo)"
	@mv $(BIN_NAME) $(INSTALL_PATH)/$(BIN_NAME) || sudo mv $(BIN_NAME) $(INSTALL_PATH)/$(BIN_NAME)
	@touch $(HISTORY_FILE)
	@echo "Installed $(BIN_NAME) to $(INSTALL_PATH)."
	@echo "History file created at $(HISTORY_FILE)."

clean:
	rm -f $(BIN_NAME)
