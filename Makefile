.PHONY: help run build build-web build-all tidy test fmt vet

SHELL   := /bin/sh
APP_NAME := gpt2api
BIN_DIR  := bin
WEB_DIR  := web
EMBED_DIR := internal/server/web

help:
	@echo "Targets:"
	@echo "  build-web   - npm run build → copy dist into embed dir"
	@echo "  build       - go build (embeds whatever is in $(EMBED_DIR))"
	@echo "  build-all   - build-web then build (single-binary release)"
	@echo "  run         - go run cmd/server (uses ./web/dist from disk if present)"
	@echo "  tidy        - go mod tidy"
	@echo "  test        - go test ./..."
	@echo "  fmt         - gofmt -w ."
	@echo "  vet         - go vet ./..."

# ---- frontend ----
build-web:
	@echo "[build-web] installing npm deps..."
	cd $(WEB_DIR) && npm install --no-audit --no-fund --loglevel=error
	@echo "[build-web] building frontend..."
	cd $(WEB_DIR) && npm run build
	@echo "[build-web] copying dist → $(EMBED_DIR)/"
	@rm -rf $(EMBED_DIR)
	@mkdir -p $(EMBED_DIR)
	@cp -r $(WEB_DIR)/dist/. $(EMBED_DIR)/
	@echo "[build-web] done."

# ---- backend (embeds current $(EMBED_DIR) content) ----
build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "-s -w" -o $(BIN_DIR)/$(APP_NAME) ./cmd/server
	@echo "[build] binary: $(BIN_DIR)/$(APP_NAME)"

# ---- combined: frontend + backend → single binary ----
build-all: build-web build

# ---- dev / util ----
run:
	go run ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...
