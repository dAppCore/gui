package display

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"forge.lthn.ai/core/config"
	coreerr "forge.lthn.ai/core/go-log"
	"gopkg.in/yaml.v3"
)

const browserStorageConfigSection = "browser_storage"

// PreloadWindow accepts injected JavaScript before app-specific code runs.
// Use: _ = svc.InjectPreload(windowHandle, "https://app.example.com")
type PreloadWindow interface {
	ExecJS(script string)
}

// ViewManifest describes app-specific preload scripts loaded from .core/view.yaml.
type ViewManifest struct {
	Preloads []ViewPreload `yaml:"preloads" json:"preloads"`
}

// ViewPreload matches an origin and appends a preload script during injection.
type ViewPreload struct {
	Origin string `yaml:"origin" json:"origin"`
	Match  string `yaml:"match" json:"match"`
	Script string `yaml:"script" json:"script"`
}

// BrowserCookie mirrors document.cookie state for one origin.
type BrowserCookie struct {
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Path      string    `json:"path,omitempty"`
	Domain    string    `json:"domain,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Secure    bool      `json:"secure,omitempty"`
}

// OriginStorageState stores browser-like storage namespaces for a single origin.
type OriginStorageState struct {
	LocalStorage   map[string]string            `json:"local_storage,omitempty"`
	SessionStorage map[string]string            `json:"session_storage,omitempty"`
	Cookies        []BrowserCookie              `json:"cookies,omitempty"`
	IndexedDB      map[string]map[string]string `json:"indexed_db,omitempty"`
	CacheStorage   map[string]string            `json:"cache_storage,omitempty"`
	StorageBuckets map[string]map[string]string `json:"storage_buckets,omitempty"`
	OPFS           map[string]string            `json:"opfs,omitempty"`
	UpdatedAt      time.Time                    `json:"updated_at,omitempty"`
}

// BrowserStorageSnapshot persists all origin-scoped storage namespaces.
type BrowserStorageSnapshot struct {
	Origins map[string]OriginStorageState `json:"origins"`
}

// BrowserStoragePolyfillState is the JSON payload injected into preload scripts.
type BrowserStoragePolyfillState struct {
	Origin         string                       `json:"origin"`
	LocalStorage   map[string]string            `json:"localStorage"`
	SessionStorage map[string]string            `json:"sessionStorage"`
	Cookies        []BrowserCookie              `json:"cookies"`
	IndexedDB      map[string]map[string]string `json:"indexedDB"`
	CacheStorage   map[string]string            `json:"cacheStorage"`
	StorageBuckets map[string]map[string]string `json:"storageBuckets"`
	OPFS           map[string]string            `json:"opfs"`
}

// BrowserStorageStore keeps origin-scoped browser data searchable from Go.
type BrowserStorageStore struct {
	mu      sync.RWMutex
	origins map[string]OriginStorageState
}

func NewBrowserStorageStore() *BrowserStorageStore {
	return &BrowserStorageStore{
		origins: make(map[string]OriginStorageState),
	}
}

func (s *BrowserStorageStore) Load(cfg *config.Config) {
	if cfg == nil {
		return
	}

	var snapshot BrowserStorageSnapshot
	if err := cfg.Get(browserStorageConfigSection, &snapshot); err != nil {
		return
	}
	if snapshot.Origins == nil {
		snapshot.Origins = make(map[string]OriginStorageState)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.origins = make(map[string]OriginStorageState, len(snapshot.Origins))
	for origin, state := range snapshot.Origins {
		s.origins[normalizeOrigin(origin)] = cloneOriginStorageState(state)
	}
}

func (s *BrowserStorageStore) Snapshot() BrowserStorageSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshotLocked()
}

func (s *BrowserStorageStore) snapshotLocked() BrowserStorageSnapshot {
	snapshot := BrowserStorageSnapshot{
		Origins: make(map[string]OriginStorageState, len(s.origins)),
	}
	for origin, state := range s.origins {
		snapshot.Origins[origin] = cloneOriginStorageState(state)
	}
	return snapshot
}

func (s *BrowserStorageStore) persist(cfg *config.Config) error {
	if cfg == nil {
		return nil
	}
	if err := cfg.Set(browserStorageConfigSection, s.Snapshot()); err != nil {
		return err
	}
	return cfg.Commit()
}

func (s *BrowserStorageStore) SaveOrigin(origin string, state OriginStorageState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	origin = normalizeOrigin(origin)
	state = cloneOriginStorageState(state)
	if state.UpdatedAt.IsZero() {
		state.UpdatedAt = time.Now().UTC()
	}
	s.origins[origin] = state
}

func (s *BrowserStorageStore) Origin(origin string) (OriginStorageState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.origins[normalizeOrigin(origin)]
	if !ok {
		return OriginStorageState{}, false
	}
	return cloneOriginStorageState(state), true
}

func (s *BrowserStorageStore) PolyfillState(origin string) BrowserStoragePolyfillState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := cloneOriginStorageState(s.origins[normalizeOrigin(origin)])
	return BrowserStoragePolyfillState{
		Origin:         normalizeOrigin(origin),
		LocalStorage:   cloneStorageStringMap(state.LocalStorage),
		SessionStorage: cloneStorageStringMap(state.SessionStorage),
		Cookies:        append([]BrowserCookie(nil), state.Cookies...),
		IndexedDB:      cloneStorageNestedStringMap(state.IndexedDB),
		CacheStorage:   cloneStorageStringMap(state.CacheStorage),
		StorageBuckets: cloneStorageNestedStringMap(state.StorageBuckets),
		OPFS:           cloneStorageStringMap(state.OPFS),
	}
}

func (s *BrowserStorageStore) Search(query string) []StoreSearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	needle := strings.ToLower(strings.TrimSpace(query))
	var results []StoreSearchResult

	for origin, state := range s.origins {
		for key, value := range state.LocalStorage {
			if matchesStorageQuery(needle, key, value) {
				results = append(results, newStorageSearchResult(origin, "localStorage", key, value, state.UpdatedAt))
			}
		}
		for key, value := range state.SessionStorage {
			if matchesStorageQuery(needle, key, value) {
				results = append(results, newStorageSearchResult(origin, "sessionStorage", key, value, state.UpdatedAt))
			}
		}
		for _, cookie := range state.Cookies {
			if matchesStorageQuery(needle, cookie.Name, cookie.Value) {
				results = append(results, newStorageSearchResult(origin, "cookie", cookie.Name, cookie.Value, state.UpdatedAt))
			}
		}
		for database, values := range state.IndexedDB {
			for key, value := range values {
				if matchesStorageQuery(needle, database+"."+key, value) {
					results = append(results, newStorageSearchResult(origin, "indexedDB", database+"."+key, value, state.UpdatedAt))
				}
			}
		}
		for key, value := range state.CacheStorage {
			if matchesStorageQuery(needle, key, value) {
				results = append(results, newStorageSearchResult(origin, "cacheStorage", key, value, state.UpdatedAt))
			}
		}
		for bucket, values := range state.StorageBuckets {
			for key, value := range values {
				if matchesStorageQuery(needle, bucket+"."+key, value) {
					results = append(results, newStorageSearchResult(origin, "storageBucket", bucket+"."+key, value, state.UpdatedAt))
				}
			}
		}
		for path, value := range state.OPFS {
			if matchesStorageQuery(needle, path, value) {
				results = append(results, newStorageSearchResult(origin, "opfs", path, value, state.UpdatedAt))
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	return results
}

func (s *Service) SaveBrowserStorageState(origin string, state OriginStorageState) error {
	if s.browserStorage == nil {
		s.browserStorage = NewBrowserStorageStore()
	}
	s.browserStorage.SaveOrigin(origin, state)
	return s.browserStorage.persist(s.configFile)
}

func (s *Service) BrowserStorageState(origin string) (OriginStorageState, bool) {
	if s.browserStorage == nil {
		return OriginStorageState{}, false
	}
	return s.browserStorage.Origin(origin)
}

// PreloadScript composes storage, background-service, Electron, and local-ML shims for an origin.
// Use: script, _ := svc.PreloadScript("https://app.example.com")
func (s *Service) PreloadScript(origin string) (string, error) {
	if s.browserStorage == nil {
		s.browserStorage = NewBrowserStorageStore()
	}

	storageScript, err := buildStoragePolyfillScript(s.browserStorage.PolyfillState(origin))
	if err != nil {
		return "", err
	}

	scripts := []string{
		storageScript,
		buildCoreMLScript(),
		buildBackgroundServiceScript(),
		buildElectronShimScript(),
	}
	scripts = append(scripts, s.appPreloadScripts(origin)...)
	return strings.Join(scripts, "\n"), nil
}

// InjectPreload installs the composed preload script into a live window handle.
// Use: _ = svc.InjectPreload(windowHandle, "https://app.example.com")
func (s *Service) InjectPreload(target PreloadWindow, origin string) error {
	if target == nil {
		return coreerr.E("display.InjectPreload", "preload target is required", nil)
	}

	script, err := s.PreloadScript(origin)
	if err != nil {
		return err
	}
	target.ExecJS(script)
	return nil
}

// InjectWindowPreload resolves a named window and installs the preload script.
// Use: _ = svc.InjectWindowPreload("main", "https://app.example.com")
func (s *Service) InjectWindowPreload(name, origin string) error {
	ws := s.windowService()
	if ws == nil {
		return coreerr.E("display.InjectWindowPreload", "window service not available", nil)
	}

	target, ok := ws.Manager().Get(name)
	if !ok {
		return coreerr.E("display.InjectWindowPreload", "window not found: "+name, nil)
	}
	return s.InjectPreload(target, origin)
}

func (s *Service) loadViewManifest(path string) {
	s.viewManifest = ViewManifest{}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var manifest ViewManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return
	}
	s.viewManifest = manifest
}

func viewManifestPath(configPath string) string {
	if configPath == "" {
		return filepath.Join(".core", "view.yaml")
	}
	return filepath.Join(filepath.Dir(filepath.Dir(configPath)), "view.yaml")
}

func (s *Service) appPreloadScripts(origin string) []string {
	if len(s.viewManifest.Preloads) == 0 {
		return nil
	}

	normalizedOrigin := normalizeOrigin(origin)
	scripts := make([]string, 0, len(s.viewManifest.Preloads))
	for _, preload := range s.viewManifest.Preloads {
		if !preloadMatchesOrigin(preload, normalizedOrigin) {
			continue
		}
		script := strings.TrimSpace(preload.Script)
		if script == "" {
			continue
		}
		scripts = append(scripts, script)
	}
	return scripts
}

func preloadMatchesOrigin(preload ViewPreload, origin string) bool {
	if strings.TrimSpace(preload.Origin) == "" && strings.TrimSpace(preload.Match) == "" {
		return true
	}

	if preload.Origin != "" && normalizeOrigin(preload.Origin) == origin {
		return true
	}

	match := strings.TrimSpace(preload.Match)
	if match == "" {
		return false
	}
	if match == "*" {
		return true
	}
	if strings.HasSuffix(match, "*") {
		return strings.HasPrefix(origin, strings.TrimSuffix(match, "*"))
	}
	return strings.Contains(origin, match)
}

func buildStoragePolyfillScript(state BrowserStoragePolyfillState) (string, error) {
	encodedState, err := json.Marshal(state)
	if err != nil {
		return "", coreerr.E("display.buildStoragePolyfillScript", "marshal browser storage state", err)
	}

	return fmt.Sprintf(`(function () {
  const initialState = %s;
  const root = globalThis;
  const queueTask = typeof root.queueMicrotask === 'function'
    ? root.queueMicrotask.bind(root)
    : (fn) => Promise.resolve().then(fn);
  const clone = (value) => JSON.parse(JSON.stringify(value));
  const bridge = () => {
    if (root.__coreGUIBridge && typeof root.__coreGUIBridge === 'object') {
      return root.__coreGUIBridge;
    }
    return null;
  };
  const ensureObject = (value) => (value && typeof value === 'object' ? value : {});
  const navigatorObject = root.navigator || {};
  const state = root.__coreStorageState = clone(initialState);
  const normaliseState = () => {
    state.localStorage = ensureObject(state.localStorage);
    state.sessionStorage = ensureObject(state.sessionStorage);
    state.indexedDB = ensureObject(state.indexedDB);
    state.cacheStorage = ensureObject(state.cacheStorage);
    state.storageBuckets = ensureObject(state.storageBuckets);
    state.opfs = ensureObject(state.opfs);
    state.cookies = Array.isArray(state.cookies) ? state.cookies : [];
  };
  const replaceState = (nextState) => {
    for (const key of Object.keys(state)) {
      delete state[key];
    }
    Object.assign(state, clone(nextState));
    normaliseState();
  };
  normaliseState();

  const exportState = () => clone({
    origin: state.origin,
    localStorage: state.localStorage,
    sessionStorage: state.sessionStorage,
    cookies: state.cookies,
    indexedDB: state.indexedDB,
    cacheStorage: state.cacheStorage,
    storageBuckets: state.storageBuckets,
    opfs: state.opfs,
  });

  const sync = () => {
    const currentBridge = bridge();
    if (currentBridge && typeof currentBridge.syncStorage === 'function') {
      Promise.resolve(currentBridge.syncStorage(state.origin, exportState())).catch(() => {});
    }
  };

  const makeStorageArea = (field) => ({
    get length() {
      return Object.keys(state[field]).length;
    },
    clear() {
      state[field] = {};
      sync();
    },
    getItem(key) {
      const value = state[field][String(key)];
      return value === undefined ? null : String(value);
    },
    key(index) {
      return Object.keys(state[field])[index] ?? null;
    },
    removeItem(key) {
      delete state[field][String(key)];
      sync();
    },
    setItem(key, value) {
      state[field][String(key)] = String(value);
      sync();
    },
  });

  const cookieString = () => state.cookies.map((cookie) => {
    if (!cookie || !cookie.name) {
      return '';
    }
    return String(cookie.name) + '=' + String(cookie.value ?? '');
  }).filter(Boolean).join('; ');

  const upsertCookie = (rawCookie) => {
    const segments = String(rawCookie ?? '')
      .split(';')
      .map((segment) => segment.trim())
      .filter(Boolean);
    if (segments.length === 0) {
      return;
    }

    const [nameValue, ...attributes] = segments;
    const separator = nameValue.indexOf('=');
    if (separator === -1) {
      return;
    }

    const cookie = {
      name: nameValue.slice(0, separator).trim(),
      value: nameValue.slice(separator + 1).trim(),
      path: '',
      domain: '',
      secure: false,
      expires_at: '',
    };

    for (const attribute of attributes) {
      const lower = attribute.toLowerCase();
      if (lower === 'secure') {
        cookie.secure = true;
        continue;
      }
      if (lower.startsWith('path=')) {
        cookie.path = attribute.slice(5);
        continue;
      }
      if (lower.startsWith('domain=')) {
        cookie.domain = attribute.slice(7);
        continue;
      }
      if (lower.startsWith('expires=')) {
        cookie.expires_at = attribute.slice(8);
      }
    }

    state.cookies = state.cookies.filter((entry) => entry.name !== cookie.name);
    state.cookies.push(cookie);
    sync();
  };

  try {
    Object.defineProperty(root, 'localStorage', {
      configurable: true,
      value: makeStorageArea('localStorage'),
    });
  } catch (_) {}

  try {
    Object.defineProperty(root, 'sessionStorage', {
      configurable: true,
      value: makeStorageArea('sessionStorage'),
    });
  } catch (_) {}

  if (root.document) {
    try {
      Object.defineProperty(root.document, 'cookie', {
        configurable: true,
        get: cookieString,
        set: upsertCookie,
      });
    } catch (_) {}
  }

  const makeRequest = (executor) => {
    const request = {
      result: undefined,
      error: undefined,
      onsuccess: null,
      onerror: null,
      onupgradeneeded: null,
    };
    queueTask(() => executor(request));
    return request;
  };

  const makeObjectStore = (container) => ({
    get(key) {
      return makeRequest((request) => {
        request.result = container[String(key)];
        if (typeof request.onsuccess === 'function') {
          request.onsuccess({ target: request });
        }
      });
    },
    put(value, key) {
      return makeRequest((request) => {
        container[String(key)] = typeof value === 'string' ? value : JSON.stringify(value);
        sync();
        request.result = String(key);
        if (typeof request.onsuccess === 'function') {
          request.onsuccess({ target: request });
        }
      });
    },
    add(value, key) {
      return this.put(value, key);
    },
    delete(key) {
      return makeRequest((request) => {
        delete container[String(key)];
        sync();
        if (typeof request.onsuccess === 'function') {
          request.onsuccess({ target: request });
        }
      });
    },
    clear() {
      return makeRequest((request) => {
        Object.keys(container).forEach((key) => delete container[key]);
        sync();
        if (typeof request.onsuccess === 'function') {
          request.onsuccess({ target: request });
        }
      });
    },
    getAll() {
      return makeRequest((request) => {
        request.result = Object.values(container);
        if (typeof request.onsuccess === 'function') {
          request.onsuccess({ target: request });
        }
      });
    },
    index() {
      return this;
    },
  });

  if (!root.indexedDB) {
    root.indexedDB = {
      open(databaseName) {
        return makeRequest((request) => {
          const name = String(databaseName);
          state.indexedDB[name] = ensureObject(state.indexedDB[name]);
          request.result = {
            name,
            close() {},
            createObjectStore(storeName) {
              const key = String(storeName);
              state.indexedDB[name][key] = state.indexedDB[name][key] || {};
              sync();
              return makeObjectStore(state.indexedDB[name][key]);
            },
            transaction(storeNames) {
              const storeList = Array.isArray(storeNames) ? storeNames : [storeNames];
              return {
                objectStore(storeName) {
                  const key = String(storeName || storeList[0] || 'default');
                  state.indexedDB[name][key] = state.indexedDB[name][key] || {};
                  return makeObjectStore(state.indexedDB[name][key]);
                },
              };
            },
          };
          if (typeof request.onupgradeneeded === 'function') {
            request.onupgradeneeded({ target: request });
          }
          if (typeof request.onsuccess === 'function') {
            request.onsuccess({ target: request });
          }
        });
      },
    };
  }

  if (!root.caches) {
    root.caches = {
      async delete(cacheName) {
        const prefix = String(cacheName) + '::';
        let deleted = false;
        for (const key of Object.keys(state.cacheStorage)) {
          if (key.startsWith(prefix)) {
            deleted = true;
            delete state.cacheStorage[key];
          }
        }
        if (deleted) {
          sync();
        }
        return deleted;
      },
      async keys() {
        return Array.from(new Set(Object.keys(state.cacheStorage).map((key) => key.split('::')[0])));
      },
      async open(cacheName) {
        const prefix = String(cacheName) + '::';
        return {
          async delete(requestKey) {
            const cacheKey = prefix + String(requestKey);
            const deleted = cacheKey in state.cacheStorage;
            delete state.cacheStorage[cacheKey];
            if (deleted) {
              sync();
            }
            return deleted;
          },
          async keys() {
            return Object.keys(state.cacheStorage)
              .filter((key) => key.startsWith(prefix))
              .map((key) => key.slice(prefix.length));
          },
          async match(requestKey) {
            const value = state.cacheStorage[prefix + String(requestKey)];
            if (value === undefined) {
              return undefined;
            }
            return typeof Response === 'function' ? new Response(value) : value;
          },
          async put(requestKey, responseValue) {
            if (responseValue && typeof responseValue.text === 'function') {
              state.cacheStorage[prefix + String(requestKey)] = await responseValue.text();
            } else {
              state.cacheStorage[prefix + String(requestKey)] = String(responseValue ?? '');
            }
            sync();
          },
        };
      },
    };
  }

  if (!navigatorObject.storageBuckets) {
    navigatorObject.storageBuckets = {
      async open(bucketName) {
        const key = String(bucketName);
        state.storageBuckets[key] = state.storageBuckets[key] || {};
        return {
          name: key,
          kv: {
            delete(itemKey) {
              delete state.storageBuckets[key][String(itemKey)];
              sync();
            },
            get(itemKey) {
              return state.storageBuckets[key][String(itemKey)] ?? null;
            },
            keys() {
              return Object.keys(state.storageBuckets[key]);
            },
            set(itemKey, value) {
              state.storageBuckets[key][String(itemKey)] = String(value);
              sync();
            },
          },
          async persist() {
            return true;
          },
        };
      },
    };
  }

  const makeFileHandle = (path) => ({
    kind: 'file',
    name: path.split('/').pop() || path,
    async createWritable() {
      return {
        async close() {
          return undefined;
        },
        async write(chunk) {
          state.opfs[path] = typeof chunk === 'string' ? chunk : JSON.stringify(chunk);
          sync();
        },
      };
    },
    async getFile() {
      const content = state.opfs[path] ?? '';
      if (typeof File === 'function') {
        return new File([content], path.split('/').pop() || 'file.txt', { type: 'text/plain' });
      }
      return new Blob([content], { type: 'text/plain' });
    },
  });

  const makeDirectoryHandle = (path) => ({
    kind: 'directory',
    name: path.split('/').pop() || '',
    async getDirectoryHandle(name) {
      const child = path ? path + '/' + String(name) : String(name);
      return makeDirectoryHandle(child);
    },
    async getFileHandle(name, options) {
      const child = path ? path + '/' + String(name) : String(name);
      if (!options || !options.create) {
        if (!(child in state.opfs)) {
          throw new Error('File not found: ' + child);
        }
      }
      state.opfs[child] = state.opfs[child] || '';
      sync();
      return makeFileHandle(child);
    },
    async keys() {
      const prefix = path ? path + '/' : '';
      return Object.keys(state.opfs)
        .filter((key) => key.startsWith(prefix))
        .map((key) => key.slice(prefix.length));
    },
  });

  navigatorObject.storage = navigatorObject.storage || {};
  if (!navigatorObject.storage.getDirectory) {
    navigatorObject.storage.getDirectory = async () => makeDirectoryHandle('');
  }

  root.core = root.core || {};
  root.core.storage = root.core.storage || {};
  root.core.storage.origin = state.origin;
  root.core.storage.exportState = exportState;
  root.core.storage.reset = () => {
    replaceState(initialState);
    sync();
    return exportState();
  };
  root.core.storage.sync = sync;
})();`, string(encodedState)), nil
}

func buildCoreMLScript() string {
	return `(function () {
  const root = globalThis;
  root.core = root.core || {};
  root.core.ml = root.core.ml || {};
  root.core.ml.generate = async (prompt) => {
    const text = String(prompt ?? '');
    try {
      const response = await fetch('http://127.0.0.1:8090/v1/chat/completions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          model: 'local',
          stream: false,
          messages: [{ role: 'user', content: text }],
        }),
      });
      if (!response.ok) {
        throw new Error('core.ml bridge returned ' + response.status);
      }
      const payload = await response.json();
      return payload?.choices?.[0]?.message?.content ?? '';
    } catch (_) {
      if (root.__coreGUIBridge && typeof root.__coreGUIBridge.generate === 'function') {
        return Promise.resolve(root.__coreGUIBridge.generate(text));
      }
      return 'Local inference bridge unavailable for prompt: ' + text;
    }
  };
})();`
}

func buildBackgroundServiceScript() string {
	return `(function () {
  const root = globalThis;
  const navigatorObject = root.navigator || {};
  const registrations = new Map();
  const syncTags = new Map();
  const periodicSyncTags = new Map();
  const backgroundFetches = new Map();
  const pushSubscriptions = new Map();

  const invoke = (action, payload = {}) => {
    const bridge = root.__coreGUIBridge;
    if (!bridge || typeof bridge !== 'object') {
      return Promise.reject(new Error('CoreGUI bridge unavailable for ' + action));
    }
    if (typeof bridge.invoke === 'function') {
      return Promise.resolve(bridge.invoke(action, payload));
    }
    if (typeof bridge.call === 'function') {
      return Promise.resolve(bridge.call(action, payload));
    }
    return Promise.reject(new Error('CoreGUI bridge unavailable for ' + action));
  };

  const nativeNotificationCtor = () => {
    if (typeof root.Notification === 'function' && root.Notification !== CoreBrowserNotification) {
      return root.Notification;
    }
    return null;
  };

  const notificationPayload = (title, options = {}) => ({
    title: String(title ?? ''),
    message: String(options.body ?? options.message ?? ''),
    subtitle: typeof options.subtitle === 'string' ? options.subtitle : '',
    icon: typeof options.icon === 'string' ? options.icon : '',
    silent: Boolean(options.silent),
    actions: Array.isArray(options.actions) ? options.actions : [],
  });

  const showNotification = (title, options = {}) => {
    const payload = notificationPayload(title, options);
    return invoke('notification:show', payload).catch(() => {
      const NativeNotification = nativeNotificationCtor();
      if (!NativeNotification) {
        return { id: '', title: payload.title, body: payload.message };
      }
      if (NativeNotification.permission === 'granted') {
        return new NativeNotification(payload.title || 'Notification', {
          body: payload.message,
          icon: payload.icon || undefined,
          silent: payload.silent,
        });
      }
      if (typeof NativeNotification.requestPermission !== 'function') {
        return { id: '', title: payload.title, body: payload.message };
      }
      return Promise.resolve(NativeNotification.requestPermission()).then((permission) => {
        if (permission !== 'granted') {
          return { id: '', title: payload.title, body: payload.message, denied: true };
        }
        return new NativeNotification(payload.title || 'Notification', {
          body: payload.message,
          icon: payload.icon || undefined,
          silent: payload.silent,
        });
      });
    });
  };

  class CoreBrowserNotification {
    constructor(title, options = {}) {
      this.title = String(title ?? '');
      this.options = options && typeof options === 'object' ? { ...options } : {};
      this.onclick = null;
      this.onerror = null;
      this.onshow = null;
      this.onclose = null;
      this._closed = false;
      this._result = showNotification(this.title, this.options)
        .then((result) => {
          if (typeof this.onshow === 'function') {
            this.onshow({ target: result });
          }
          return result;
        })
        .catch((error) => {
          if (typeof this.onerror === 'function') {
            this.onerror(error);
          }
          throw error;
        });
    }

    static get permission() {
      const NativeNotification = nativeNotificationCtor();
      if (NativeNotification && typeof NativeNotification.permission === 'string') {
        return NativeNotification.permission;
      }
      return 'granted';
    }

    static requestPermission(callback) {
      const NativeNotification = nativeNotificationCtor();
      const request = NativeNotification && typeof NativeNotification.requestPermission === 'function'
        ? Promise.resolve(NativeNotification.requestPermission())
        : Promise.resolve('granted');
      return request
        .catch(() => 'granted')
        .then((permission) => {
          if (typeof callback === 'function') {
            callback(permission);
          }
          return permission;
        });
    }

    close() {
      if (this._closed) {
        return;
      }
      this._closed = true;
      if (typeof this.onclose === 'function') {
        this.onclose({ target: this });
      }
    }
  }

  if (typeof root.Notification !== 'function') {
    root.Notification = CoreBrowserNotification;
  }

  const createPushSubscription = (scope, options = {}) => ({
    endpoint: 'core://push/' + encodeURIComponent(scope || '/'),
    options,
    async unsubscribe() {
      pushSubscriptions.delete(String(scope || '/'));
      return true;
    },
    toJSON() {
      return {
        endpoint: this.endpoint,
        options: this.options,
      };
    },
  });

  const createRegistration = (scope, scriptURL) => {
    const registration = {
      scope,
      scriptURL,
      active: {
        scriptURL,
        state: 'activated',
        postMessage() {},
      },
      installing: null,
      waiting: null,
      navigationPreload: {
        async disable() {
          return undefined;
        },
        async enable() {
          return undefined;
        },
        async setHeaderValue() {
          return undefined;
        },
      },
      async showNotification(title, options = {}) {
        return showNotification(title, options);
      },
      async unregister() {
        syncTags.delete(scope);
        periodicSyncTags.delete(scope);
        pushSubscriptions.delete(scope);
        registrations.delete(scope);
        return true;
      },
      async update() {
        return registration;
      },
      sync: {
        async getTags() {
          return Array.from(syncTags.get(scope) || []);
        },
        async register(tag) {
          const tags = syncTags.get(scope) || new Set();
          tags.add(String(tag));
          syncTags.set(scope, tags);
          return undefined;
        },
      },
      periodicSync: {
        async getTags() {
          const tags = periodicSyncTags.get(scope);
          return tags ? Array.from(tags.keys()) : [];
        },
        async register(tag, options = {}) {
          const tags = periodicSyncTags.get(scope) || new Map();
          tags.set(String(tag), options);
          periodicSyncTags.set(scope, tags);
          return undefined;
        },
        async unregister(tag) {
          const tags = periodicSyncTags.get(scope);
          if (tags) {
            tags.delete(String(tag));
          }
          return undefined;
        },
      },
      backgroundFetch: {
        async fetch(id, requests, options = {}) {
          const entries = (Array.isArray(requests) ? requests : [requests]).map((request) => {
            if (typeof request === 'string') {
              return { url: request, init: {} };
            }
            return {
              url: String(request?.url ?? request),
              init: request?.init ?? {},
            };
          });

          const record = {
            id: String(id || 'fetch-' + Date.now()),
            downloaded: 0,
            downloads: [],
            failureReason: '',
            options,
            recordsAvailable: false,
            requests: entries,
            result: 'success',
          };

          try {
            if (typeof root.fetch === 'function') {
              const downloads = await Promise.all(entries.map(async (entry) => {
                const response = await root.fetch(entry.url, entry.init);
                const body = response && typeof response.text === 'function'
                  ? await response.text()
                  : '';
                return {
                  body,
                  ok: Boolean(response && response.ok),
                  status: response ? response.status : 0,
                  url: entry.url,
                };
              }));
              record.downloads = downloads;
              record.downloaded = downloads.reduce((total, download) => total + download.body.length, 0);
              record.recordsAvailable = downloads.length > 0;
            }
          } catch (error) {
            record.result = 'failure';
            record.failureReason = error && error.message ? error.message : 'background fetch failed';
          }

          backgroundFetches.set(scope + '::' + record.id, record);
          return record;
        },
        async get(id) {
          return backgroundFetches.get(scope + '::' + String(id)) || null;
        },
      },
      pushManager: {
        async getSubscription() {
          return pushSubscriptions.get(scope) || null;
        },
        async permissionState() {
          return 'granted';
        },
        async subscribe(options = {}) {
          const subscription = createPushSubscription(scope, options);
          pushSubscriptions.set(scope, subscription);
          return subscription;
        },
      },
      paymentManager: {
        instruments: {
          async clear() {
            return undefined;
          },
          async delete() {
            return true;
          },
          async get() {
            return null;
          },
          async has() {
            return false;
          },
          async keys() {
            return [];
          },
          async set() {
            return undefined;
          },
        },
        async enableDelegations() {
          return undefined;
        },
      },
    };

    return registration;
  };

  const ensureDefaultRegistration = () => {
    if (!registrations.has('/')) {
      registrations.set('/', createRegistration('/', '/coregui-service-worker.js'));
    }
    return registrations.get('/');
  };

  const serviceWorkerContainer = {
    controller: null,
    ready: Promise.resolve().then(() => ensureDefaultRegistration()),
    addEventListener() {},
    dispatchEvent() {
      return true;
    },
    async getRegistration(clientURL) {
      if (typeof clientURL === 'string' && clientURL) {
        return registrations.get(clientURL) || null;
      }
      return ensureDefaultRegistration();
    },
    async getRegistrations() {
      ensureDefaultRegistration();
      return Array.from(registrations.values());
    },
    async register(scriptURL, options = {}) {
      const scope = String(options.scope ?? '/');
      const registration = createRegistration(scope, String(scriptURL ?? '/coregui-service-worker.js'));
      registrations.set(scope, registration);
      serviceWorkerContainer.controller = registration.active;
      return registration;
    },
    removeEventListener() {},
  };

  if (!navigatorObject.serviceWorker) {
    navigatorObject.serviceWorker = serviceWorkerContainer;
  }

  root.core = root.core || {};
  root.core.background = root.core.background || {
    fetch(id, requests, options) {
      return navigatorObject.serviceWorker
        .register('/coregui-service-worker.js')
        .then((registration) => registration.backgroundFetch.fetch(id, requests, options));
    },
    notify(title, options) {
      return showNotification(title, options);
    },
    sync(tag) {
      return navigatorObject.serviceWorker
        .register('/coregui-service-worker.js')
        .then((registration) => registration.sync.register(tag));
    },
  };
})();`
}

func buildElectronShimScript() string {
	return `(function () {
  const root = globalThis;
  const invoke = (action, payload = {}) => {
    const bridge = root.__coreGUIBridge;
    if (!bridge) {
      return Promise.reject(new Error('CoreGUI bridge unavailable for ' + action));
    }
    if (typeof bridge.invoke === 'function') {
      return Promise.resolve(bridge.invoke(action, payload));
    }
    if (typeof bridge.call === 'function') {
      return Promise.resolve(bridge.call(action, payload));
    }
    return Promise.reject(new Error('CoreGUI bridge unavailable for ' + action));
  };

  const subscribe = (channel, handler) => {
    const bridge = root.__coreGUIBridge;
    if (bridge && typeof bridge.on === 'function') {
      return bridge.on(channel, handler);
    }
    return () => {};
  };

  class CoreElectronNotification {
    constructor(options = {}) {
      this.options = options && typeof options === 'object' ? { ...options } : {};
      this.onclick = null;
      this.onshow = null;
      this.onclose = null;
      this.onerror = null;
      this._fallback = null;
    }

    static isSupported() {
      return Boolean(root.__coreGUIBridge) || Boolean(fallbackNotification());
    }

    show() {
      const payload = {
        title: String(this.options.title ?? ''),
        message: String(this.options.body ?? this.options.message ?? ''),
        subtitle: typeof this.options.subtitle === 'string' ? this.options.subtitle : '',
        icon: typeof this.options.icon === 'string' ? this.options.icon : '',
        silent: Boolean(this.options.silent),
        actions: Array.isArray(this.options.actions) ? this.options.actions : [],
      };

      return invoke('notification:show', payload)
        .catch((error) => {
          const NativeNotification = fallbackNotification();
          if (!NativeNotification) {
            throw error;
          }

          const nativeNotification = new NativeNotification(payload.title || 'Notification', {
            body: payload.message,
            icon: payload.icon || undefined,
            silent: payload.silent,
          });

          this._fallback = nativeNotification;
          if (typeof nativeNotification.addEventListener === 'function') {
            nativeNotification.addEventListener('click', () => {
              if (typeof this.onclick === 'function') {
                this.onclick({ target: nativeNotification });
              }
            });
            nativeNotification.addEventListener('close', () => {
              if (typeof this.onclose === 'function') {
                this.onclose({ target: nativeNotification });
              }
            });
          } else {
            nativeNotification.onclick = (event) => {
              if (typeof this.onclick === 'function') {
                this.onclick(event);
              }
            };
          }

          return nativeNotification;
        })
        .then((result) => {
          if (typeof this.onshow === 'function') {
            this.onshow({ target: result });
          }
          return result;
        })
        .catch((error) => {
          if (typeof this.onerror === 'function') {
            this.onerror(error);
          }
          throw error;
        });
    }

    close() {
      if (this._fallback && typeof this._fallback.close === 'function') {
        this._fallback.close();
      }
      if (typeof this.onclose === 'function') {
        this.onclose({ target: this });
      }
    }
  }

  const fallbackNotification = () => {
    if (typeof root.Notification === 'function' && root.Notification !== CoreElectronNotification) {
      return root.Notification;
    }
    return null;
  };

  root.electron = root.electron || {};
  root.electron.Notification = root.electron.Notification || CoreElectronNotification;
  root.electron.clipboard = root.electron.clipboard || {
    readText() {
      return invoke('clipboard:read').then((result) => {
        if (result && typeof result === 'object' && 'text' in result) {
          return String(result.text ?? '');
        }
        return String(result ?? '');
      });
    },
    writeText(text) {
      return invoke('clipboard:write', { text: String(text ?? '') });
    },
  };
  root.electron.dialog = root.electron.dialog || {
    showMessageBox(options) {
      return invoke('dialog:message', options ?? {});
    },
    showOpenDialog(options) {
      return invoke('dialog:open-file', options ?? {});
    },
    showSaveDialog(options) {
      return invoke('dialog:save-file', options ?? {});
    },
  };
  root.electron.ipcRenderer = root.electron.ipcRenderer || {
    invoke(channel, payload) {
      return invoke('core.ipc.query', { channel, payload });
    },
    on(channel, handler) {
      return subscribe(channel, handler);
    },
    send(channel, payload) {
      return invoke('core.ipc.action', { channel, payload });
    },
  };
  root.electron.shell = root.electron.shell || {
    openExternal(target) {
      return invoke('browser:open-url', { url: String(target ?? '') });
    },
    openPath(target) {
      return invoke('browser:open-file', { path: String(target ?? '') });
    },
  };

  root.require = root.require || ((name) => {
    if (name === 'electron') {
      return root.electron;
    }
    throw new Error('Unsupported module: ' + String(name));
  });
})();`
}

func deriveOrigin(rawURL string) string {
	normalized := strings.TrimSpace(rawURL)
	if normalized == "" {
		return "app://local"
	}

	parsed, err := url.Parse(normalized)
	if err != nil {
		return "app://local"
	}

	if parsed.Scheme == "" && parsed.Host == "" {
		return "app://local"
	}

	origin := parsed.Scheme + "://"
	if parsed.Host != "" {
		origin += parsed.Host
		return origin
	}

	return origin + strings.Trim(parsed.Path, "/")
}

func normalizeOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return "app://local"
	}
	return deriveOrigin(origin)
}

func matchesStorageQuery(query string, key string, value string) bool {
	if query == "" {
		return true
	}
	needle := strings.ToLower(query)
	return strings.Contains(strings.ToLower(key), needle) || strings.Contains(strings.ToLower(value), needle)
}

func newStorageSearchResult(origin string, area string, key string, value string, updatedAt time.Time) StoreSearchResult {
	return StoreSearchResult{
		Origin:      origin,
		StorageArea: area,
		Key:         key,
		Snippet:     trimSnippet(value),
		UpdatedAt:   updatedAt,
	}
}

func cloneOriginStorageState(state OriginStorageState) OriginStorageState {
	return OriginStorageState{
		LocalStorage:   cloneStorageStringMap(state.LocalStorage),
		SessionStorage: cloneStorageStringMap(state.SessionStorage),
		Cookies:        append([]BrowserCookie(nil), state.Cookies...),
		IndexedDB:      cloneStorageNestedStringMap(state.IndexedDB),
		CacheStorage:   cloneStorageStringMap(state.CacheStorage),
		StorageBuckets: cloneStorageNestedStringMap(state.StorageBuckets),
		OPFS:           cloneStorageStringMap(state.OPFS),
		UpdatedAt:      state.UpdatedAt,
	}
}

func cloneStorageStringMap(source map[string]string) map[string]string {
	if len(source) == 0 {
		return map[string]string{}
	}
	cloned := make(map[string]string, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}

func cloneStorageNestedStringMap(source map[string]map[string]string) map[string]map[string]string {
	if len(source) == 0 {
		return map[string]map[string]string{}
	}
	cloned := make(map[string]map[string]string, len(source))
	for key, values := range source {
		cloned[key] = cloneStorageStringMap(values)
	}
	return cloned
}
