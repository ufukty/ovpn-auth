VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X 'main.Version=$(VERSION)'"

.PHONY: build

build:
	mkdir -p build/$(VERSION)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-darwin-x64 ./cmd/ovpn-auth
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-darwin-arm ./cmd/ovpn-auth
	GOOS=linux  GOARCH=amd64 go build $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-linux-x64  ./cmd/ovpn-auth
