default_target: usage
.PHONY : default_target upload

usage:
	@echo "The cryptic Makefile"
	@echo ""
	@echo "Usage : make <command> "
	@echo ""
	@echo "commands"
	@echo ""
	@echo "  clean                 - cleans temp files"
	@echo "  test                  - builds and runs tests"
	@echo "  build                 - creates binary"
	@echo "  run                   - runs in go run"
	@echo "  install               - builds and installs"
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

run: 
	go run main.go $*

all: clean build test
	go install
