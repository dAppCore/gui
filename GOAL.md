<!-- SPDX-License-Identifier: EUPL-1.2 -->

# GOAL — core/gui Wails v3 alpha.91 upgrade

**Owner:** Codex (goal-mode autonomous)
**Working dir:** `/Users/snider/Code/core/gui`
**Branch:** `dev`
**Spawned:** 2026-05-13
**Depends on:** `lthn-desktop` audit passing first (see
`/Users/snider/Code/lthn/desktop/GOAL.md` — the lessons learned about
v0.9.0 compliance shapes feed back into this work as preflight)

---

## Objective

Drop the `replace github.com/wailsapp/wails/v3 => ./stubs/wails`
directive in `go/go.mod`, bump the pinned wails version from
`v3.0.0-alpha.90` to **`v3.0.0-alpha.91`** (already on Wails upstream
`master`), let the real upstream package resolve, and fix every
breaking-API call site in `pkg/window/` / `pkg/systray/` / consumer
packages until `go build ./...` + `go test ./...` are green.

This unblocks the **window-position-remember** port from `core/ide`
into `lthn-desktop` — once `core/gui` builds against the same Wails
version `lthn-desktop` already consumes (`alpha.91`), the
`window.Service` LayoutManager pattern can plug into the per-surface
windows we just landed in `lthn-desktop`.

---

## Success condition

After the work lands:

```bash
cd /Users/snider/Code/core/gui
GOWORK=on go build ./...           # workspace mode (dev)
GOWORK=off go build ./...          # tag mode (CI parity)
go test ./...                      # full test suite green
```

