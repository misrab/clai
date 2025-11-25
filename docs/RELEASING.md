# Release Guide for clai

This guide explains how to create releases for clai using GoReleaser.

## Prerequisites

1. **Install GoReleaser**
   ```bash
   # macOS
   brew install goreleaser/tap/goreleaser
   
   # or using Go
   go install github.com/goreleaser/goreleaser/v2@latest
   ```

2. **GitHub Token** (for publishing releases)
   ```bash
   # Create a GitHub token with repo permissions at:
   # https://github.com/settings/tokens/new
   
   # Export it in your shell:
   export GITHUB_TOKEN="your_github_token_here"
   
   # Or add to your ~/.zshrc or ~/.bashrc:
   echo 'export GITHUB_TOKEN="your_github_token_here"' >> ~/.zshrc
   ```

## Release Process

### 1. Test the Release Process (Dry Run)

First, test that everything is configured correctly:

```bash
make release-dry-run
```

This will build binaries for all platforms but won't publish anything.

### 2. Create a Snapshot Release (No Git Tag Required)

For testing the full release process locally without creating a GitHub release:

```bash
make release-snapshot
```

This creates binaries in the `dist/` folder that you can share directly with friends.

### 3. Create a Production Release

When you're ready to create an official release:

```bash
# 1. Make sure all changes are committed
git add .
git commit -m "Prepare for release"

# 2. Create and push a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 3. Run the release
make release
```

This will:
- Build binaries for Linux, macOS, and Windows (both amd64 and arm64)
- Create archives (.tar.gz for Linux/macOS, .zip for Windows)
- Generate checksums
- Create a GitHub release with all the assets
- Generate a changelog from your git commits

## Quick Distribution to Friends

### Option 1: Snapshot Release (Easiest)

```bash
# Create snapshot builds
make release-snapshot

# Share the binaries from dist/ folder
# For example, for macOS ARM64:
# dist/clai_darwin_arm64/clai
```

### Option 2: GitHub Release

After running `make release`, your friends can download from:
```
https://github.com/misrab/clai/releases
```

### Option 3: Direct Install via Go

Once released, friends with Go installed can run:
```bash
go install github.com/misrab/clai@latest
```

## File Locations

After running a release:
- `dist/` - Contains all build artifacts
- `dist/checksums.txt` - SHA256 checksums for verification
- `dist/clai_<os>_<arch>/` - Individual binary folders
- `dist/*.tar.gz` or `dist/*.zip` - Distribution archives

## Versioning

Follow semantic versioning (semver):
- `v1.0.0` - Major release
- `v1.1.0` - Minor release (new features)
- `v1.0.1` - Patch release (bug fixes)

## Troubleshooting

**"git tag not found"**: Make sure you've created and pushed a git tag before running `make release`.

**"GITHUB_TOKEN not set"**: Export your GitHub token as shown in prerequisites.

**Build fails**: Run `make test` and `make build` first to ensure everything compiles correctly.

## Advanced: Homebrew Distribution

To distribute via Homebrew, uncomment the `brews:` section in `.goreleaser.yaml` and create a tap repository at `github.com/misrab/homebrew-tap`.

