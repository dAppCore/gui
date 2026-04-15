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
	storageBootstrap := map[string]map[string]string{}
	if s.storage != nil {
		storageBootstrap = s.storage.Snapshot(pageURL)
	}
	parts := []string{
		s.injectStoragePolyfills(pageURL, storageBootstrap),
		s.injectBackgroundServiceShims(),
		s.injectElectronShim(),
		s.injectCoreMLShim(),
		s.buildHLCRFComponents(pageURL),
	}
	if appPreloads, err := s.injectAppPreloads(pageURL); err == nil && strings.TrimSpace(appPreloads) != "" {
		parts = append(parts, appPreloads)
	}
	return strings.Join(parts, "\n"), nil
}

func (s *Service) injectStoragePolyfills(pageOrigin string, bootstrap map[string]map[string]string) string {
	return `(function() {
  const __corePageURL = ` + core.JSONMarshalString(pageOrigin) + `;
  const __coreOrigin = ` + core.JSONMarshalString(storageOriginForPageURL(pageOrigin)) + ` || __corePageURL;
  const __coreBootstrapStorage = ` + core.JSONMarshalString(bootstrap) + `;
  const __coreScopes = globalThis.__coreStorageScopes || (globalThis.__coreStorageScopes = {});
  const __scope = __coreScopes[__coreOrigin] || (__coreScopes[__coreOrigin] = { localStorage: {}, sessionStorage: {}, cookies: {}, indexedDB: {}, caches: {}, buckets: {}, opfs: {} });
  const __coreBridge = globalThis.__coreBridge || (globalThis.__coreBridge = {
    invoke(route, payload) {
      if (typeof globalThis.__CORE_GUI_INVOKE__ === 'function') {
        return Promise.resolve(globalThis.__CORE_GUI_INVOKE__(route, payload));
      }
      return Promise.resolve({ route, payload, bridged: false });
    }
  });
  const hydrateBucket = (bucketName, target, source) => {
    if (!source || typeof source !== "object") {
      return;
    }
    const parseMaybeJSON = (value) => {
      if (typeof value !== "string") {
        return value;
      }
      try {
        return JSON.parse(value);
      } catch (_) {
        return value;
      }
    };
    for (const [key, value] of Object.entries(source)) {
      if (bucketName === "cookies") {
        try {
          target[key] = typeof value === "string" ? JSON.parse(value) : value;
        } catch (_) {
          target[key] = { value: String(value ?? "") };
        }
        continue;
      }
      if (bucketName.startsWith("cache:")) {
        const cacheName = bucketName.slice("cache:".length);
        target[cacheName] = target[cacheName] || {};
        target[cacheName][key] = parseMaybeJSON(value);
        continue;
      }
      if (bucketName.startsWith("indexeddb:")) {
        const databaseName = bucketName.slice("indexeddb:".length);
        const [storeName, recordKey] = key.split(":", 2);
        target[databaseName] = target[databaseName] || { stores: {} };
        target[databaseName].stores = target[databaseName].stores || {};
        target[databaseName].stores[storeName] = target[databaseName].stores[storeName] || {};
        target[databaseName].stores[storeName][recordKey] = parseMaybeJSON(value);
        continue;
      }
      if (bucketName.startsWith("storageBucket:")) {
        const bucketNameValue = bucketName.slice("storageBucket:".length);
        target[bucketNameValue] = target[bucketNameValue] || { kv: {}, files: {} };
        target[bucketNameValue].kv = target[bucketNameValue].kv || {};
        target[bucketNameValue].kv[key] = value;
        continue;
      }
      if (bucketName === "opfs") {
        target[key] = { contents: value };
        continue;
      }
      target[key] = value;
    }
  };
  Object.entries(__coreBootstrapStorage || {}).forEach(([bucketName, bucket]) => {
    if (bucketName.startsWith("cache:")) {
      hydrateBucket(bucketName, __scope.caches, bucket);
      return;
    }
    if (bucketName.startsWith("indexeddb:")) {
      hydrateBucket(bucketName, __scope.indexedDB, bucket);
      return;
    }
    if (bucketName.startsWith("storageBucket:")) {
      hydrateBucket(bucketName, __scope.buckets, bucket);
      return;
    }
    if (bucketName === "opfs") {
      hydrateBucket(bucketName, __scope.opfs, bucket);
      return;
    }
    if (!__scope[bucketName]) {
      __scope[bucketName] = {};
    }
    hydrateBucket(bucketName, __scope[bucketName], bucket);
  });
  const persist = (bucket, key, value) => {
    if (bucket === "sessionStorage") {
      return;
    }
    __coreBridge.invoke('display.storage.set', { origin: __coreOrigin, bucket, key, value }).catch(() => undefined);
  };
  const createStorage = (bucketName, bucket) => ({
    getItem(key) { return Object.prototype.hasOwnProperty.call(bucket, key) ? String(bucket[key]) : null; },
    setItem(key, value) { bucket[key] = String(value); persist(bucketName, key, bucket[key]); },
    removeItem(key) { delete bucket[key]; persist(bucketName, key, ''); },
    clear() { Object.keys(bucket).forEach((key) => { delete bucket[key]; persist(bucketName, key, ''); }); },
    key(index) { return Object.keys(bucket)[index] ?? null; },
    get length() { return Object.keys(bucket).length; }
  });
  globalThis.core = globalThis.core || {};
  globalThis.core.storage = globalThis.core.storage || {};
  globalThis.core.storage.local = createStorage('localStorage', __scope.localStorage);
  globalThis.core.storage.session = createStorage('sessionStorage', __scope.sessionStorage);
  globalThis.core.storage.cookies = createStorage('cookies', __scope.cookies);
  const currentLocation = () => globalThis.location || {};
  const cookieValue = (record) => {
    if (record && typeof record === "object" && !Array.isArray(record)) {
      return String(record.value ?? "");
    }
    return String(record ?? "");
  };
  const cookieIsExpired = (record) => {
    if (!record || typeof record !== "object" || Array.isArray(record)) {
      return false;
    }
    const expiresAt = record.expires ? Date.parse(record.expires) : NaN;
    return Number.isFinite(expiresAt) && expiresAt <= Date.now();
  };
  const cookieMatchesPath = (record) => {
    const path = String(currentLocation().pathname || "/");
    const cookiePath = (record && typeof record === "object" && record.path) ? String(record.path) : "/";
    return path === cookiePath || path.startsWith(cookiePath.endsWith("/") ? cookiePath : cookiePath + "/") || cookiePath === "/";
  };
  const cookieMatchesDomain = (record) => {
    const hostname = String(currentLocation().hostname || "");
    const domain = (record && typeof record === "object" && record.domain) ? String(record.domain).replace(/^\./, "") : "";
    if (!domain) {
      return true;
    }
    return hostname === domain || hostname.endsWith("." + domain);
  };
  const cookieMatchesSecure = (record) => {
    if (!record || typeof record !== "object" || Array.isArray(record) || !record.secure) {
      return true;
    }
    const protocol = String(currentLocation().protocol || "");
    return protocol === "https:" || protocol === "wss:" || currentLocation().hostname === "localhost";
  };
  const cookieEntries = () => Object.entries(__scope.cookies)
    .filter(([, record]) => !cookieIsExpired(record) && cookieMatchesPath(record) && cookieMatchesDomain(record) && cookieMatchesSecure(record))
    .map(([name, record]) => name + "=" + cookieValue(record))
    .join("; ");
  const setCookie = (value) => {
    const rawCookie = String(value ?? "");
    const parts = rawCookie.split(";").map((part) => part.trim()).filter(Boolean);
    const [pair, ...attributes] = parts;
    if (!pair) {
      return;
    }
    const separatorIndex = pair.indexOf("=");
    if (separatorIndex < 0) {
      return;
    }
    const name = pair.slice(0, separatorIndex).trim();
    if (!name) {
      return;
    }
    const record = {
      value: pair.slice(separatorIndex + 1).trim(),
      path: "/",
      domain: String(currentLocation().hostname || ""),
      secure: false,
      sameSite: "Lax"
    };
    attributes.forEach((attribute) => {
      const [rawKey, ...rawValue] = attribute.split("=");
      const key = String(rawKey || "").trim().toLowerCase();
      const value = rawValue.join("=").trim();
      switch (key) {
        case "expires":
          if (value) {
            const expiresAt = new Date(value);
            if (!Number.isNaN(expiresAt.getTime())) {
              record.expires = expiresAt.toISOString();
            }
          }
          break;
        case "max-age": {
          const seconds = Number.parseInt(value, 10);
          if (Number.isFinite(seconds)) {
            record.expires = new Date(Date.now() + (seconds * 1000)).toISOString();
          }
          break;
        }
        case "path":
          record.path = value || "/";
          break;
        case "domain":
          record.domain = value.replace(/^\./, "");
          break;
        case "secure":
          record.secure = true;
          break;
        case "samesite":
          record.sameSite = value || "Lax";
          break;
      }
    });
    if (cookieIsExpired(record)) {
      delete __scope.cookies[name];
      persist('cookies', name, '');
      return;
    }
    __scope.cookies[name] = record;
    persist('cookies', name, JSON.stringify(record));
  };
  const asyncResult = (value) => Promise.resolve(value);
  const createIDBRequest = (result) => {
    const request = { result, error: null, onsuccess: null, onerror: null, onupgradeneeded: null };
    queueMicrotask(() => {
      request.onupgradeneeded?.({ target: request });
      request.onsuccess?.({ target: request });
    });
    return request;
  };
  const createObjectStore = (databaseName, database, storeName) => ({
    put(value, key) {
      database.stores[storeName] = database.stores[storeName] || {};
      const resolvedKey = key ?? value?.id ?? crypto.randomUUID?.() ?? String(Date.now());
      database.stores[storeName][resolvedKey] = value;
      persist('indexeddb:' + databaseName, storeName + ':' + resolvedKey, JSON.stringify(value));
      return createIDBRequest(resolvedKey);
    },
    get(key) {
      return createIDBRequest(database.stores?.[storeName]?.[key]);
    },
    getAll() {
      return createIDBRequest(Object.values(database.stores?.[storeName] || {}));
    },
    delete(key) {
      if (database.stores?.[storeName]) {
        delete database.stores[storeName][key];
      }
      persist('indexeddb:' + databaseName, storeName + ':' + key, '');
      return createIDBRequest(undefined);
    },
    clear() {
      database.stores[storeName] = {};
      return createIDBRequest(undefined);
    },
    createIndex() {
      return this;
    }
  });
  const createDB = (name) => {
    const database = __scope.indexedDB[name] || (__scope.indexedDB[name] = { stores: {} });
    return {
      name,
      createObjectStore(storeName) {
        database.stores[storeName] = database.stores[storeName] || {};
        return createObjectStore(name, database, storeName);
      },
      transaction(storeNames) {
        const names = Array.isArray(storeNames) ? storeNames : [storeNames];
        return {
          objectStore(storeName) {
            const target = names.includes(storeName) ? storeName : names[0];
            return createObjectStore(name, database, target);
          }
        };
      },
      close() {}
    };
  };
  const cachesAPI = {
    async open(name) {
      const bucket = __scope.caches[name] || (__scope.caches[name] = {});
      return {
        async match(request) {
          return bucket[typeof request === 'string' ? request : request?.url] ?? undefined;
        },
        async put(request, response) {
          const key = typeof request === 'string' ? request : request?.url;
          bucket[key] = response;
          persist('cache:' + name, key, JSON.stringify(response));
        },
        async delete(request) {
          const key = typeof request === 'string' ? request : request?.url;
          delete bucket[key];
          persist('cache:' + name, key, '');
          return true;
        },
        async keys() {
          return Object.keys(bucket);
        }
      };
    },
    async keys() {
      return Object.keys(__scope.caches);
    }
  };
  const bucketAPI = {
    async open(name) {
      const bucket = __scope.buckets[name] || (__scope.buckets[name] = { kv: {}, files: {} });
      return {
        name,
        persisted() { return asyncResult(true); },
        durability: 'strict',
        async getDirectory() { return bucket.files; },
        storage: createStorage('storageBucket:' + name, bucket.kv)
      };
    }
  };
  const opfsRoot = {
    async getDirectoryHandle(name, options) {
      if (!__scope.opfs[name] || options?.create) {
        __scope.opfs[name] = __scope.opfs[name] || {};
      }
      return __scope.opfs[name];
    },
    async getFileHandle(name, options) {
      if (!__scope.opfs[name] || options?.create) {
        __scope.opfs[name] = __scope.opfs[name] || { contents: '' };
      }
      return {
        async createWritable() {
          return {
            async write(contents) {
              __scope.opfs[name].contents = String(contents);
              persist('opfs', name, __scope.opfs[name].contents);
            },
            async close() {}
          };
        },
        async getFile() {
          return new File([__scope.opfs[name].contents || ''], name);
        }
      };
    }
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
  try {
    if (!globalThis.indexedDB) {
      globalThis.indexedDB = {
        open(name) { return createIDBRequest(createDB(name)); },
        deleteDatabase(name) {
          delete __scope.indexedDB[name];
          return createIDBRequest(undefined);
        }
      };
    }
  } catch (_) {}
  try {
    if (!globalThis.caches) {
      globalThis.caches = cachesAPI;
    }
  } catch (_) {}
  try {
    globalThis.navigator = globalThis.navigator || {};
    globalThis.navigator.storageBuckets = globalThis.navigator.storageBuckets || bucketAPI;
    globalThis.navigator.storage = globalThis.navigator.storage || {};
    if (!globalThis.navigator.storage.getDirectory) {
      globalThis.navigator.storage.getDirectory = () => asyncResult(opfsRoot);
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
      return invokeBridge('gui.browser.open', { url }).then(() => undefined);
    },
    openPath(path) {
      return invokeBridge('gui.browser.openFile', { path }).then(() => "");
    }
  };
  const clipboard = {
    readText() {
      return invokeBridge('gui.clipboard.read', {}).then((value) => {
        if (typeof value === "string") {
          return value;
        }
        return value?.text ?? value?.Text ?? "";
      });
    },
    writeText(text) {
      return invokeBridge('gui.clipboard.write', { text }).then(() => undefined);
    }
  };
  const invokeBridge = (route, payload) => (globalThis.__coreBridge?.invoke?.(route, payload) ?? Promise.resolve({ route, payload }));
  const dialog = {
    showOpenDialog(options) {
      return invokeBridge('gui.dialog.open', options);
    },
    showSaveDialog(options) {
      return invokeBridge('gui.dialog.save', options);
    },
    showMessageBox(options) {
      return invokeBridge('gui.dialog.message', options);
    }
  };
  class CoreNotification {
    constructor(title, options = {}) {
      this.title = title;
      this.options = options;
    }
    show() {
      invokeBridge('gui.notification.send', { title: this.title, ...this.options });
    }
    close() {}
    static requestPermission() {
      return Promise.resolve('granted');
    }
  }
  class BrowserWindow {
    constructor(options = {}) {
      this.options = options;
      this.id = options.id || ('core-window-' + Math.random().toString(36).slice(2));
      invokeBridge('window.open', { name: this.id, options });
    }
    loadURL(url) { return invokeBridge('webview.navigate', { name: this.id, url }); }
    show() { return invokeBridge('window.setVisibility', { name: this.id, visible: true }); }
    hide() { return invokeBridge('window.setVisibility', { name: this.id, visible: false }); }
    close() { return invokeBridge('window.close', { name: this.id }); }
  }
  globalThis.Notification = globalThis.Notification || CoreNotification;
  globalThis.electron = { ipcRenderer, shell, clipboard, dialog, BrowserWindow, Notification: CoreNotification };
  globalThis.require = (name) => name === "electron" ? globalThis.electron : undefined;
})();`
}

