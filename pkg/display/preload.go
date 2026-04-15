package display

import (
	"strings"

	core "dappco.re/go/core"
)

type PreloadTarget interface {
	ExecJS(string)
}

// InjectPreload injects the page preload script into a live webview.
// Use: _ = display.InjectPreload(webview, "https://example.com")
func (s *Service) InjectPreload(webview PreloadTarget, origin string) error {
	script, err := s.BuildPreloadScript(origin)
	if err != nil {
		return err
	}
	if strings.TrimSpace(script) == "" {
		return nil
	}
	webview.ExecJS(script)
	return nil
}

// BuildPreloadScript returns the JavaScript bootstrap that CoreGUI injects
// before page code runs.
// Use: script, _ := display.BuildPreloadScript("https://example.com")
func (s *Service) BuildPreloadScript(pageURL string) (string, error) {
	parts := []string{
		s.injectStoragePolyfills(pageURL),
		s.injectElectronShim(),
		s.injectCoreMLShim(),
		s.injectAppPreloads(pageURL),
	}
	return strings.Join(parts, "\n"), nil
}

func (s *Service) injectStoragePolyfills(pageOrigin string) string {
	return `(function() {
  const __coreOrigin = ` + core.JSONMarshalString(pageOrigin) + `;
  const __coreScopes = globalThis.__coreStorageScopes || (globalThis.__coreStorageScopes = {});
  const __scope = __coreScopes[__coreOrigin] || (__coreScopes[__coreOrigin] = { local: {}, session: {}, cookies: {} });
  const createStorage = (bucket) => ({
    getItem(key) { return Object.prototype.hasOwnProperty.call(bucket, key) ? String(bucket[key]) : null; },
    setItem(key, value) { bucket[key] = String(value); },
    removeItem(key) { delete bucket[key]; },
    clear() { Object.keys(bucket).forEach((key) => delete bucket[key]); },
    key(index) { return Object.keys(bucket)[index] ?? null; },
    get length() { return Object.keys(bucket).length; }
  });
  globalThis.core = globalThis.core || {};
  globalThis.core.storage = globalThis.core.storage || {};
  globalThis.core.storage.local = createStorage(__scope.local);
  globalThis.core.storage.session = createStorage(__scope.session);
  globalThis.core.storage.cookies = createStorage(__scope.cookies);
  const cookieEntries = () => Object.entries(__scope.cookies).map(([name, value]) => name + "=" + value).join("; ");
  const setCookie = (value) => {
    const rawCookie = String(value ?? "");
    const [pair] = rawCookie.split(";", 1);
    const separatorIndex = pair.indexOf("=");
    if (separatorIndex < 0) {
      return;
    }
    const name = pair.slice(0, separatorIndex).trim();
    if (!name) {
      return;
    }
    __scope.cookies[name] = pair.slice(separatorIndex + 1).trim();
  };
  try {
    if (!globalThis.localStorage) {
      Object.defineProperty(globalThis, 'localStorage', {
        configurable: true,
        enumerable: true,
        get() { return globalThis.core.storage.local; }
      });
    }
  } catch (_) {}
  try {
    if (!globalThis.sessionStorage) {
      Object.defineProperty(globalThis, 'sessionStorage', {
        configurable: true,
        enumerable: true,
        get() { return globalThis.core.storage.session; }
      });
    }
  } catch (_) {}
  try {
    if (typeof Document !== 'undefined') {
      Object.defineProperty(Document.prototype, 'cookie', {
        configurable: true,
        enumerable: true,
        get() { return cookieEntries(); },
        set(value) { setCookie(value); }
      });
    }
  } catch (_) {}
})();`
}

func (s *Service) injectElectronShim() string {
	return `(function() {
  if (globalThis.electron) {
    return;
  }
  const listeners = new Map();
  const toEventName = (channel) => "__core_electron__:" + channel;
  const ipcRenderer = {
    send(channel, ...args) {
      globalThis.dispatchEvent(new CustomEvent(toEventName(channel), { detail: args }));
    },
    invoke(channel, ...args) {
      globalThis.dispatchEvent(new CustomEvent(toEventName(channel), { detail: args }));
      return Promise.resolve({ channel, args });
    },
    on(channel, listener) {
      const handler = (event) => listener(event, ...(event.detail || []));
      listeners.set(listener, handler);
      globalThis.addEventListener(toEventName(channel), handler);
      return () => ipcRenderer.removeListener(channel, listener);
    },
    once(channel, listener) {
      const off = ipcRenderer.on(channel, (event, ...args) => {
        off();
        listener(event, ...args);
      });
      return off;
    },
    removeListener(channel, listener) {
      const handler = listeners.get(listener);
      if (handler) {
        globalThis.removeEventListener(toEventName(channel), handler);
        listeners.delete(listener);
      }
    }
  };
  const shell = {
    openExternal(url) {
      globalThis.open?.(url, "_blank", "noopener,noreferrer");
      return Promise.resolve();
    },
    openPath(path) {
      globalThis.location.href = path;
      return Promise.resolve("");
    }
  };
  const clipboard = {
    readText() {
      return globalThis.navigator?.clipboard?.readText?.() ?? Promise.resolve("");
    },
    writeText(text) {
      return globalThis.navigator?.clipboard?.writeText?.(text) ?? Promise.resolve();
    }
  };
  globalThis.electron = { ipcRenderer, shell, clipboard };
  globalThis.require = (name) => name === "electron" ? globalThis.electron : undefined;
})();`
}

func (s *Service) injectCoreMLShim() string {
	return `(function() {
  const __coreMLApiURL = ` + core.JSONMarshalString(strings.TrimRight(core.Env("CORE_ML_API_URL"), "/")) + ` || "http://localhost:8090";
  globalThis.core = globalThis.core || {};
  globalThis.core.ml = globalThis.core.ml || {
    async generate(input) {
      const payload = typeof input === "string"
        ? { messages: [{ role: "user", content: input }], stream: false }
        : { ...input, stream: false };
      const response = await fetch((__coreMLApiURL || "http://localhost:8090") + "/v1/chat/completions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
      return response.json();
    },
    async stream(input) {
      const payload = typeof input === "string"
        ? { messages: [{ role: "user", content: input }], stream: true }
        : { ...input, stream: true };
      return fetch((__coreMLApiURL || "http://localhost:8090") + "/v1/chat/completions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
    }
  };
})();`
}

func (s *Service) injectAppPreloads(_ string) string {
	return ""
}
