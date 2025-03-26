BinaryName := domain-checker
LDFLAGS := '-w -s -extldflags "-static"'

.PHONY: build

build:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o bin/$(BinaryName) cmd/main.go
	upx --lzma bin/$(BinaryName)