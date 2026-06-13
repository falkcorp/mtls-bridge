<!-- file: .github/copilot-instructions.md -->
<!-- version: 1.0.0 -->
<!-- guid: a3f7c2d1-8b4e-4f9a-b6e3-2d5c1a9f7e84 -->
<!-- last-edited: 2026-06-13 -->

# mtls-bridge — Additional Context

Org-wide coding standards (file headers, language rules, commit format) are at
**https://github.com/falkcorp/.github** and apply automatically to this repo.

For full project context: **CLAUDE.md** at the repo root.

## Project overview

mTLS stdio bridge — wrap any stdin/stdout process with mutual TLS. Language: Go.

## Key directories

| Path | Purpose |
|---|---|
| `cmd/mtls-bridge/` | Cobra CLI entry point (serve, connect, provision, update, version) |
| `internal/mtls/` | Core library: certs, config, transport, bridge, provisioning, updater |

## Critical constraints

- **Pure Go:** No CGO. Must cross-compile cleanly.
- **TLS 1.3 minimum:** All TLS configs enforce `MinVersion: tls.VersionTLS13`.
