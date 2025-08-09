BINARY=to
PKG=github.com/textonlyio/textonly-cli
LDFLAGS=-s -w -X main.version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev) -X main.commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: build
build:
	mkdir -p dist
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY) ./cmd/to

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	@echo "No linter configured; add golangci-lint if desired"

.PHONY: release
release:
	goreleaser release --clean --snapshot
