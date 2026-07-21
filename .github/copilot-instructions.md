<!-- file: .github/copilot-instructions.md -->
<!-- version: 1.1.0 -->
<!-- guid: a3f7c2d1-8b4e-4f9a-b6e3-2d5c1a9f7e84 -->
<!-- last-edited: 2026-07-21 -->

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
