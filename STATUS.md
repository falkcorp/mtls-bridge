# mtls-bridge — Implementation Status

**Date:** 2026-04-07
**Repo:** github.com/jdfalk/mtls-bridge

## Completed

### Task 1-3: Repo Creation + Core Code ✅
- GitHub repo created and public
- go.mod, Makefile, README, CLAUDE.md, LICENSE, .gitignore
- All `internal/mtls/` files extracted from audiobook-organizer:
  - `certs.go` — ECDSA P-256 CA + signed cert generation
  - `config.go` — .mtls/ directory management, PSK, server.json, cert expiry
  - `transport.go` — TLS 1.3 server/client/ephemeral configs
  - `bridge.go` — Bidirectional TCP-to-stdio with half-close
  - `provisioning.go` — PSK exchange protocol
- CLI (`cmd/mtls-bridge/main.go`) with: serve, connect, provision, version
- Import path updated to `github.com/jdfalk/mtls-bridge`
- 18 tests passing

### Task 4: Self-Update Library ✅
- `internal/mtls/updater.go` — GitHub Releases API check, semver compare, self-replace, throttle
- `internal/mtls/updater_test.go` — 5 tests (parse, needsUpdate, check new/current, throttle)
- 23 total tests passing

## Remaining

### Task 5: Wire Update into CLI
- Add `update` subcommand to main.go
- Add auto-update on `serve` startup (check + download + re-exec)
- Add update notification on `connect` startup
- See plan: `docs/superpowers/plans/2026-04-06-mtls-bridge-repo-extraction.md` (Task 5)

### Task 6: Reconnect with Exponential Backoff
- Wrap `connect` command's BridgeStdio call in a reconnect loop
- Backoff: 1s, 2s, 4s, 8s, 16s, 30s cap
- Re-read server.json on each retry (server may have restarted on new port)
- Exit cleanly on stdin EOF (Claude Code closed)
- See plan Task 6

### Task 7: CI Workflow
- `.github/workflows/ci.yml` using `jdfalk/ghcommon` reusable CI
- See plan Task 7

### Task 8: CodeQL Workflow
- `.github/workflows/codeql.yml` — Go security scanning
- Weekly schedule + push/PR triggers
- See plan Task 8

### Task 9: Release Workflow + GoReleaser
- `.github/workflows/release.yml` using `jdfalk/ghcommon` reusable release
- `.goreleaser.yml` — multi-platform builds, checksums, changelog
- See plan Task 9

### Task 10: CODEOWNERS + Final Push
- `.github/CODEOWNERS`
- See plan Task 10

### Task 11: Cleanup audiobook-organizer
- Delete `internal/mtls/` and `cmd/mtls-bridge/` from audiobook-organizer
- Remove Makefile targets, .gitignore entries
- Update `.mcp.json` to use `mtls-bridge` from PATH
- See plan Task 11

### Task 12: Create Initial Release
- Tag v1.0.0, trigger GoReleaser
- Verify binaries + checksums published
- Test auto-update from v0.0.1 → v1.0.0
- See plan Task 12

## References
- **Design spec:** audiobook-organizer `docs/superpowers/specs/2026-04-06-mtls-bridge-repo-extraction-design.md`
- **Implementation plan:** audiobook-organizer `docs/superpowers/plans/2026-04-06-mtls-bridge-repo-extraction.md`
