.PHONY: run build test tidy gorm-gen validate

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

gorm-gen:
	go run ./cmd/gorm-gen

validate:
	ddd-bootstrap validate --project-root .
