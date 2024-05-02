VERSION=$(shell git describe --tags --always --dirty)

.PHONY: build

build:
	mkdir -p build
	GOOS=darwin GOARCH=amd64 go build -o build/ovpn-auth-darwin-x64-$(VERSION) .
	GOOS=darwin GOARCH=arm64 go build -o build/ovpn-auth-darwin-arm-$(VERSION) .
	GOOS=linux GOARCH=amd64 go build -o build/ovpn-auth-linux-x64-$(VERSION) .
