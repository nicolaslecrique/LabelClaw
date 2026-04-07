
.PHONY: install fmt lint test build dev

install:
	cd backend && go mod tidy

fmt:
	cd backend && gofmt -w .

lint: fmt
	cd backend && go vet ./...

test:
	cd backend && go test ./...

build:
	mkdir -p backend/bin
	cd backend && go build -o ./bin/labelclaw ./cmd/labelclaw

dev: build
	cd backend && ./bin/labelclaw
