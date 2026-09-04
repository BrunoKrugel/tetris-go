.PHONY: build test race cover lint windows run format icon

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

windows: icon
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H windowsgui" -o tetris.exe .

icon:
	go run github.com/tc-hib/go-winres@v0.3.3 simply --icon internal/icon/icon.ico --manifest gui --arch amd64,arm64

format:
	goimports -w .
	go fmt ./...
	fieldalignment -fix ./...
