# make debug PKG=./pkg/util
PKG ?= ./...
DLV_PORT ?= 38697

test:
	go test $(PKG)

debug:
	dlv test $(PKG) --headless --listen=:$(DLV_PORT) --api-version=2

build:
	go build

.PHONY: test debug build
