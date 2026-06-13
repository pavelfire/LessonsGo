GATEWAY_BIN := bin/gateway
BANKSIM_BIN := bin/banksim

.PHONY: build test lint run-gateway run-banksim run tidy docker-up docker-down

build:
	CGO_ENABLED=0 go build -o $(GATEWAY_BIN) ./cmd/gateway
	CGO_ENABLED=0 go build -o $(BANKSIM_BIN) ./cmd/banksim

test:
	go test ./...

lint:
	go vet ./...

run-banksim: build
	./$(BANKSIM_BIN) -addr :8081

run-gateway: build
	./$(GATEWAY_BIN) -addr :8080 -bank-url http://localhost:8081

run:
	@echo "Start bank simulator in one terminal: make run-banksim"
	@echo "Start gateway in another terminal: make run-gateway"

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down
