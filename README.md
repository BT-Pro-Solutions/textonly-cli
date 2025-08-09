TextOnly CLI (to)

A fast, single-binary CLI for TextOnly. One-line install, magic-link login, secure token storage, shell completions, and auto-updates.

Quick install
curl -fsSL https://textonly.io/install.sh | sh

First steps

to login
to notes create --title "Hello" --stdin <<<'This is my note'
to notes list

Commands

Auth:
- to login [--no-open]
- to logout
- to whoami [--json]

Notes:
- to notes list [--public|--private] [--limit N] [--json]
- to notes view <id|slug> [--raw] [--open] [--json]
- to notes create [--title ...] [--file F|--stdin] [--public|--private] [--json]
- to notes update <id|slug> [--title ...] [--file F|--stdin] [--public|--private]
- to notes delete <id|slug> [--yes]
- to notes visibility <id|slug> --public|--private
- to notes stats <id|slug> [--json]
- to notes link <id|slug>

Tooling:
- to config get|set|unset|path
- to completion [zsh|bash|fish]
- to update [--check]
- to version
- to doctor

Authentication (magic-link, no subdomain)

to login starts a device-style flow:
CLI calls POST https://textonly.io/api/cli/login/start → gets user_code, device_code, and verification_uri=https://textonly.io/activate
CLI opens the browser to the activation page. After approval, CLI polls POST /api/cli/login/poll until it receives a Bearer token.

Tokens:
Short-lived access sent as a single token (Bearer). Stored in the platform keychain (macOS Keychain, Linux Secret Service, Windows Credential Manager when supported). File fallback uses 0600 perms.
to logout revokes the token server-side and clears local storage.

Configuration

Precedence: flags > env > config file.

Env:
TO_API (default https://textonly.io)
TO_TOKEN (override for CI/PAT)
TO_NO_TELEMETRY=1

Config file: ~/.config/textonly/config.yaml

Proxy: honors HTTPS_PROXY, NO_PROXY.

Shell completions

zshrc: run `to completion zsh` and add to your fpath.

Update and uninstall

update: to update
uninstall: remove the binary and ~/.config/textonly

Build from source (contributors)

Requirements: Go 1.22+, Git, make.

Useful targets:
make build – local build
make test – run unit tests
make lint – static checks
make release – dry-run GoReleaser build

Project layout

cmd/to/ – main
internal/auth/ – login flow, token/keychain
internal/api/ – HTTP client, retries, error handling
internal/notes/ – commands for notes
internal/update/ – self-update/downloader/verification
internal/config/ – env/flags/config resolution
pkg/ui/ – prompts/spinners, JSON output helpers
scripts/install.sh – one-liner installer (published to the website)
.goreleaser.yaml – builds darwin/linux amd64/arm64, produces checksums, SBOM, cosign signatures
.github/workflows/release.yml – tag-driven release pipeline

Releasing (maintainers)

Tag to release: vX.Y.Z

GitHub Actions runs GoReleaser to:
Build static binaries for darwin/linux (arm64/amd64)
Generate checksums and SBOM
Sign artifacts via cosign keyless (OIDC)
Publish release notes
Update/install script channel metadata (optional latest.json)

Homebrew tap: a separate textonly-homebrew-tap repo auto-updated by the pipeline.

Security model

All traffic to https://textonly.io over TLS with HSTS.
Magic-link login requires a first-party browser session to approve.
Tokens stored in platform keychains; fallback file is 0600.
Signed releases (cosign), installer verifies signature+checksum before install.
Minimal telemetry and opt-in via TO_NO_TELEMETRY=1.

Supported platforms

macOS 13+ (arm64, amd64/Intel)
Linux x86_64/arm64 (glibc ≥ 2.31)
Windows not yet supported (planned)

Troubleshooting

to doctor – runs network/auth checks.
“Unauthorized”: run to login again or to logout && to login.
PATH issues: ensure ~/.local/bin or /usr/local/bin precedes other dirs.
Corporate proxies: set HTTPS_PROXY.

License

MIT.
