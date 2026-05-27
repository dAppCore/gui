<!-- SPDX-License-Identifier: EUPL-1.2 -->

# Event re-broadcast pattern

core/gui sub-services dispatch native lifecycle / window / environment
events onto Core's action bus as typed `gui<package>.Action*` messages
(e.g. `lifecycle.ActionApplicationStarted`, `window.ActionFilesDropped`,
`environment.ActionThemeChanged`). Apps that want a single uniform
event namespace on the frontend re-broadcast these as custom events
via `events.emit`.

This pattern is **app-specific** — every app picks its own namespace
prefix and payload shape. core/gui doesn't bake the re-broadcast
because:

1. The namespace is the app's brand (`lthn:*`, `core-ide:*`, …).
2. The payload shapes mirror the frontend's contract, which varies by
   app's Lit element conventions.
3. The set of events to bridge is a per-app choice (some apps care
   about file-drop, others don't).

## Canonical example: lthn/desktop

See `lthn/desktop/go/pkg/desktop/sysevents.go` for the working pattern:

```go
func registerSystemEvents(c *core.Core) {
    c.RegisterAction(func(c *core.Core, msg core.Message) core.Result {
        switch event := msg.(type) {
        case lifecycle.ActionApplicationStarted:
            return emitCoreEvent(c, "lthn:app:started", nil)
        case lifecycle.ActionOpenedWithFile:
            return emitCoreEvent(c, "lthn:app:opened-file", event.Path)
        case environment.ActionThemeChanged:
            mode := "light"
            if event.IsDark { mode = "dark" }
            return emitCoreEvent(c, "lthn:theme", mode)
        case window.ActionWindowFocused:
            return emitWindowEvent(c, "focus", event.Name, nil)
        case window.ActionWindowBlurred:
            return emitWindowEvent(c, "blur", event.Name, nil)
        case window.ActionWindowMoved:
            return emitWindowEvent(c, "move", event.Name,
                map[string]any{"x": event.X, "y": event.Y})
        case window.ActionWindowResized:
            return emitWindowEvent(c, "resize", event.Name,
                map[string]any{"width": event.Width, "height": event.Height})
        case window.ActionFilesDropped:
            payload := map[string]any{"files": event.Paths}
            if event.TargetID != "" {
                payload["target"] = map[string]any{"id": event.TargetID}
            }
            return emitWindowEvent(c, "files-dropped", event.Name, payload)
        default:
            return core.Ok(nil)
        }
    })
}

func emitWindowEvent(c *core.Core, verb, window string, payload any) core.Result {
    return emitCoreEvent(c, "lthn:window:"+verb,
        map[string]any{"window": window, "payload": payload})
}

func emitCoreEvent(c *core.Core, name string, data any) core.Result {
    return c.Action("events.emit").Run(core.Background(), core.NewOptions(
        core.Option{Key: "task", Value: events.TaskEmit{Name: name, Data: data}}))
}
```

## How to adapt for your app

1. Pick a namespace prefix (e.g. `myapp:*`).
2. Decide which gui Action* events your frontend needs.
3. Decide the payload shape for each — what fields, what structure.
4. Register a single `RegisterAction` handler that type-switches on
   the gui Action* types and emits via `events.emit`.
5. Subscribe on the frontend via `Events.On("myapp:foo")`.

## Frontend side

The frontend imports `@wailsio/runtime`'s `Events` module:

```ts
import { Events } from "@wailsio/runtime";

Events.On("lthn:window:files-dropped", (e) => {
    const { window, payload } = e.data;
    console.log(`files dropped on ${window}`, payload.files);
});

Events.On("lthn:theme", (e) => {
    document.body.dataset.theme = e.data; // "light" | "dark"
});
```

## When core/gui will help

A future `gui.EventBridge` declarative slice (`[]EventBridge{From: …,
To: …, PayloadField: …}`) could absorb the type-switch boilerplate
when payload shapes are simple field projections. Apps with custom
payload composition (lthn's `{window, payload}` envelope) will still
want the imperative `RegisterAction` shape — that flexibility is the
reason this is doc-only today.

If your app's payload shapes ARE simple field projections, file a
ticket and we'll add the declarative bridge.
