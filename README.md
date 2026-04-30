<!-- SPDX-License-Identifier: EUPL-1.2 -->

# Core GUI

> Wails-based GUI runtime — webview bridge, stubs, package primitives

[![CI](https://github.com/dappcore/gui/actions/workflows/ci.yml/badge.svg?branch=dev)](https://github.com/dappcore/gui/actions/workflows/ci.yml)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=dappcore_gui&metric=alert_status)](https://sonarcloud.io/dashboard?id=dappcore_gui)
[![Coverage](https://codecov.io/gh/dappcore/gui/branch/dev/graph/badge.svg)](https://codecov.io/gh/dappcore/gui)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=dappcore_gui&metric=security_rating)](https://sonarcloud.io/dashboard?id=dappcore_gui)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=dappcore_gui&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=dappcore_gui)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=dappcore_gui&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=dappcore_gui)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=dappcore_gui&metric=code_smells)](https://sonarcloud.io/dashboard?id=dappcore_gui)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=dappcore_gui&metric=ncloc)](https://sonarcloud.io/dashboard?id=dappcore_gui)
[![Go Reference](https://pkg.go.dev/badge/dappco.re/go/gui.svg)](https://pkg.go.dev/dappco.re/go/gui)
[![License: EUPL-1.2](https://img.shields.io/badge/License-EUPL--1.2-blue.svg)](https://eupl.eu/1.2/en/)


Core GUI is the Go backend surface for a desktop GUI built around Core services
and Wails-style application primitives. It provides service packages for browser
control, chat, clipboard, container/TIM lifecycle, context menus, dialogs,
display orchestration, dock integration, environment state, events, keybindings,
menus, notifications, peer-to-peer messaging, preload injection, screens,
system tray integration, webview automation, and window management.

The module is `dappco.re/go/gui`. It depends on `dappco.re/go` for Core
primitives and wraps Wails through `stubs/wails` so ordinary Go tests can run
without native desktop bindings. The packages expose small registration
functions and service objects that can be mounted into a Core runtime.

Typical local verification:

```sh
GOWORK=off go mod tidy
GOWORK=off go vet ./...
GOWORK=off go test -count=1 ./...
bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh .
```

Use `pkg/display` when wiring the full GUI experience. Use individual packages
when testing or embedding a narrow capability, such as `pkg/window` for window
state/layout, `pkg/preload` for trusted preload policy, or `pkg/mcp` for MCP tool
surfaces over the GUI services.
