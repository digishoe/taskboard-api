.PHONY: test dev build

test:
	go test ./...

dev:
	go run .

build:
	go build -o taskboard-api .