func (s *Service) injectBackgroundServiceShims() string {
	return `(function() {
  const invokeBridge = (route, payload) => (globalThis.__coreBridge?.invoke?.(route, payload) ?? Promise.resolve({ route, payload }));
  if (!globalThis.navigator) {
    globalThis.navigator = {};
  }
  globalThis.navigator.serviceWorker = globalThis.navigator.serviceWorker || {
    register(scriptURL, options) {
      return invokeBridge('core.background.serviceWorker.register', { scriptURL, options });
    },
    ready: Promise.resolve({ active: true })
  };
  globalThis.BackgroundFetchManager = globalThis.BackgroundFetchManager || class {
    fetch(id, requests, options) { return invokeBridge('core.background.fetch', { id, requests, options }); }
  };
  globalThis.registration = globalThis.registration || {
    sync: {
      register(tag) { return invokeBridge('core.background.sync', { tag }); }
    },
    periodicSync: {
      register(tag, options) { return invokeBridge('core.background.periodicSync', { tag, options }); }
    },
    pushManager: {
      subscribe(options) { return invokeBridge('core.background.push.subscribe', options); }
    },
    paymentManager: {
      instruments: {
        set(key, details) { return invokeBridge('core.payment.instrument.set', { key, details }); }
      }
    }
  };
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
      if (!response.ok) {
        throw new Error("Core ML request failed: " + response.status + " " + response.statusText);
      }
      const body = await response.text();
      try {
        const parsed = JSON.parse(body);
        const content = parsed?.choices?.[0]?.message?.content;
        if (typeof content === "string") {
          return content;
        }
        if (Array.isArray(content)) {
          return content.map((part) => {
            if (typeof part === "string") {
              return part;
            }
            return part?.text ?? "";
          }).join("");
        }
        if (typeof parsed?.content === "string") {
          return parsed.content;
        }
        if (typeof parsed === "string") {
          return parsed;
        }
        return body;
      } catch (_) {
        return body;
      }
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
    },
    async state() {
      return invokeBridge('display.models.state', {}).then((value) => value);
    },
    async models() {
      const state = await this.state();
      return state.available || state.models || [];
    }
  };
})();`
}

func (s *Service) injectAppPreloads(pageURL string) (string, error) {
	loaded, err := s.loadManifestForOrigin(pageURL)
	if err != nil || loaded == nil {
		return "", err
	}
	scripts := make([]string, 0, len(loaded.Manifest.Preloads))
	for _, preload := range loaded.Manifest.Preloads {
		if preload.Enabled != nil && !*preload.Enabled {
			continue
		}
		if inline := strings.TrimSpace(preload.Inline); inline != "" {
			scripts = append(scripts, inline)
			continue
		}
		if path := strings.TrimSpace(preload.Path); path != "" {
			body, readErr := s.readManifestPreload(loaded.BaseDir, path)
			if readErr != nil {
				return "", readErr
			}
			scripts = append(scripts, string(body))
		}
	}
	return strings.Join(scripts, "\n"), nil
}
