.PHONY: install lint test type-check build dev

GO ?= go
PNPM ?= pnpm

install:
	@command -v $(GO) >/dev/null 2>&1 || (echo "Go 1.24.x is required but '$(GO)' was not found in PATH." && exit 1)
	@command -v node >/dev/null 2>&1 || (echo "Node.js 22.x is required but 'node' was not found in PATH." && exit 1)
	@command -v $(PNPM) >/dev/null 2>&1 || (echo "pnpm 10.x is required but '$(PNPM)' was not found in PATH." && exit 1)
	cd backend && $(GO) mod download
	$(PNPM) --dir frontend install

lint:
	@command -v $(GO) >/dev/null 2>&1 || (echo "Go 1.24.x is required but '$(GO)' was not found in PATH." && exit 1)
	cd backend && test -z "$$(gofmt -l .)" || (echo "Run gofmt on backend sources before linting." && exit 1)
	cd backend && $(GO) vet ./...
	cd backend && golangci-lint run ./...
	$(PNPM) --dir frontend lint

test:
	@command -v $(GO) >/dev/null 2>&1 || (echo "Go 1.24.x is required but '$(GO)' was not found in PATH." && exit 1)
	cd backend && $(GO) test ./...
	$(PNPM) --dir frontend test
	$(PNPM) --dir frontend test:e2e

type-check:
	@command -v $(GO) >/dev/null 2>&1 || (echo "Go 1.24.x is required but '$(GO)' was not found in PATH." && exit 1)
	cd backend && $(GO) build ./...
	$(PNPM) --dir frontend type-check

build:
	@command -v $(GO) >/dev/null 2>&1 || (echo "Go 1.24.x is required but '$(GO)' was not found in PATH." && exit 1)
	$(PNPM) --dir frontend build
	mkdir -p bin
	cd backend && $(GO) build -o ../bin/labelclaw ./cmd/server

dev:
	@command -v $(GO) >/dev/null 2>&1 || (echo "Go 1.24.x is required but '$(GO)' was not found in PATH." && exit 1)
	./scripts/dev.sh
