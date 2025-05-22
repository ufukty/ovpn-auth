VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X 'github.com/ufukty/ovpn-auth/cmd/ovpn-auth/commands/version.Version=$(VERSION)'"

.PHONY: build

build:
	mkdir -p build/$(VERSION)
	GOOS=darwin  GOARCH=amd64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-darwin-amd64  ./cmd/ovpn-auth
	GOOS=darwin  GOARCH=arm64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-darwin-arm64  ./cmd/ovpn-auth
	GOOS=linux   GOARCH=amd64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-linux-amd64   ./cmd/ovpn-auth
	GOOS=linux   GOARCH=386   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-linux-386     ./cmd/ovpn-auth
	GOOS=linux   GOARCH=arm   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-linux-arm     ./cmd/ovpn-auth
	GOOS=linux   GOARCH=arm64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-linux-arm64   ./cmd/ovpn-auth
	GOOS=freebsd GOARCH=amd64 go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-freebsd-amd64 ./cmd/ovpn-auth
	GOOS=freebsd GOARCH=386   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-freebsd-386   ./cmd/ovpn-auth
	GOOS=freebsd GOARCH=arm   go build -trimpath $(LDFLAGS) -o build/$(VERSION)/ovpn-auth-$(VERSION)-freebsd-arm   ./cmd/ovpn-auth

install:
	go install -trimpath $(LDFLAGS) ./cmd/ovpn-auth