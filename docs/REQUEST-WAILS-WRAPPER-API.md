# Request: gui wrapper API + rewrite test cascade against wails3 alpha.83

**Owner:** Snider's call (DeepSeek lane).
**Repo:** `dappcore/gui` (this repo) — **dev branch only**, no public squash yet.
**Pairs with:** `dappcore/ide` commit `0fcd4c6` (alpha.83 + AlphaAssets + IPv4 ng host).
**Cladius preflight:** `chore(gui): drop wails3 stub-replace, bump to alpha.83` (commit `22a0d40d`).

---

## Why this exists

Two coupled problems just surfaced:

1. **gui currently leaks wails3 as a public API surface.** Consumers (`ide`, future apps) `import "github.com/wailsapp/wails/v3/pkg/application"` directly to call `application.New(...)` + `app.Run()`, plus `application.WebviewWindowOptions{}` to declare windows. This means every app bringing up a window has to learn wails3 and stay in lockstep with its alpha churn. Snider's intent is the opposite: **gui implements the wails api so consumers don't need to worry about wails.**

2. **The stub-replace block was just removed.** `external/gui/go/go.mod` previously had `replace github.com/wailsapp/wails/v3 => ./stubs/wails`. That stub had ~3000 lines of fake-Wails types (`application.Handler`, `application.Logger{}`, `Menu.Items` field, `WebviewWindow.Title()` getter) that diverged from real wails3. Five test files that reached into those stub-only surfaces via reflection were deleted to unblock the alpha.83 bump:
   - `pkg/menu/wails_test.go`
   - `pkg/systray/wails_test.go`
   - `pkg/window/wails_test.go`
   - `pkg/display/scheme_test.go`
   - `pkg/display/interfaces_test.go`

   Total deleted: ~3,500 lines. They tested gui's wrappers by reflecting into wails3 internals — the right replacement is tests against gui's own wrapper API, which means **the wrapper API has to exist first.**

---

## What to build

### Phase 1 — wrapper API (the load-bearing piece)

Add a single canonical entry point that owns the wails3 application lifecycle. Consumers should be able to launch a window without importing `github.com/wailsapp/wails/v3/...` at all.

**Suggested shape (open to adaptation — pick what reads cleanest):**

```go
// dappco.re/go/gui/pkg/app

type AppOptions struct {
    Name        string
    Description string
    Mac         MacOptions          // gui-owned alias of wails3's
    Services    []core.Service      // takes core.Service registrations, NOT wails3 services
    // Any other top-level wails3 Options fields — re-exposed as gui types.
}

type WindowOptions struct {
    Name      string
    Title     string
    Width     int
    Height    int
    MinWidth  int
    MinHeight int
    // ... etc — re-export of WebviewWindowOptions fields, gui-owned
}

// New constructs the underlying wails3 App + asset server (defaults to
// AlphaAssets so dev mode FRONTEND_DEVSERVER_URL works) + registers the
// caller's Core services as wails3 services.
//
//	gApp := gui.New(gui.AppOptions{Name: "core-ide", Services: []core.Service{...}})
func New(opts AppOptions) *App

// NewWindow opens a window. Idempotent — call after New().
//
//	gApp.NewWindow(gui.WindowOptions{Name: "main", Title: "core/ide", Width: 1180, Height: 780})
func (a *App) NewWindow(opts WindowOptions) Window

// Run blocks on the wails3 event loop. Returns when the app quits or
// ctx is cancelled.
//
//	return gApp.Run(ctx)
func (a *App) Run(ctx context.Context) error

// Quit signals the wails3 event loop to exit cleanly. Safe to call from
// any goroutine, idempotent.
func (a *App) Quit()
```

The point is: **`gui.New` + `gApp.Run(ctx)` should be the entire wails3-touching surface a consumer needs.** Today ide's `pkg/server/gui.go` is the consumer pattern to validate against — after Phase 1 lands, that file should drop its `"github.com/wailsapp/wails/v3/pkg/application"` import entirely.

#### Specific wrapping concerns to handle

- **Asset middleware**: `pkg/display/scheme.go` already migrated to `net/http` types (Handler / ResponseWriter / Request) and exposes `(s *Service) AssetMiddleware() application.Middleware`. The wrapper should let consumers add middleware via gui types so display can ship its `core://` scheme handler without consumers touching wails3.
- **Services**: wails3's `application.NewService(&bridge{})` mechanism is how the frontend gets bindings. The wrapper should accept consumer service handles (e.g. `chatBridge` from ide) and register them with wails3 transparently.
- **Dev-mode default**: `Assets: application.AlphaAssets` is the load-bearing value that activates the dev-server proxy when `FRONTEND_DEVSERVER_URL` is set. Default to it unless the consumer overrides.

### Phase 2 — rewrite the deleted test files against the wrapper API

