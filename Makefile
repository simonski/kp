default_target: usage
.PHONY : default_target usage clean test build install publish format setup all

PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

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
	@echo "  publish               - build, cross-compile, tag, push, gh release, update homebrew tap"
	@echo "  setup                 - install dev dependencies (staticcheck, bn)"
	@echo ""

setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go get github.com/simonski/bn
	go install github.com/simonski/bn

clean:
	go clean
	rm -rf dist

test:
	go test

build: clean test format
	bn revision
	go fmt
	go build
	codesign -s - ./kp

install:
	go install

dist: clean test format
	bn revision
	go fmt
	mkdir -p dist
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		echo "Building $${os}/$${arch}..."; \
		GOOS=$${os} GOARCH=$${arch} go build -o "dist/kp_$${os}_$${arch}/kp" . ; \
		tar -czf "dist/kp_$${os}_$${arch}.tar.gz" -C "dist/kp_$${os}_$${arch}" kp ; \
	done

publish: dist
	$(eval VERSION := $(shell cat Buildnumber))
	git add -A
	git commit -m "release v$(VERSION)" || true
	git tag -a "v$(VERSION)" -m "v$(VERSION)"
	git push origin main --tags
	gh release create "v$(VERSION)" dist/kp_*.tar.gz --title "v$(VERSION)" --generate-notes
	@# Update homebrew tap
	$(eval DARWIN_AMD64_SHA := $(shell shasum -a 256 dist/kp_darwin_amd64.tar.gz | awk '{print $$1}'))
	$(eval DARWIN_ARM64_SHA := $(shell shasum -a 256 dist/kp_darwin_arm64.tar.gz | awk '{print $$1}'))
	$(eval LINUX_AMD64_SHA := $(shell shasum -a 256 dist/kp_linux_amd64.tar.gz | awk '{print $$1}'))
	$(eval LINUX_ARM64_SHA := $(shell shasum -a 256 dist/kp_linux_arm64.tar.gz | awk '{print $$1}'))
	rm -rf /tmp/homebrew-tap
	git clone git@github.com:simonski/homebrew-tap.git /tmp/homebrew-tap
	mkdir -p /tmp/homebrew-tap/Formula
	@echo 'class Kp < Formula' > /tmp/homebrew-tap/Formula/kp.rb
	@echo '  desc "A terminal tool to manage encrypted key/value pairs"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  homepage "https://github.com/simonski/kp"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  license "MIT"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  version "$(VERSION)"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  on_macos do' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    if Hardware::CPU.arm?' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      url "https://github.com/simonski/kp/releases/download/v$(VERSION)/kp_darwin_arm64.tar.gz"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      sha256 "$(DARWIN_ARM64_SHA)"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    else' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      url "https://github.com/simonski/kp/releases/download/v$(VERSION)/kp_darwin_amd64.tar.gz"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      sha256 "$(DARWIN_AMD64_SHA)"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    end' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  end' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  on_linux do' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    if Hardware::CPU.arm?' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      url "https://github.com/simonski/kp/releases/download/v$(VERSION)/kp_linux_arm64.tar.gz"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      sha256 "$(LINUX_ARM64_SHA)"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    else' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      url "https://github.com/simonski/kp/releases/download/v$(VERSION)/kp_linux_amd64.tar.gz"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '      sha256 "$(LINUX_AMD64_SHA)"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    end' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  end' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  def install' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    bin.install "kp"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  end' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  test do' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '    system "#{bin}/kp", "version"' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo '  end' >> /tmp/homebrew-tap/Formula/kp.rb
	@echo 'end' >> /tmp/homebrew-tap/Formula/kp.rb
	cd /tmp/homebrew-tap && git add Formula/kp.rb && git commit -m "Update kp to v$(VERSION)" && git push
	rm -rf /tmp/homebrew-tap
	@echo ""
	@echo "Published v$(VERSION) — brew install simonski/tap/kp"

format:
	staticcheck ./...
	go fmt ./...

all: clean build test install
