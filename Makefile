default_target: usage
.PHONY : default_target usage clean test build install publish release release-snapshot format setup all

usage:
	@echo "The kp Makefile"
	@echo ""
	@echo "Usage : make <command> "
	@echo ""
	@echo "commands"
	@echo ""
	@echo "  clean                 - cleans temp files"
	@echo "  test                  - runs tests"
	@echo "  build                 - cleans, tests, formats, builds binary"
	@echo "  install               - go install"
	@echo "  publish               - tags, pushes, and releases via goreleaser (brew installable)"
	@echo "  release               - goreleaser release (no tag/push)"
	@echo "  release-snapshot      - local snapshot build for testing"
	@echo "  setup                 - install dev dependencies (staticcheck, bn)"
	@echo ""

setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go get github.com/simonski/bn
	go install github.com/simonski/bn
	brew install --cask goreleaser/tap/goreleaser

clean:
	go clean

test:
	go test

build: clean test format
	bn revision
	go fmt
	go build
	codesign -s - ./kp

install:
	go install

publish: build
	$(eval VERSION := $(shell cat Buildnumber))
	git add -A
	git commit -m "release v$(VERSION)" || true
	git tag -a "v$(VERSION)" -m "v$(VERSION)"
	git push origin main --tags
	goreleaser release --clean

release:
	goreleaser release --clean

release-snapshot:
	goreleaser release --snapshot --clean --skip=publish

format:
	staticcheck ./...
	go fmt ./...

all: clean build test install
