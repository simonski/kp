VERSION := $(shell cat Buildnumber)
DIST    := dist

.PHONY: usage clean build test install setup release release-formula format

usage:
	@echo "The kp Makefile"
	@echo ""
	@echo "Usage : make <command>"
	@echo ""
	@echo "commands"
	@echo ""
	@echo "  usage                 - show this help (default)"
	@echo "  clean                 - run go clean"
	@echo "  build                 - creates native binary"
	@echo "  test                  - runs go test"
	@echo "  install               - go install"
	@echo "  setup                 - install dev dependencies (staticcheck, bn)"
	@echo "  release               - cross-compile and create GitHub release"
	@echo "  release-formula       - generate and push Homebrew formula to tap"
	@echo ""

clean:
	go clean
	rm -rf $(DIST)

build:
	staticcheck ./...
	bn revision
	go fmt .
	go build
	codesign -s - ./kp

test:
	go test

install:
	go install

setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go get github.com/simonski/bn
	go install github.com/simonski/bn

release: clean build test
	$(eval VERSION := $(shell cat Buildnumber))
	rm -rf $(DIST)
	mkdir -p $(DIST)
	GOOS=darwin  GOARCH=arm64 go build -o $(DIST)/kp && tar -czf $(DIST)/kp_$(VERSION)_darwin_arm64.tar.gz  -C $(DIST) kp && rm $(DIST)/kp
	GOOS=darwin  GOARCH=amd64 go build -o $(DIST)/kp && tar -czf $(DIST)/kp_$(VERSION)_darwin_amd64.tar.gz  -C $(DIST) kp && rm $(DIST)/kp
	GOOS=linux   GOARCH=arm64 go build -o $(DIST)/kp && tar -czf $(DIST)/kp_$(VERSION)_linux_arm64.tar.gz   -C $(DIST) kp && rm $(DIST)/kp
	GOOS=linux   GOARCH=amd64 go build -o $(DIST)/kp && tar -czf $(DIST)/kp_$(VERSION)_linux_amd64.tar.gz   -C $(DIST) kp && rm $(DIST)/kp
	gh release create v$(VERSION) $(DIST)/*.tar.gz --title "v$(VERSION)" --notes "Release v$(VERSION)"

release-formula:
	$(eval VERSION := $(shell cat Buildnumber))
	$(eval SHA_DARWIN_ARM64 := $(shell shasum -a 256 $(DIST)/kp_$(VERSION)_darwin_arm64.tar.gz | cut -d' ' -f1))
	$(eval SHA_DARWIN_AMD64 := $(shell shasum -a 256 $(DIST)/kp_$(VERSION)_darwin_amd64.tar.gz | cut -d' ' -f1))
	$(eval SHA_LINUX_ARM64  := $(shell shasum -a 256 $(DIST)/kp_$(VERSION)_linux_arm64.tar.gz  | cut -d' ' -f1))
	$(eval SHA_LINUX_AMD64  := $(shell shasum -a 256 $(DIST)/kp_$(VERSION)_linux_amd64.tar.gz  | cut -d' ' -f1))
	sed -e 's/{{VERSION}}/$(VERSION)/g' \
	    -e 's/{{SHA_DARWIN_ARM64}}/$(SHA_DARWIN_ARM64)/g' \
	    -e 's/{{SHA_DARWIN_AMD64}}/$(SHA_DARWIN_AMD64)/g' \
	    -e 's/{{SHA_LINUX_ARM64}}/$(SHA_LINUX_ARM64)/g' \
	    -e 's/{{SHA_LINUX_AMD64}}/$(SHA_LINUX_AMD64)/g' \
	    homebrew/kp.rb.tmpl > homebrew/kp.rb
	@echo "--- Generated Formula ---"
	@cat homebrew/kp.rb
	@echo "--- Pushing to simonski/homebrew-tap ---"
	gh api repos/simonski/homebrew-tap/contents/Formula/kp.rb \
		--method PUT \
		-f message="Update kp formula to $(VERSION)" \
		-f content="$$(base64 < homebrew/kp.rb)" \
		-f sha="$$(gh api repos/simonski/homebrew-tap/contents/Formula/kp.rb --jq .sha 2>/dev/null || echo '')" \
		|| gh api repos/simonski/homebrew-tap/contents/Formula/kp.rb \
			--method PUT \
			-f message="Add kp formula $(VERSION)" \
			-f content="$$(base64 < homebrew/kp.rb)"

format:
	staticcheck ./...
	go fmt ./...
