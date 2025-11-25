#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  scripts/next_tag.sh [patch|minor|major]
  scripts/next_tag.sh version <vX.Y.Z>

Examples:
  scripts/next_tag.sh patch
  scripts/next_tag.sh version v1.4.0
EOF
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

TYPE="$1"
shift || true

normalize_tag() {
  local raw="$1"
  if [[ "$raw" =~ ^v ]]; then
    echo "$raw"
  else
    echo "v$raw"
  fi
}

latest_tag() {
  if git describe --tags --abbrev=0 >/dev/null 2>&1; then
    git describe --tags --abbrev=0
  else
    echo "v0.0.0"
  fi
}

if [[ "$TYPE" == "version" ]]; then
  if [[ $# -lt 1 ]]; then
    echo "Error: version not provided" >&2
    usage
    exit 1
  fi
  TAG="$(normalize_tag "$1")"
else
  CURRENT="$(latest_tag)"
  CURRENT="${CURRENT#v}"
  IFS='.' read -r MAJOR MINOR PATCH <<<"$CURRENT"

  case "$TYPE" in
    patch)
      PATCH=$((PATCH + 1))
      ;;
    minor)
      MINOR=$((MINOR + 1))
      PATCH=0
      ;;
    major)
      MAJOR=$((MAJOR + 1))
      MINOR=0
      PATCH=0
      ;;
    *)
      echo "Error: unknown type '$TYPE'" >&2
      usage
      exit 1
      ;;
  esac

  TAG="v${MAJOR}.${MINOR}.${PATCH}"
fi

echo "$TAG"

