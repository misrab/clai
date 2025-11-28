.PHONY: test build build-webui dev-webui dev-go dev run clean install \
	release release-snapshot release-dry-run \
	release-patch release-minor release-major release-version \
	ensure-clean

NEXT_TAG_SCRIPT := $(abspath scripts/next_tag.sh)
GORELEASER ?= goreleaser
LOAD_RELEASE_ENV = if [ -f .env.release ]; then \
		echo "Loading .env.release"; \
		set -a; \
		. .env.release; \
		set +a; \
	fi;

test:
	go test ./...

build-webui:
	@echo "Building web UI..."
	@cd webui && npm install && npm run build

dev-webui:
	@cd webui && npm install && npm run dev

dev-go:
	@echo "Starting Go backend with hot reload..."
	@which air > /dev/null || (echo "Error: 'air' not found. Install it with: go install github.com/air-verse/air@latest" && exit 1)
	@air

dev:
	@echo "Starting development environment..."
	@echo "  - Go backend with hot reload (port 8080)"
	@echo "  - React frontend with HMR (port 5173)"
	@echo ""
	@which air > /dev/null || (echo "Error: 'air' not found. Install it with: go install github.com/air-verse/air@latest" && exit 1)
	@trap 'kill 0' EXIT; \
	air & \
	sleep 2 && \
	(cd webui && npm install && npm run dev)

build: build-webui
	@mkdir -p bin
	go build -o bin/clai .

run: build
	./bin/clai $(ARGS)

install: build-webui
	go install .

clean:
	rm -f bin/clai
	rm -rf dist
	rm -rf webui/dist

# Release commands (requires goreleaser to be installed)
# Install goreleaser: brew install goreleaser/tap/goreleaser
# or: go install github.com/goreleaser/goreleaser/v2@latest

release: ensure-clean
	@echo "Using existing git tag for release."
	@echo "Make sure you have:"
	@echo "  1. Created the desired tag (e.g., git tag -a v1.0.0 -m 'Release v1.0.0')"
	@echo "  2. Set GITHUB_TOKEN environment variable"
	@echo ""
	@$(LOAD_RELEASE_ENV) $(GORELEASER) release --clean

release-snapshot:
	@echo "Creating a snapshot release (no git tag required)..."
	@$(LOAD_RELEASE_ENV) $(GORELEASER) release --snapshot --clean

release-dry-run:
	@echo "Dry run of the release process..."
	@$(LOAD_RELEASE_ENV) $(GORELEASER) release --skip=publish --clean

release-patch: ensure-clean
	@TAG=$$($(NEXT_TAG_SCRIPT) patch); \
	echo "Auto bump patch: $$TAG"; \
	git tag -a $$TAG -m "Release $$TAG"; \
	if { $(LOAD_RELEASE_ENV) $(GORELEASER) release --clean; }; then \
		git push origin $$TAG 2>/dev/null || echo "Tag $$TAG already exists on remote (GoReleaser created it)"; \
		echo "Release $$TAG complete."; \
		echo "View release at https://github.com/misrab/clai/releases/tag/$$TAG"; \
	else \
		echo "GoReleaser failed; deleting tag $$TAG"; \
		git tag -d $$TAG >/dev/null 2>&1 || true; \
		exit 1; \
	fi

release-minor: ensure-clean
	@TAG=$$($(NEXT_TAG_SCRIPT) minor); \
	echo "Auto bump minor: $$TAG"; \
	git tag -a $$TAG -m "Release $$TAG"; \
	if { $(LOAD_RELEASE_ENV) $(GORELEASER) release --clean; }; then \
		git push origin $$TAG 2>/dev/null || echo "Tag $$TAG already exists on remote (GoReleaser created it)"; \
		echo "Release $$TAG complete."; \
		echo "View release at https://github.com/misrab/clai/releases/tag/$$TAG"; \
	else \
		echo "GoReleaser failed; deleting tag $$TAG"; \
		git tag -d $$TAG >/dev/null 2>&1 || true; \
		exit 1; \
	fi

release-major: ensure-clean
	@TAG=$$($(NEXT_TAG_SCRIPT) major); \
	echo "Auto bump major: $$TAG"; \
	git tag -a $$TAG -m "Release $$TAG"; \
	if { $(LOAD_RELEASE_ENV) $(GORELEASER) release --clean; }; then \
		git push origin $$TAG 2>/dev/null || echo "Tag $$TAG already exists on remote (GoReleaser created it)"; \
		echo "Release $$TAG complete."; \
		echo "View release at https://github.com/misrab/clai/releases/tag/$$TAG"; \
	else \
		echo "GoReleaser failed; deleting tag $$TAG"; \
		git tag -d $$TAG >/dev/null 2>&1 || true; \
		exit 1; \
	fi

release-version: ensure-clean
ifndef VERSION
	$(error VERSION is required, e.g. make release-version VERSION=v1.2.3)
endif
	@TAG=$$($(NEXT_TAG_SCRIPT) version $(VERSION)); \
	echo "Releasing explicit version: $$TAG"; \
	git tag -a $$TAG -m "Release $$TAG"; \
	if { $(LOAD_RELEASE_ENV) $(GORELEASER) release --clean; }; then \
		git push origin $$TAG 2>/dev/null || echo "Tag $$TAG already exists on remote (GoReleaser created it)"; \
		echo "Release $$TAG complete."; \
		echo "View release at https://github.com/misrab/clai/releases/tag/$$TAG"; \
	else \
		echo "GoReleaser failed; deleting tag $$TAG"; \
		git tag -d $$TAG >/dev/null 2>&1 || true; \
		exit 1; \
	fi

ensure-clean:
	@if ! git diff --quiet; then \
		echo "Working tree has unstaged changes. Please commit or stash them."; \
		exit 1; \
	fi
	@if ! git diff --cached --quiet; then \
		echo "Working tree has staged but uncommitted changes. Please commit or unstage."; \
		exit 1; \
	fi
