default_target: build
.PHONY : default_target upload

usage:
	@echo "The kp Makefile"
	@echo ""
	@echo "Usage : make <command> "
	@echo ""
	@echo "commands"
	@echo ""
	@echo "  clean                 - cleans temp files"
	@echo "  test                  - builds and runs tests"
	@echo "  build                 - creates binary"
	@echo "  install               - builds and installs"
	@echo "  release               - creates the crossplatorm releases"
	@echo ""
	@echo "  all                   - all of the above"
	@echo ""

setup:	
	go install honnef.co/go/tools/cmd/staticcheck@latest
	
clean:
	go clean
	
build: clean test format
	go fmt
	go build
	
test:
	go test

install:
	go install

all: clean build test install release 
	go install

release:
	goreleaser --snapshot --skip-publish --rm-dist

format:
	staticcheck ./...
	go fmt ./...
