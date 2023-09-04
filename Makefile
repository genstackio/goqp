ARCH ?= darwin-arm64
platform_temp = $(subst -, ,$(ARCH))
GOOS = $(word 1, $(platform_temp))
GOARCH = $(word 2, $(platform_temp))
GOPROXY = https://proxy.golang.org
BIN ?= goqp

export CI

arch:
	@echo $(shell go env GOOS)-$(shell go env GOARCH)

build: build-binary
build-binary: prepare-build-dir
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o build/bin/$(GOOS)/$(GOARCH)/$(BIN) ./cli

clean: clean-binary clean-build-dir
clean-binary:
	@rm -rf build/bin
clean-build-dir:
	@rm -rf build


format:
	@gofmt -w *.go

install:
	@go mod vendor

prepare-build-dir:
	@mkdir -p build

update-dependencies:
	@go mod tidy -v

test:
	@go test *.go

.DEFAULT_GOAL := install
.PHONY: arch \
		build build-binary \
		clean clean-binary clean-build-dir \
		format \
		install \
		prepare-build-dir \
		test
