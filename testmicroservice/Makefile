BINARY := bin/portloader
IMAGE  := portloader

.PHONY: build test lint run docker-build docker-run tidy

build:
	CGO_ENABLED=0 go build -o $(BINARY) ./cmd/portloader

test:
	go test ./...

lint:
	go vet ./...

run: build
	./$(BINARY) -file ports.json

tidy:
	go mod tidy

docker-build: build
	docker build -t $(IMAGE) .

docker-run: docker-build
	docker run --rm $(IMAGE)
