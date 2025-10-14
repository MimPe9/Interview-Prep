.PHONY: run build docker-up docker-down

run:
    go run cmd/server/main.go

build:
    go build -o bin/server cmd/server/main.go

docker-up:
    docker-compose -f compose.yml up -d

docker-down:
    docker-compose -f compose.yml down

test:
    go test ./...