Re-author the 5 deleted tests as exercises of the **wrapper** API, not wails3 internals. Use exported wrapper getters (which you'll need to add — e.g. `gWindow.Title()`, `gMenu.Items()`) instead of reflection into wails3 unexported fields.

Files to recreate (paths preserved):

| Path | What it tested before |
|------|------------------------|
| `pkg/menu/wails_test.go` | menu construction, item label/accelerator/tooltip/checked/enabled, click handlers |
| `pkg/systray/wails_test.go` | systray menu construction, item state, click |
| `pkg/window/wails_test.go` | window creation, name/title/position/size/visibility/zoom getters + setters, fullscreen toggles, file-drop handler, window events (focus, move, resize, close) |
| `pkg/display/scheme_test.go` | `core://` scheme handler middleware over `http.Handler`, request body parsing, route resolution |
| `pkg/display/interfaces_test.go` | `wailsApp` wrapper around `*application.App` + Logger access |

For each file:
- **Test through gui's wrappers, not wails3 directly.** No `reflect.ValueOf(...).FieldByName(...)` against wails3 unexported fields.
- **Add wrapper getters where missing.** If you need to assert a window's title, gui's `Window` interface should have a `Title() string` method (which already exists in `pkg/window/wails.go` — wraps the locally-stored title since wails3 has no public getter).
- **AX-1**: predictable, descriptive test names. Avoid the `TestX_Good / Bad / Ugly` AX-7 triplet pattern when the three variants have the same body — Snider's flagged it as agent-overworking. Three variants only when they exercise three real branches.

---

## Hard constraints

These come from Cladius's memory of Snider's decisions — break any of them and the diff bounces:

1. **NO `replace` directives in `go.mod`.** Especially no `replace github.com/wailsapp/wails/v3 => ...`. The audit dimension `replace-directives` flags this. We just removed the stub for exactly this reason; do not reintroduce a workaround.

2. **NO stub-Wails revival.** The `stubs/wails/` directory is gone. Don't recreate it. If gui needs CGO-free testing for some surface, use Go interfaces + small testable wrappers, not a parallel API surface.

3. **NO new audit dimensions getting violated.** Run the audit script before declaring done:
   `core/go/tests/cli/v090-upgrade/audit.sh` from the canonical repo. Counters target zero. See `feedback_audit_is_contract.md` in Cladius's memory if uncertain — every counter must be 0, no scope-narrowing.

4. **NO `interface{ Error() string }` anonymous types.** If you need the `resultFailure` shape, declare or import a named type. The existing `pkg/chat/result_failure.go` etc. is the canonical pattern — `type resultFailure = interface { Error() string }` per package.

5. **AX principles** apply throughout — predictable names, usage-example comments (`// Usage example: ...`), path-as-documentation, declarative-over-imperative, core primitives. See `RFC-CORE-008-AGENT-EXPERIENCE.md` if you need the full canon.

6. **UK English** — colour, organisation, centre. Never American.

7. **Workspace mode is the bar.** `go vet ./...` and `go build ./...` from `external/gui/go/` (or with `GOWORK=off` if running standalone). Module proxy / `GOWORK=off` consumers are NOT a session concern — Snider's stated stance.

8. **No bulk perl/sed.** Edit one wrapper site at a time, vet between. Bulk replaces have caused regressions before; they're disallowed except for purely structural single-identifier renames.

---

## Verification (what "done" looks like)

From the **gui repo** (`external/gui/go/`):
```
GOWORK=off go vet ./...      # clean
GOWORK=off go build ./...    # clean
GOWORK=off go test ./...     # all 5 rewritten test files pass
```

From the **ide repo** (`dappcore/ide` workspace mode):
```
cd ide/go && go vet ./...    # clean
cd ide/go && go build ./...  # clean
```

Plus: edit `ide/go/pkg/server/gui.go` to drop `"github.com/wailsapp/wails/v3/pkg/application"` import entirely and use gui's new wrapper instead. The whole point of Phase 1 is that ide doesn't touch wails3 directly anymore — when that file is wails3-import-free and `task dev` still launches a working window, the wrapper API has earned its place.

End-to-end smoke (run from `dappcore/ide`):
```
task dev   →   builds, signs core-ide.dev.app, launches Wails window,
               renders the Angular IDE shell (sidebar / notifications /
               search / version footer all visible — same as today's
               working state).
```

---

## Scope guardrails (don't do these now)

- Don't refactor `pkg/display/scheme.go` further — it's freshly migrated to `net/http`. Leave it.
- Don't bump wails3 past alpha.83 in this PR — that's a separate session.
- Don't touch ide's `frontend/` Angular code. Pure Go-side wrapper work + tests.
- Don't restructure submodules. The `go/` subtree is canonical; `external/<dep>/go` references in `go.work` are the standard.
- Don't open public-forge mirrors. dev branch only on github (origin), pushes squashed to forge.lthn.ai later by Snider.

---

## Hand-off notes

- Branch: work directly on `dev`. Force-push allowed if you rebase, but coordinate with Snider before doing so.
- Commit style: conventional, one cohesive `feat(gui): wails wrapper api` commit + one `test(gui): rewrite wrapper-API tests against alpha.83` commit. Both with `Co-Authored-By: DeepSeek v4 Pro <noreply@deepseek.com>` trailer.
- If a constraint above blocks a clean shape, write the question into the PR description rather than working around — Snider will adjudicate.
