// SPDX-License-Identifier: EUPL-1.2

package gui

import (
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Binding is the typed wails Service value the consumer constructs via
// gui.Bind[T]. Aliased so consumers can hold a typed slice
// ([]gui.Binding) without importing wails.
type Binding = application.Service

// Bind wraps a domain service pointer as a Wails IPC binding. Consumers
// pass *gui.Bind(svc)* per service in the GuiConfig.Bindings slice; the
// Service's OnStartup forwards the values into the wails App.
//
//	cfg := gui.GuiConfig{
//	    Bindings: []gui.Binding{
//	        gui.Bind(runnerSvc),
//	        gui.Bind(serverSvc),
//	    },
//	}
func Bind[T any](instance *T) Binding {
	return application.NewService(instance)
}

// buildWailsOptions translates a GuiConfig into the wails
// application.Options shape. Consumers stay free of wails imports —
// this file is the single translation point.
func buildWailsOptions(cfg GuiConfig) application.Options {
	name := cfg.Name
	if name == "" {
		name = "core-gui"
	}

	opts := application.Options{
		Name:        name,
		Description: cfg.Description,
		Icon:        cfg.Icon,
		Services:    cfg.Bindings,
		Assets: application.AssetOptions{
			Handler:    cfg.Assets.Handler,
			Middleware: translateMiddleware(cfg.Assets.Middleware),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: cfg.Mac.ApplicationShouldTerminateAfterLastWindowClosed,
			ActivationPolicy: translateActivationPolicy(cfg.Mac.ActivationPolicy),
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: cfg.Windows.DisableQuitOnLastWindowClosed,
			EnabledFeatures:               cfg.Windows.EnabledFeatures,
		},
	}

	if cfg.SingleInstance != nil {
		opts.SingleInstance = translateSingleInstance(*cfg.SingleInstance)
	}
	if cfg.ShouldQuit != nil {
		opts.ShouldQuit = cfg.ShouldQuit
	}
	if cfg.OnShutdown != nil {
		opts.OnShutdown = cfg.OnShutdown
	}
	if cfg.PostShutdown != nil {
		opts.PostShutdown = cfg.PostShutdown
	}
	if cfg.OnPanic != nil {
		opts.PanicHandler = translatePanicHandler(cfg.OnPanic)
	}

	return opts
}

// translateMiddleware converts a core/gui MiddlewareFunc to the wails
// application.Middleware type. Nil input → nil output (wails uses no
// middleware).
func translateMiddleware(mw MiddlewareFunc) application.Middleware {
	if mw == nil {
		return nil
	}
	return func(next http.Handler) http.Handler {
		return mw(next)
	}
}

// translateActivationPolicy maps the core/gui enum to the wails enum.
// Unrecognised values fall back to wails's default (Regular).
func translateActivationPolicy(p ActivationPolicy) application.ActivationPolicy {
	switch p {
	case ActivationPolicyAccessory:
		return application.ActivationPolicyAccessory
	case ActivationPolicyProhibited:
		return application.ActivationPolicyProhibited
	default:
		return application.ActivationPolicyRegular
	}
}

// translateSingleInstance maps the core/gui SingleInstanceOptions to
// the wails type. OnSecondInstanceLaunch is wrapped so the consumer
// receives the core/gui SecondInstanceData shape instead of wails's.
func translateSingleInstance(opts SingleInstanceOptions) *application.SingleInstanceOptions {
	wails := &application.SingleInstanceOptions{
		UniqueID:       opts.UniqueID,
		EncryptionKey:  opts.EncryptionKey,
		AdditionalData: opts.AdditionalData,
	}
	if opts.OnSecondInstanceLaunch != nil {
		consumer := opts.OnSecondInstanceLaunch
		wails.OnSecondInstanceLaunch = func(d application.SecondInstanceData) {
			consumer(SecondInstanceData{
				Args:           d.Args,
				WorkingDir:     d.WorkingDir,
				AdditionalData: d.AdditionalData,
			})
		}
	}
	return wails
}

// translatePanicHandler wraps a core/gui panic callback so it receives
// the core/gui PanicDetails shape instead of wails's.
func translatePanicHandler(consumer func(PanicDetails)) func(*application.PanicDetails) {
	return func(d *application.PanicDetails) {
		if d == nil {
			consumer(PanicDetails{})
			return
		}
		consumer(PanicDetails{
			Error:          d.Error,
			StackTrace:     d.StackTrace,
			FullStackTrace: d.FullStackTrace,
		})
	}
}
