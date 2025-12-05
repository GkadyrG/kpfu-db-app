.PHONY: help run build tidy compose-up compose-down clean

help:
	@echo "Available commands:"
	@echo "  make run          - Run application locally"
	@echo "  make build        - Build application"
	@echo "  make tidy         - Tidy go modules"
	@echo "  make compose-up   - Start with Docker Compose"
	@echo "  make compose-down - Stop Docker Compose"
	@echo "  make clean        - Clean build artifacts"

run:
	go run cmd/main.go

build:
	go build -o bin/app cmd/main.go

tidy:
	go mod tidy

compose-up:
	docker-compose up --build

compose-down:
	docker-compose down

clean:
	rm -rf bin/
	docker-compose down -v

