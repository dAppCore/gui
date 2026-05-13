# GOAL Status - Wails v3 alpha.91 upgrade

Date: 2026-05-13
Branch: codex/gui-wails91-upgrade

## Result

The implementation is complete against the actual Go module in `go/`.
`github.com/wailsapp/wails/v3` now resolves from upstream at
`v3.0.0-alpha.91`, the local `go/stubs/wails` tree has been removed,
and the package build/test gates pass from the module root.

One judgement call remains for Snider: the literal success commands in
`GOAL.md` are not executable from the repository root as written. This
checkout has a workspace root at `/Users/snider/Code/core/gui` and a Go
module root at `/Users/snider/Code/core/gui/go`; additionally, this Go
toolchain rejects `GOWORK=on` because `GOWORK` must be `off`, `auto`, or
an absolute `go.work` path. I treated the module-root equivalents as the
actionable verification lane and recorded the literal failures below.

## Audit

Baseline:

```text
verdict: NON-COMPLIANT - 10 findings
```

Post-upgrade:

```text
verdict: NON-COMPLIANT - 10 findings
```

Delta: 0 findings. The upgrade did not regress the v0.9.0 audit.

## Verification

Passed:

```bash
cd /Users/snider/Code/core/gui/go
GOWORK=/Users/snider/Code/core/gui/go.work go build ./...
GOWORK=off go build ./...
go test ./... -count=1
```

Passed:

```bash
cd /Users/snider/Code/core/gui
rg -n "stubs/wails" -g "*.go" .
```

The `rg` command returned no matches.

Literal `GOAL.md` commands from the repository root:

```bash
GOWORK=on go build ./...
# go: go: invalid GOWORK: not an absolute path

GOWORK=off go build ./...
# pattern ./...: directory prefix . does not contain main module or its selected dependencies

go test ./...
# pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies
```

## Files Touched

Upgrade files touched before this status file: 76.

- 13 modified files in `go/` and `go.work.sum`
- 63 deleted files under `go/stubs/wails`
- 1 status artifact: `GOAL-STATUS.md`

## Commit List

- `70f08697` - resolve upstream Wails alpha.91 deps
- `588e4fc7` - adapt display middleware and bridge routes
- `88486111` - guard the upstream window manager access
- `ec095376` - guard the upstream tray manager access
- `e462899b` - fix stub-dependent tests
- `3007cebb` - drop `go/stubs/wails`

## Deferred Judgements

- Confirm whether future `GOAL.md` success commands should use
  `cd go && GOWORK=/Users/snider/Code/core/gui/go.work ...` for workspace
  mode, or whether the repository should gain a different root-level Go
  invocation convention.
