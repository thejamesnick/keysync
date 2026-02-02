BINARY_NAME=keysync
INSTALL_PATH=/usr/local/bin

.PHONY: all build install clean uninstall test verify

all: build

build:
	@echo "🏗️  Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) ./cmd/keysync

install: build
	@echo "📦 Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo mv bin/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installed! You can now run '$(BINARY_NAME)' from anywhere."

uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Uninstalled."

clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f verify_flow.sh
	@rm -rf tmp_verify

test:
	@go test ./...

verify:
	@chmod +x verify_flow.sh
	@./verify_flow.sh
