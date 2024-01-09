VERSION=$(shell git describe --tags --always --dirty)

.PHONY: build-mac build-linux install-dependencies-mac install-dependencies-linux

build-mac:
	mkdir -p build
	CGO_ENABLED=1 CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib" \
		go build -o build/ovpn-auth-darwin-x64-$(VERSION) .

build-linux:
	mkdir -p build
	CGO_ENABLED=1 go build -o build/ovpn-auth-linux-x64-$(VERSION) .

install-dependencies-mac:
	cd "$$(mktemp -d)" && \
		pwd -P && \
		curl -sSLo 20190702.tar.gz https://github.com/P-H-C/phc-winner-argon2/archive/refs/tags/20190702.tar.gz && \
		tar -xvf 20190702.tar.gz && \
		cd phc-winner-argon2-20190702 && \
		sudo make install PREFIX=/usr/local
	brew install oath-toolkit
	pip install qrcode

install-dependencies-linux:
	cd "$$(mktemp -d)" && \
		pwd -P && \
		wget -sSLo 20190702.tar.gz https://github.com/P-H-C/phc-winner-argon2/archive/refs/tags/20190702.tar.gz && \
		tar -xvf 20190702.tar.gz && \
		cd phc-winner-argon2-20190702 && \
		sudo make install
	apt install oathtool
	pip install qrcode