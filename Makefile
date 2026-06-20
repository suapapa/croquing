BINARY_NAME=croquis-king

all: build

web:
	cd frontend && npm install && npm run build
	go build -o bin/$(BINARY_NAME) ./cmd/server/main.go

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server/main.go

run:
	go run ./cmd/server/main.go

clean:
	rm -rf bin

test:
	go test -race ./...

progress:
	go run ./scripts/update_progress