All three must exit 0. Plus the v0.9.0 audit must remain at its
current count or below (this work shouldn't regress it):

```bash
bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh .
```

Record the audit baseline at start (will likely be >0); record again
at end; the delta must be ≤0.

---

## Baseline (run on 2026-05-13 before this goal started)

### go.mod state

```
require github.com/wailsapp/wails/v3 v3.0.0-alpha.90
replace github.com/wailsapp/wails/v3 => ./stubs/wails
```

The `./stubs/wails` directory holds a hand-rolled subset of the
Wails v3 surface (`pkg/application/{window,services,event,...}.go`)
that core/gui currently builds against. This was the load-bearing
hack while upstream Wails v3 was unstable; alpha.91 is stable enough
to drop the stub.

### Consumer of `core/gui` already on alpha.91

- `lthn-desktop` (`/Users/snider/Code/lthn/desktop/go/go.mod` line 10):
  `github.com/wailsapp/wails/v3 v3.0.0-alpha.91`
- That's the proof point — anything alpha.91 changed that core/gui
  uses, lthn-desktop has already had to deal with.

---

## Approach

### Phase 1 — Snapshot

1.1. Record the v0.9.0 audit count baseline:
     ```bash
     bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh . 2>&1 \
       | grep '^verdict:' > /tmp/gui-audit-pre.txt
     ```

1.2. Record what currently uses the stub:
     ```bash
     grep -rln 'wailsapp/wails/v3' --include='*.go' go/ > /tmp/gui-wails-callsites.txt
     ```

1.3. Commit baseline snapshot work to a branch:
     ```bash
     git switch -c codex/gui-wails91-upgrade
     ```

### Phase 2 — Drop the stub, pin alpha.91

2.1. Edit `go/go.mod`:
     - Change `github.com/wailsapp/wails/v3 v3.0.0-alpha.90` →
       `github.com/wailsapp/wails/v3 v3.0.0-alpha.91`
     - Delete the line `replace github.com/wailsapp/wails/v3 => ./stubs/wails`

2.2. Run `cd go && GOWORK=off go mod tidy` to pull real
     alpha.91 + rebalance indirect deps.

2.3. Don't delete `go/stubs/wails/` yet — leave it for now. It can
     come out in Phase 5 after the rest is green.

### Phase 3 — Fix breaking API call sites

3.1. Run `cd /Users/snider/Code/core/gui && go build ./...`.
     The compiler will surface every alpha.90 → alpha.91 breakage.

3.2. Common breaking-change shapes seen in other repos that already
     made this jump (use as a reference, NOT a script — actual
     breakages will vary):

     | alpha.90                                  | alpha.91                                  |
     |-------------------------------------------|-------------------------------------------|
     | `app.NewWindow(opts)`                     | `app.Window.NewWithOptions(opts)`         |
     | `application.SystemTray{}`                | `app.SystemTray.New()`                    |
     | `win.SetPosition(x, y)`                   | (verify still present; rename possible)   |
     | `win.SetSize(w, h)`                       | `win.SetSize(w, h)` (likely unchanged)    |
     | `events.WindowClosed`                     | `events.Common.WindowClosing`             |
     | `application.MenuItem{}`                  | menu builder via `app.Menu.New().Add(...)`|
     | `application.Tray{}.SetMenu(menu)`        | `systray.SetMenu(menu)`                   |
     | `application.WebviewWindowOptions{Hidden}`| (still present; verify field name)        |
     | `app.Run()` signature                     | usually unchanged; verify err return      |

3.3. For each breakage, look at `/Users/snider/Code/lthn/desktop/go/pkg/desktop/*.go`
     for the working alpha.91 pattern. That codebase is already on
     alpha.91 and exercises most of the same surface (windows,
     systray, menus, events).

3.4. After each cluster of fixes:
     ```bash
     cd /Users/snider/Code/core/gui
     go build ./...
     ```
     Commit per-package with conventional prefix
     `fix(gui/<pkg>): alpha.91 — <breakage shape>`.

### Phase 4 — Test pass

4.1. Run the test suite:
     ```bash
     go test ./...
     ```

4.2. Fix any test that relied on stub-specific behaviour. The stub
     was a faithful subset of alpha.90 surface, so most tests should
     pass directly — but tests written against stub-internal types
     (e.g. assertions on `application.MockWindow{}`) need rewriting.

4.3. Commit as `test(gui): alpha.91 — fix stub-dependent tests`.

### Phase 5 — Remove the stub

5.1. `git rm -r go/stubs/wails/`

5.2. Verify nothing still imports from it:
     ```bash
     grep -rln 'stubs/wails' --include='*.go' .
     ```
     Must return empty.

5.3. Run the full build + test cycle one more time. Both must be green.

5.4. Commit as `chore(gui): drop stubs/wails — alpha.91 lands real upstream`.

### Phase 6 — Audit + ship

6.1. Re-run the v0.9.0 audit:
     ```bash
     bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh . 2>&1 \
       | grep '^verdict:'
     ```
     Compare to `/tmp/gui-audit-pre.txt`. Delta must be ≤0 — this work
     shouldn't ADD findings. If it does, identify which new findings
     come from the upgrade (likely err-shape-funcs from refactored
     handlers) and address before merge.

6.2. Push the branch:
     ```bash
     git push -u homelab codex/gui-wails91-upgrade
     git push -u github codex/gui-wails91-upgrade
     ```

6.3. Write `GOAL-STATUS.md` (in this repo root) summarising:
     - Pre/post audit counts
     - Number of files touched
     - Per-package commit list
     - Any judgement calls deferred to Snider

---

## Constraints — DO NOT touch

- **`external/*` submodules** — pinned. Never edit.
- **`go.work`** — leave as is. The stubs/wails workspace use comes
  out of `go.mod`'s replace, not go.work.
- **`pkg/window/LayoutManager`** — this is the load-bearing surface
  for the window-position-remember feature lthn-desktop will consume.
  If it needs refactoring for alpha.91 compatibility, do it minimally;
  do NOT rewrite the storage format or the Action contract
  (`window.restore_layout` / `window.save_layout`).
- **No new packages.** All work happens inside existing files.

---

## Stop conditions — write GOAL-STATUS.md and exit

- alpha.91 introduces a breakage that requires a Snider-class design
  call (e.g. SystemTray contract redesigned, requires API decision).
- v0.9.0 audit count goes UP after the upgrade and codex can't bring
  it back down without changing public API.
- 4 hours elapsed total runtime regardless of progress.

---

## References

- **Wails v3 release notes:** check upstream `wailsapp/wails` on GitHub
  for `v3.0.0-alpha.91` changelog vs alpha.90.
- **lthn-desktop alpha.91 usage:** `/Users/snider/Code/lthn/desktop/go/pkg/desktop/`
  (working reference; copy patterns from there)
- **Audit script:** `/Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh`
- **Sibling clean repos for v0.9.0 reference:** `core/agent`, `core/lint`

---

## Next goal (post-completion)

Once this lands and `lthn-desktop/GOAL.md` is also at 0:

**Window-position-remember port into lthn-desktop.** Concrete shape:

1. In `lthn-desktop`, add `external/gui` as a submodule pointing at
   `https://github.com/dappcore/gui.git` (dev branch).
2. Add `./external/gui/go` to `lthn-desktop/go.work`.
3. Register `gui.Service` on Core via
   `core.WithName("gui", gui.NewService(gui.Options{}))` in
   `/Users/snider/Code/lthn/desktop/go/cmd/lthn/app.go`.
4. In `pkg/desktop/desktop.go`, after `preCreateWindows()`, dispatch
   `c.Action("window.restore_layout").Run(ctx, core.NewOptions(...
   {Name: "default"}))` to pull saved positions onto the pre-created
   windows.
5. Register a SIGTERM hook that calls
   `c.Action("window.save_layout")` with the same Name before quit.
6. The per-plugin windows we open dynamically via `openPluginWindow`
   should ALSO persist position — call save_layout on the
   `WindowClosing` hook and restore_layout when re-opened.

Storage lands at `~/Lethean/conf/layouts.json` (per the no-hidden-
bloat principle — symmetric with how plugins live at
`~/Lethean/conf/plugins/`).
