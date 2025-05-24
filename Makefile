run:
	go run ./cmd/main.go

build:
	go build -o bin/tixgo ./cmd/main.go

.PHONY: run build 