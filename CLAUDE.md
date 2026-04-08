# CLAUDE.md

**mtls-bridge** — wrap any stdin/stdout process with mutual TLS.

## Build & Test

```bash
make build          # Build for current platform
make build-all      # Cross-compile for all platforms
make test           # Run tests
make coverage       # Run tests with coverage
make lint           # Run go vet
```

## Architecture

- `cmd/mtls-bridge/main.go` — Cobra CLI (serve, connect, provision, update, version)
- `internal/mtls/` — Core library (certs, config, transport, bridge, provisioning, updater)

## Subcommands

- `serve --powershell <path>` — mTLS server wrapping a subprocess
- `connect` — mTLS client bridging stdio
- `provision --generate-psk | --renew | --reset` — Certificate management
- `update` — Self-update from GitHub Releases
- `version` — Print version info

## Critical Rules

1. **Git:** Conventional commits mandatory.
2. **Pure Go:** No CGO. Must cross-compile cleanly.
3. **TLS 1.3 minimum:** All TLS configs enforce `MinVersion: tls.VersionTLS13`.
