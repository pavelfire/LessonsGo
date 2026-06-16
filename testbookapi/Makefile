BINARY := bin/bookshop

.PHONY: build test lint run dc tidy

build:
	CGO_ENABLED=0 go build -o $(BINARY) ./cmd/server

test:
	go test ./...

lint:
	go vet ./...

run: build
	./$(BINARY) -addr :8080

dc:
	docker compose up --build

tidy:
	go mod tidy
