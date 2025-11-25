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
   ```

   Create `.env.release` (already gitignored) at the repo root:
   ```bash
   cat <<'EOF' > .env.release
   GITHUB_TOKEN=ghp_yourtokenhere
   EOF
   ```

   Every release-related `make` target will automatically source this file and export the variables it defines (no need to add `export`). Keep the file private—never commit it.

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

Pick the command that matches the bump you want:

```bash
# Patch bump (v1.2.3 -> v1.2.4)
make release-patch

# Minor bump (v1.2.3 -> v1.3.0)
make release-minor

# Major bump (v1.2.3 -> v2.0.0)
make release-major

# Use an exact version
make release-version VERSION=v1.4.0
```

Each command will:

- Ensure the git working tree is clean.
- Auto-create the git tag (annotated).
- Run GoReleaser against that tag.
- Push the tag to `origin` once the release succeeds.
- Remove the local tag again if GoReleaser fails, so you can retry cleanly.
- Leave the release artifacts in `dist/`.

Need to use a tag that already exists? Create/tag/push manually, then run `make release`.

Every full release will:
- Build binaries for Linux, macOS, and Windows (both amd64 and arm64)
- Create archives (.tar.gz for Linux/macOS, .zip for Windows)
- Generate checksums
- Create a GitHub release with all the assets
- Generate a changelog from your git commits

## Quick Distribution to Friends

### Option 1: Quick Share (No GitHub Release)

1. Run `make release-snapshot`.
2. Zip the binary your friend needs from `dist/` (for example `dist/clai_darwin_arm64/clai`).
3. Tell them to drop the binary somewhere on their `$PATH` (e.g. `/usr/local/bin`) and run `chmod +x clai` once.

This is the fastest way to hand someone a working build.

### Option 2: GitHub Release Download

1. Run `make release` (requires a tag + `GITHUB_TOKEN`).
2. Share the link `https://github.com/misrab/clai/releases`.
3. Friends download the archive that matches their OS/arch, unzip, move `clai` to their `$PATH`, run `chmod +x clai`.

### Option 3: `go install`

If they already have Go installed:

```bash
go install github.com/misrab/clai@latest
```

The `clai` binary lands in `$GOPATH/bin` (or `~/go/bin` by default). Have them ensure that directory is on their `$PATH`.

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

