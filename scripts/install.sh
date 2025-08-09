#!/usr/bin/env sh
set -euo pipefail

# Configuration
REPO="textonlyio/textonly-cli"
BINARY="to"
BASE_URL="${TO_BASE_URL:-https://textonly.io}"
INSTALL_DIR="/usr/local/bin"
USER_BIN="$HOME/.local/bin"

uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    linux|darwin) echo "$os" ;;
    *) echo "unsupported OS: $os" >&2; exit 1 ;;
  esac
}

uname_arch() {
  arch=$(uname -m)
  case "$arch" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) echo "unsupported arch: $arch" >&2; exit 1 ;;
  esac
}

checksum_cmd() {
  if command -v shasum >/dev/null 2>&1; then
    echo "shasum -a 256 -c -"
  elif command -v sha256sum >/dev/null 2>&1; then
    echo "sha256sum -c -"
  else
    echo "no sha256 tool found (need shasum or sha256sum)" >&2
    exit 1
  fi
}

latest_tag_textonly() {
  # Expects JSON: {"tag_name":"vX.Y.Z"}
  curl -fsSL "$BASE_URL/cli/latest.json" \
    | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"}]\+\)".*/\1/p' | head -n1
}

latest_tag_github() {
  curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"}]\+\)".*/\1/p' | head -n1
}

resolve_tag() {
  if [ -n "${TO_VERSION:-}" ]; then
    echo "$TO_VERSION"
    return
  fi
  tag="$(latest_tag_textonly || true)"
  if [ -n "$tag" ]; then
    echo "$tag"
    return
  fi
  latest_tag_github
}

have_url() {
  # Return 0 if URL exists
  curl -IfsL "$1" >/dev/null 2>&1
}

main() {
  OS=$(uname_os)
  ARCH=$(uname_arch)
  TAG=$(resolve_tag)
  if [ -z "$TAG" ]; then
    echo "could not resolve latest version tag" >&2
    exit 1
  fi

  VER="${TAG#v}"
  TARBALL="${BINARY}_${VER}_${OS}_${ARCH}.tar.gz"

  # Prefer textonly.io; fall back to GitHub
  TXT_TARBALL_URL="$BASE_URL/downloads/$TAG/$TARBALL"
  TXT_CHECKSUMS_URL="$BASE_URL/downloads/$TAG/checksums.txt"

  GH_TARBALL_URL="https://github.com/$REPO/releases/download/$TAG/$TARBALL"
  GH_CHECKSUMS_URL="https://github.com/$REPO/releases/download/$TAG/checksums.txt"

  TARBALL_URL="$TXT_TARBALL_URL"
  CHECKSUMS_URL="$TXT_CHECKSUMS_URL"
  if ! have_url "$CHECKSUMS_URL" || ! have_url "$TARBALL_URL"; then
    TARBALL_URL="$GH_TARBALL_URL"
    CHECKSUMS_URL="$GH_CHECKSUMS_URL"
  fi

  TMP=$(mktemp -d)
  trap 'rm -rf "$TMP"' EXIT

  echo "Downloading $TARBALL_URL"
  curl -fL "$TARBALL_URL" -o "$TMP/$TARBALL"

  echo "Verifying checksum"
  curl -fL "$CHECKSUMS_URL" -o "$TMP/checksums.txt"
  CHECKER="$(checksum_cmd)"
  (
    cd "$TMP"
    # Support both sha256sum and shasum formats (with or without leading \*)
    grep -E "[[:space:]][*]?$TARBALL$" checksums.txt | sh -c "$CHECKER"
  )

  echo "Extracting"
  tar -xzf "$TMP/$TARBALL" -C "$TMP"

  DEST="$INSTALL_DIR"
  if [ ! -w "$DEST" ]; then
    DEST="$USER_BIN"
    mkdir -p "$DEST"
    case :$PATH: in
      *:$USER_BIN:*) : ;;
      *) echo "Note: add $USER_BIN to your PATH" ;;
    esac
  fi

  install -m 0755 "$TMP/$BINARY" "$DEST/$BINARY"
  echo "Installed to $DEST/$BINARY"
}

main "$@"
