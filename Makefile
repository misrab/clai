.PHONY: test build run clean install release release-snapshot release-dry-run

test:
	go test ./...

build:
	@mkdir -p bin
	go build -o bin/clai .

run: build
	./bin/clai $(ARGS)

install:
	go install .

clean:
	rm -f bin/clai
	rm -rf dist

# Release commands (requires goreleaser to be installed)
# Install goreleaser: brew install goreleaser/tap/goreleaser
# or: go install github.com/goreleaser/goreleaser/v2@latest

release:
	@echo "Creating a new release..."
	@echo "Make sure you have:"
	@echo "  1. Committed all changes"
	@echo "  2. Tagged the release (e.g., git tag -a v1.0.0 -m 'Release v1.0.0')"
	@echo "  3. Set GITHUB_TOKEN environment variable"
	@echo ""
	goreleaser release --clean

release-snapshot:
	@echo "Creating a snapshot release (no git tag required)..."
	goreleaser release --snapshot --clean

release-dry-run:
	@echo "Dry run of the release process..."
	goreleaser release --skip=publish --clean
