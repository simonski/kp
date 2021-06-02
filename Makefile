default_target: usage
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

init:
	go mod github.com/rakyll/statik

clean:
	go clean
	
build:
	go build
	
test:
	go test

install:
	go install

all: clean build test install release 
	go install

release:
	goreleaser --snapshot --skip-publish --rm-dist

docker: build
	GOOS=linux GOARCH=amd64 go build -o kp_linux
	docker build -t kp .

