<!-- file: CLAUDE.md -->
<!-- version: 1.2.0 -->
<!-- guid: b8d4e2f3-9c5a-4b7e-a1f2-3e6d8c2b4a7f -->
<!-- last-edited: 2026-07-21 -->

# CLAUDE.md

**mtls-bridge** — wrap any stdin/stdout process with mutual TLS.

## Coding Standards

Org-wide coding standards are in the `.standards/` git submodule (cloned from `https://github.com/falkcorp/.github`).
Always clone with `git clone --recurse-submodules` so these are available.

Key files:
- **File headers (MANDATORY):** `.standards/instructions/file-headers.md`
- **Commit format:** `.standards/instructions/commit-messages.md`

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


## 📝 Changelog & TODO — Use the Fragment System (MANDATORY)

**Do not hand-edit `CHANGELOG.md`, and do not add new tasks straight into the
`TODO.md` inbox.** Both files are assembled from per-change fragments so that
parallel PRs never collide on them.

- **`CHANGELOG.md` is assembled, not hand-edited.** Add a fragment under
  `changelog.d/` (run `scriv create`, or write the Markdown file by hand). The
  fragments are folded into `CHANGELOG.md` at release time by `scriv`, and a CI
  check (`changelog-check.yml`) requires one on each PR. See `changelog.d/README.md`.
- **New `TODO.md` tasks are added via fragments.** Drop a Markdown fragment in
  `todo.d/` (see `todo.d/README.md`) instead of editing the `## 📥 Inbox`
  section. `scripts/assemble_todo.py` folds fragments in daily. This is
  **add-only**: checking a task off or removing it is a normal direct edit of
  `TODO.md`.
