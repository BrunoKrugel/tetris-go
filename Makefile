.PHONY: build test race cover lint windows run format

run:
	go run .

build:
	go build ./...

test:
	go test ./...

race:
	go test -race ./...

cover:
	go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

lint:
	golangci-lint run

windows:
	GOOS=windows GOARCH=amd64 go build -o tetris.exe .

format:
	goimports -w .
	go fmt ./...
	fieldalignment -fix ./...
