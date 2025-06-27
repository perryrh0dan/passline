# make debug PKG=./pkg/util
PKG ?= ./...
DLV_PORT ?= 38697
COMMAND ?=

test:
	go test $(PKG)

test-debug:
	dlv test $(PKG) --headless --listen=:$(DLV_PORT) --api-version=2

build:
	go build

dev:
	go run main.go app.go $(COMMAND)

.PHONY: test debug build dev
