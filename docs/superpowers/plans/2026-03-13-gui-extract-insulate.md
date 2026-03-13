# GUI Extract & Insulate Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extract clipboard, dialog, notification, environment, and screen from pkg/display into independent core.Service packages with Platform interface insulation and full IPC coverage.

**Architecture:** Each package follows the three-layer pattern (IPC Bus → Service → Platform Interface) established by window/systray/menu. The display orchestrator sheds ~420 lines and gains HandleIPCEvents + WS bridge cases for the new Action types.

**Tech Stack:** Go, core/go DI framework, Wails v3 (behind Platform interfaces), gorilla/websocket (WSEventManager)

**Spec:** `docs/superpowers/specs/2026-03-13-gui-extract-insulate-design.md`

---

## File Structure

### New files (per package, 5 packages × 4 files = 20 files)

| Package | File | Responsibility |
|---------|------|---------------|
| `pkg/clipboard/` | `platform.go` | Platform interface (2 methods) |
| | `messages.go` | IPC message types (QueryText, TaskSetText, TaskClear) |
| | `service.go` | Service struct, Register(), OnStartup(), handlers |
| | `service_test.go` | Tests with mock platform |
| `pkg/dialog/` | `platform.go` | Platform interface (4 methods) + option types |
| | `messages.go` | IPC message types (4 Tasks) |
| | `service.go` | Service struct, Register(), OnStartup(), handlers |
| | `service_test.go` | Tests with mock platform |
| `pkg/notification/` | `platform.go` | Platform interface (3 methods) + option types |
| | `messages.go` | IPC message types (QueryPermission, TaskSend, TaskRequestPermission) |
| | `service.go` | Service struct, Register(), OnStartup(), handlers, fallback logic |
| | `service_test.go` | Tests with mock platform + fallback test |
| `pkg/environment/` | `platform.go` | Platform interface (5 methods) + own types |
| | `messages.go` | IPC message types (3 Queries, 1 Task, 1 Action) |
| | `service.go` | Service struct, Register(), OnStartup(), theme callback, handlers |
| | `service_test.go` | Tests with mock platform + theme change test |
| `pkg/screen/` | `platform.go` | Platform interface (2 methods) + own types |
| | `messages.go` | IPC message types (4 Queries, 1 Action) |
| | `service.go` | Service struct, Register(), OnStartup(), computed queries |
| | `service_test.go` | Tests with mock platform + computed query tests |

### Modified files

| File | Change |
|------|--------|
| `pkg/display/display.go` | Remove screen methods (~120 lines), remove ShowEnvironmentDialog, remove ScreenInfo/WorkArea types |
| `pkg/display/events.go` | Change `NewWSEventManager` signature to drop `EventSource`, remove `SetupWindowEventListeners()` |
| `pkg/display/interfaces.go` | Remove `wailsEventSource`, `wailsDialogManager` (moved to new packages) |
| `pkg/display/types.go` | Remove `EventSource` interface |
| `pkg/display/clipboard.go` | DELETE (moved to pkg/clipboard) |
| `pkg/display/dialog.go` | DELETE (moved to pkg/dialog) |
| `pkg/display/notification.go` | DELETE (moved to pkg/notification) |
| `pkg/display/theme.go` | DELETE (moved to pkg/environment) |

---

## Chunk 1: Clipboard + Dialog

### Task 1: Create pkg/clipboard

**Files:**
- Create: `pkg/clipboard/platform.go`
- Create: `pkg/clipboard/messages.go`
- Create: `pkg/clipboard/service.go`
- Create: `pkg/clipboard/service_test.go`
- Delete: `pkg/display/clipboard.go` (after Task 7)

- [ ] **Step 1: Create platform.go**

```go
// pkg/clipboard/platform.go
package clipboard

// Platform abstracts the system clipboard backend.
type Platform interface {
	Text() (string, bool)
	SetText(text string) bool
}

// ClipboardContent is the result of QueryText.
type ClipboardContent struct {
	Text       string `json:"text"`
	HasContent bool   `json:"hasContent"`
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/clipboard/messages.go
package clipboard

// QueryText reads the clipboard. Result: ClipboardContent
type QueryText struct{}

// TaskSetText writes text to the clipboard. Result: bool (success)
type TaskSetText struct{ Text string }

// TaskClear clears the clipboard. Result: bool (success)
type TaskClear struct{}
```

- [ ] **Step 3: Write failing test**

```go
// pkg/clipboard/service_test.go
package clipboard

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	text string
	ok   bool
}

func (m *mockPlatform) Text() (string, bool) { return m.text, m.ok }
func (m *mockPlatform) SetText(text string) bool {
	m.text = text
	m.ok = text != ""
	return true
}

func newTestService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(&mockPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "clipboard")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestService(t)
	assert.NotNil(t, svc)
}

func TestQueryText_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryText{})
	require.NoError(t, err)
	assert.True(t, handled)
	content := result.(ClipboardContent)
	assert.Equal(t, "hello", content.Text)
	assert.True(t, content.HasContent)
}

func TestQueryText_Bad(t *testing.T) {
	// No clipboard service registered
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryText{})
	assert.False(t, handled)
}

func TestTaskSetText_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskSetText{Text: "world"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, true, result)

	// Verify via query
	r, _, _ := c.QUERY(QueryText{})
	assert.Equal(t, "world", r.(ClipboardContent).Text)
}

func TestTaskClear_Good(t *testing.T) {
	_, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskClear{})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify empty
	r, _, _ := c.QUERY(QueryText{})
	assert.Equal(t, "", r.(ClipboardContent).Text)
	assert.False(t, r.(ClipboardContent).HasContent)
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/clipboard/ -v`
Expected: FAIL — `Service` type not defined

- [ ] **Step 5: Create service.go**

```go
// pkg/clipboard/service.go
package clipboard

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the clipboard service.
type Options struct{}

// Service is a core.Service managing clipboard operations via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryText:
		text, ok := s.platform.Text()
		return ClipboardContent{Text: text, HasContent: ok && text != ""}, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetText:
		return s.platform.SetText(t.Text), true, nil
	case TaskClear:
		return s.platform.SetText(""), true, nil
	default:
		return nil, false, nil
	}
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/clipboard/ -v`
Expected: PASS (5 tests)

- [ ] **Step 7: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/clipboard/
git commit -m "feat(clipboard): add clipboard core.Service with Platform interface and IPC"
```

---

### Task 2: Create pkg/dialog

**Files:**
- Create: `pkg/dialog/platform.go`
- Create: `pkg/dialog/messages.go`
- Create: `pkg/dialog/service.go`
- Create: `pkg/dialog/service_test.go`
- Delete: `pkg/display/dialog.go` (after Task 7)

- [ ] **Step 1: Create platform.go with types**

```go
// pkg/dialog/platform.go
package dialog

// Platform abstracts the native dialog backend.
type Platform interface {
	OpenFile(opts OpenFileOptions) ([]string, error)
	SaveFile(opts SaveFileOptions) (string, error)
	OpenDirectory(opts OpenDirectoryOptions) (string, error)
	MessageDialog(opts MessageDialogOptions) (string, error)
}

// DialogType represents the type of message dialog.
type DialogType int

const (
	DialogInfo DialogType = iota
	DialogWarning
	DialogError
	DialogQuestion
)

// OpenFileOptions contains options for the open file dialog.
type OpenFileOptions struct {
	Title         string       `json:"title,omitempty"`
	Directory     string       `json:"directory,omitempty"`
	Filename      string       `json:"filename,omitempty"`
	Filters       []FileFilter `json:"filters,omitempty"`
	AllowMultiple bool         `json:"allowMultiple,omitempty"`
}

// SaveFileOptions contains options for the save file dialog.
type SaveFileOptions struct {
	Title     string       `json:"title,omitempty"`
	Directory string       `json:"directory,omitempty"`
	Filename  string       `json:"filename,omitempty"`
	Filters   []FileFilter `json:"filters,omitempty"`
}

// OpenDirectoryOptions contains options for the directory picker.
type OpenDirectoryOptions struct {
	Title         string `json:"title,omitempty"`
	Directory     string `json:"directory,omitempty"`
	AllowMultiple bool   `json:"allowMultiple,omitempty"`
}

// MessageDialogOptions contains options for a message dialog.
type MessageDialogOptions struct {
	Type    DialogType `json:"type"`
	Title   string     `json:"title"`
	Message string     `json:"message"`
	Buttons []string   `json:"buttons,omitempty"`
}

// FileFilter represents a file type filter for dialogs.
type FileFilter struct {
	DisplayName string   `json:"displayName"`
	Pattern     string   `json:"pattern"`
	Extensions  []string `json:"extensions,omitempty"`
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/dialog/messages.go
package dialog

// TaskOpenFile shows an open file dialog. Result: []string (paths)
type TaskOpenFile struct{ Opts OpenFileOptions }

// TaskSaveFile shows a save file dialog. Result: string (path)
type TaskSaveFile struct{ Opts SaveFileOptions }

// TaskOpenDirectory shows a directory picker. Result: string (path)
type TaskOpenDirectory struct{ Opts OpenDirectoryOptions }

// TaskMessageDialog shows a message dialog. Result: string (button clicked)
type TaskMessageDialog struct{ Opts MessageDialogOptions }
```

- [ ] **Step 3: Write failing test**

```go
// pkg/dialog/service_test.go
package dialog

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	openFilePaths  []string
	saveFilePath   string
	openDirPath    string
	messageButton  string
	openFileErr    error
	saveFileErr    error
	openDirErr     error
	messageErr     error
	lastOpenOpts   OpenFileOptions
	lastSaveOpts   SaveFileOptions
	lastDirOpts    OpenDirectoryOptions
	lastMsgOpts    MessageDialogOptions
}

func (m *mockPlatform) OpenFile(opts OpenFileOptions) ([]string, error) {
	m.lastOpenOpts = opts
	return m.openFilePaths, m.openFileErr
}
func (m *mockPlatform) SaveFile(opts SaveFileOptions) (string, error) {
	m.lastSaveOpts = opts
	return m.saveFilePath, m.saveFileErr
}
func (m *mockPlatform) OpenDirectory(opts OpenDirectoryOptions) (string, error) {
	m.lastDirOpts = opts
	return m.openDirPath, m.openDirErr
}
func (m *mockPlatform) MessageDialog(opts MessageDialogOptions) (string, error) {
	m.lastMsgOpts = opts
	return m.messageButton, m.messageErr
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		openFilePaths: []string{"/tmp/file.txt"},
		saveFilePath:  "/tmp/save.txt",
		openDirPath:   "/tmp/dir",
		messageButton: "OK",
	}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	return mock, c
}

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "dialog")
	assert.NotNil(t, svc)
}

func TestTaskOpenFile_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/a.txt", "/b.txt"}

	result, handled, err := c.PERFORM(TaskOpenFile{
		Opts: OpenFileOptions{Title: "Pick", AllowMultiple: true},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	paths := result.([]string)
	assert.Equal(t, []string{"/a.txt", "/b.txt"}, paths)
	assert.Equal(t, "Pick", mock.lastOpenOpts.Title)
	assert.True(t, mock.lastOpenOpts.AllowMultiple)
}

func TestTaskSaveFile_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskSaveFile{
		Opts: SaveFileOptions{Filename: "out.txt"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/save.txt", result)
}

func TestTaskOpenDirectory_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskOpenDirectory{
		Opts: OpenDirectoryOptions{Title: "Pick Dir"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/dir", result)
}

func TestTaskMessageDialog_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Yes"

	result, handled, err := c.PERFORM(TaskMessageDialog{
		Opts: MessageDialogOptions{
			Type: DialogQuestion, Title: "Confirm",
			Message: "Sure?", Buttons: []string{"Yes", "No"},
		},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "Yes", result)
	assert.Equal(t, DialogQuestion, mock.lastMsgOpts.Type)
}

func TestTaskOpenFile_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenFile{})
	assert.False(t, handled)
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/dialog/ -v`
Expected: FAIL — `Service` type not defined

- [ ] **Step 5: Create service.go**

```go
// pkg/dialog/service.go
package dialog

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the dialog service.
type Options struct{}

// Service is a core.Service managing native dialogs via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenFile:
		paths, err := s.platform.OpenFile(t.Opts)
		return paths, true, err
	case TaskSaveFile:
		path, err := s.platform.SaveFile(t.Opts)
		return path, true, err
	case TaskOpenDirectory:
		path, err := s.platform.OpenDirectory(t.Opts)
		return path, true, err
	case TaskMessageDialog:
		button, err := s.platform.MessageDialog(t.Opts)
		return button, true, err
	default:
		return nil, false, nil
	}
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/dialog/ -v`
Expected: PASS (6 tests)

- [ ] **Step 7: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/dialog/
git commit -m "feat(dialog): add dialog core.Service with Platform interface and IPC"
```

---

## Chunk 2: Notification + Environment

### Task 3: Create pkg/notification

**Files:**
- Create: `pkg/notification/platform.go`
- Create: `pkg/notification/messages.go`
- Create: `pkg/notification/service.go`
- Create: `pkg/notification/service_test.go`
- Delete: `pkg/display/notification.go` (after Task 7)

- [ ] **Step 1: Create platform.go with types**

```go
// pkg/notification/platform.go
package notification

// Platform abstracts the native notification backend.
type Platform interface {
	Send(opts NotificationOptions) error
	RequestPermission() (bool, error)
	CheckPermission() (bool, error)
}

// NotificationSeverity indicates the severity for dialog fallback.
type NotificationSeverity int

const (
	SeverityInfo NotificationSeverity = iota
	SeverityWarning
	SeverityError
)

// NotificationOptions contains options for sending a notification.
type NotificationOptions struct {
	ID       string               `json:"id,omitempty"`
	Title    string               `json:"title"`
	Message  string               `json:"message"`
	Subtitle string               `json:"subtitle,omitempty"`
	Severity NotificationSeverity `json:"severity,omitempty"`
}

// PermissionStatus indicates whether notifications are authorised.
type PermissionStatus struct {
	Granted bool `json:"granted"`
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/notification/messages.go
package notification

// QueryPermission checks notification authorisation. Result: PermissionStatus
type QueryPermission struct{}

// TaskSend sends a notification. Falls back to dialog if platform fails.
type TaskSend struct{ Opts NotificationOptions }

// TaskRequestPermission requests notification authorisation. Result: bool (granted)
type TaskRequestPermission struct{}

// ActionNotificationClicked is broadcast when a notification is clicked (future).
type ActionNotificationClicked struct{ ID string }
```

- [ ] **Step 3: Write failing test**

```go
// pkg/notification/service_test.go
package notification

import (
	"context"
	"errors"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	sendErr       error
	permGranted   bool
	permErr       error
	lastOpts      NotificationOptions
	sendCalled    bool
}

func (m *mockPlatform) Send(opts NotificationOptions) error {
	m.sendCalled = true
	m.lastOpts = opts
	return m.sendErr
}
func (m *mockPlatform) RequestPermission() (bool, error) { return m.permGranted, m.permErr }
func (m *mockPlatform) CheckPermission() (bool, error)   { return m.permGranted, m.permErr }

// mockDialogPlatform tracks whether MessageDialog was called (for fallback test).
type mockDialogPlatform struct {
	messageCalled bool
	lastMsgOpts   dialog.MessageDialogOptions
}

func (m *mockDialogPlatform) OpenFile(opts dialog.OpenFileOptions) ([]string, error) { return nil, nil }
func (m *mockDialogPlatform) SaveFile(opts dialog.SaveFileOptions) (string, error)   { return "", nil }
func (m *mockDialogPlatform) OpenDirectory(opts dialog.OpenDirectoryOptions) (string, error) {
	return "", nil
}
func (m *mockDialogPlatform) MessageDialog(opts dialog.MessageDialogOptions) (string, error) {
	m.messageCalled = true
	m.lastMsgOpts = opts
	return "OK", nil
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{permGranted: true}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	return mock, c
}

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "notification")
	assert.NotNil(t, svc)
}

func TestTaskSend_Good(t *testing.T) {
	mock, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskSend{
		Opts: NotificationOptions{Title: "Test", Message: "Hello"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.sendCalled)
	assert.Equal(t, "Test", mock.lastOpts.Title)
}

func TestTaskSend_Fallback_Good(t *testing.T) {
	// Platform fails → falls back to dialog via IPC
	mockNotify := &mockPlatform{sendErr: errors.New("no permission")}
	mockDlg := &mockDialogPlatform{}
	c, err := core.New(
		core.WithService(dialog.Register(mockDlg)),
		core.WithService(Register(mockNotify)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	_, handled, err := c.PERFORM(TaskSend{
		Opts: NotificationOptions{Title: "Warn", Message: "Oops", Severity: SeverityWarning},
	})
	assert.True(t, handled)
	assert.NoError(t, err) // fallback succeeds even though platform failed
	assert.True(t, mockDlg.messageCalled)
	assert.Equal(t, dialog.DialogWarning, mockDlg.lastMsgOpts.Type)
}

func TestQueryPermission_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryPermission{})
	require.NoError(t, err)
	assert.True(t, handled)
	status := result.(PermissionStatus)
	assert.True(t, status.Granted)
}

func TestTaskRequestPermission_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskRequestPermission{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, true, result)
}

func TestTaskSend_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskSend{})
	assert.False(t, handled)
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/notification/ -v`
Expected: FAIL — `Service` type not defined

- [ ] **Step 5: Create service.go**

```go
// pkg/notification/service.go
package notification

import (
	"context"
	"fmt"
	"time"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
)

// Options holds configuration for the notification service.
type Options struct{}

// Service is a core.Service managing notifications via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryPermission:
		granted, err := s.platform.CheckPermission()
		return PermissionStatus{Granted: granted}, true, err
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSend:
		return nil, true, s.send(t.Opts)
	case TaskRequestPermission:
		granted, err := s.platform.RequestPermission()
		return granted, true, err
	default:
		return nil, false, nil
	}
}

// send attempts native notification, falls back to dialog via IPC.
func (s *Service) send(opts NotificationOptions) error {
	// Generate ID if not provided
	if opts.ID == "" {
		opts.ID = fmt.Sprintf("core-%d", time.Now().UnixNano())
	}

	if err := s.platform.Send(opts); err != nil {
		// Fallback: show as dialog via IPC
		return s.fallbackDialog(opts)
	}
	return nil
}

// fallbackDialog shows a dialog via IPC when native notifications fail.
func (s *Service) fallbackDialog(opts NotificationOptions) error {
	// Map severity to dialog type
	var dt dialog.DialogType
	switch opts.Severity {
	case SeverityWarning:
		dt = dialog.DialogWarning
	case SeverityError:
		dt = dialog.DialogError
	default:
		dt = dialog.DialogInfo
	}

	msg := opts.Message
	if opts.Subtitle != "" {
		msg = opts.Subtitle + "\n\n" + msg
	}

	_, _, err := s.Core().PERFORM(dialog.TaskMessageDialog{
		Opts: dialog.MessageDialogOptions{
			Type:    dt,
			Title:   opts.Title,
			Message: msg,
			Buttons: []string{"OK"},
		},
	})
	return err
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/notification/ -v`
Expected: PASS (6 tests)

- [ ] **Step 7: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/notification/
git commit -m "feat(notification): add notification core.Service with fallback to dialog via IPC"
```

---

### Task 4: Create pkg/environment

**Files:**
- Create: `pkg/environment/platform.go`
- Create: `pkg/environment/messages.go`
- Create: `pkg/environment/service.go`
- Create: `pkg/environment/service_test.go`
- Delete: `pkg/display/theme.go` (after Task 7)

- [ ] **Step 1: Create platform.go with types**

```go
// pkg/environment/platform.go
package environment

// Platform abstracts environment and theme backend queries.
type Platform interface {
	IsDarkMode() bool
	Info() EnvironmentInfo
	AccentColour() string
	OpenFileManager(path string, selectFile bool) error
	OnThemeChange(handler func(isDark bool)) func() // returns cancel func
}

// EnvironmentInfo contains system environment details.
type EnvironmentInfo struct {
	OS       string       `json:"os"`
	Arch     string       `json:"arch"`
	Debug    bool         `json:"debug"`
	Platform PlatformInfo `json:"platform"`
}

// PlatformInfo contains platform-specific details.
type PlatformInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ThemeInfo contains the current theme state.
type ThemeInfo struct {
	IsDark bool   `json:"isDark"`
	Theme  string `json:"theme"` // "dark" or "light"
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/environment/messages.go
package environment

// QueryTheme returns the current theme. Result: ThemeInfo
type QueryTheme struct{}

// QueryInfo returns environment information. Result: EnvironmentInfo
type QueryInfo struct{}

// QueryAccentColour returns the system accent colour. Result: string
type QueryAccentColour struct{}

// TaskOpenFileManager opens the system file manager. Result: error only
type TaskOpenFileManager struct {
	Path   string `json:"path"`
	Select bool   `json:"select"`
}

// ActionThemeChanged is broadcast when the system theme changes.
type ActionThemeChanged struct {
	IsDark bool `json:"isDark"`
}
```

- [ ] **Step 3: Write failing test**

```go
// pkg/environment/service_test.go
package environment

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	isDark       bool
	info         EnvironmentInfo
	accentColour string
	openFMErr    error
	themeHandler func(isDark bool)
	mu           sync.Mutex
}

func (m *mockPlatform) IsDarkMode() bool        { return m.isDark }
func (m *mockPlatform) Info() EnvironmentInfo    { return m.info }
func (m *mockPlatform) AccentColour() string     { return m.accentColour }
func (m *mockPlatform) OpenFileManager(path string, selectFile bool) error {
	return m.openFMErr
}
func (m *mockPlatform) OnThemeChange(handler func(isDark bool)) func() {
	m.mu.Lock()
	m.themeHandler = handler
	m.mu.Unlock()
	return func() {
		m.mu.Lock()
		m.themeHandler = nil
		m.mu.Unlock()
	}
}

// simulateThemeChange triggers the stored handler (test helper).
func (m *mockPlatform) simulateThemeChange(isDark bool) {
	m.mu.Lock()
	h := m.themeHandler
	m.mu.Unlock()
	if h != nil {
		h(isDark)
	}
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		isDark:       true,
		accentColour: "rgb(0,122,255)",
		info: EnvironmentInfo{
			OS: "darwin", Arch: "arm64",
			Platform: PlatformInfo{Name: "macOS", Version: "14.0"},
		},
	}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	return mock, c
}

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "environment")
	assert.NotNil(t, svc)
}

func TestQueryTheme_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryTheme{})
	require.NoError(t, err)
	assert.True(t, handled)
	theme := result.(ThemeInfo)
	assert.True(t, theme.IsDark)
	assert.Equal(t, "dark", theme.Theme)
}

func TestQueryInfo_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryInfo{})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(EnvironmentInfo)
	assert.Equal(t, "darwin", info.OS)
	assert.Equal(t, "arm64", info.Arch)
}

func TestQueryAccentColour_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryAccentColour{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "rgb(0,122,255)", result)
}

func TestTaskOpenFileManager_Good(t *testing.T) {
	_, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskOpenFileManager{Path: "/tmp", Select: true})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestThemeChange_ActionBroadcast_Good(t *testing.T) {
	mock, c := newTestService(t)

	// Register a listener that captures the action
	var received *ActionThemeChanged
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionThemeChanged); ok {
			mu.Lock()
			received = &a
			mu.Unlock()
		}
		return nil
	})

	// Simulate theme change
	mock.simulateThemeChange(false)

	mu.Lock()
	r := received
	mu.Unlock()
	require.NotNil(t, r)
	assert.False(t, r.IsDark)
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/environment/ -v`
Expected: FAIL — `Service` type not defined

- [ ] **Step 5: Create service.go**

```go
// pkg/environment/service.go
package environment

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the environment service.
type Options struct{}

// Service is a core.Service providing environment queries and theme change events via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	platform    Platform
	cancelTheme func() // cancel function for theme change listener
}

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers IPC handlers and the theme change listener.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)

	// Register theme change callback — broadcasts ActionThemeChanged via IPC
	s.cancelTheme = s.platform.OnThemeChange(func(isDark bool) {
		_ = s.Core().ACTION(ActionThemeChanged{IsDark: isDark})
	})
	return nil
}

// OnShutdown cancels the theme change listener.
func (s *Service) OnShutdown(ctx context.Context) error {
	if s.cancelTheme != nil {
		s.cancelTheme()
	}
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryTheme:
		isDark := s.platform.IsDarkMode()
		theme := "light"
		if isDark {
			theme = "dark"
		}
		return ThemeInfo{IsDark: isDark, Theme: theme}, true, nil
	case QueryInfo:
		return s.platform.Info(), true, nil
	case QueryAccentColour:
		return s.platform.AccentColour(), true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenFileManager:
		return nil, true, s.platform.OpenFileManager(t.Path, t.Select)
	default:
		return nil, false, nil
	}
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/environment/ -v`
Expected: PASS (6 tests)

- [ ] **Step 7: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/environment/
git commit -m "feat(environment): add environment core.Service with theme change broadcasts"
```

---

## Chunk 3: Screen + Display Orchestrator Update

### Task 5: Create pkg/screen

**Files:**
- Create: `pkg/screen/platform.go`
- Create: `pkg/screen/messages.go`
- Create: `pkg/screen/service.go`
- Create: `pkg/screen/service_test.go`

- [ ] **Step 1: Create platform.go with types**

```go
// pkg/screen/platform.go
package screen

// Platform abstracts the screen/display backend.
type Platform interface {
	GetAll() []Screen
	GetPrimary() *Screen
}

// Screen describes a display/monitor.
type Screen struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ScaleFactor float64 `json:"scaleFactor"`
	Size        Size    `json:"size"`
	Bounds      Rect    `json:"bounds"`
	WorkArea    Rect    `json:"workArea"`
	IsPrimary   bool    `json:"isPrimary"`
	Rotation    float64 `json:"rotation"`
}

// Rect represents a rectangle with position and dimensions.
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Size represents dimensions.
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/screen/messages.go
package screen

// QueryAll returns all screens. Result: []Screen
type QueryAll struct{}

// QueryPrimary returns the primary screen. Result: *Screen (nil if not found)
type QueryPrimary struct{}

// QueryByID returns a screen by ID. Result: *Screen (nil if not found)
type QueryByID struct{ ID string }

// QueryAtPoint returns the screen containing a point. Result: *Screen (nil if none)
type QueryAtPoint struct{ X, Y int }

// QueryWorkAreas returns work areas for all screens. Result: []Rect
type QueryWorkAreas struct{}

// ActionScreensChanged is broadcast when displays change (future).
type ActionScreensChanged struct{ Screens []Screen }
```

- [ ] **Step 3: Write failing test**

```go
// pkg/screen/service_test.go
package screen

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	screens []Screen
}

func (m *mockPlatform) GetAll() []Screen { return m.screens }
func (m *mockPlatform) GetPrimary() *Screen {
	for i := range m.screens {
		if m.screens[i].IsPrimary {
			return &m.screens[i]
		}
	}
	return nil
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		screens: []Screen{
			{
				ID: "1", Name: "Built-in", IsPrimary: true,
				Bounds:   Rect{X: 0, Y: 0, Width: 2560, Height: 1600},
				WorkArea: Rect{X: 0, Y: 38, Width: 2560, Height: 1562},
				Size:     Size{Width: 2560, Height: 1600},
			},
			{
				ID: "2", Name: "External",
				Bounds:   Rect{X: 2560, Y: 0, Width: 1920, Height: 1080},
				WorkArea: Rect{X: 2560, Y: 0, Width: 1920, Height: 1080},
				Size:     Size{Width: 1920, Height: 1080},
			},
		},
	}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	return mock, c
}

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "screen")
	assert.NotNil(t, svc)
}

func TestQueryAll_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryAll{})
	require.NoError(t, err)
	assert.True(t, handled)
	screens := result.([]Screen)
	assert.Len(t, screens, 2)
}

func TestQueryPrimary_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryPrimary{})
	require.NoError(t, err)
	assert.True(t, handled)
	scr := result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "Built-in", scr.Name)
	assert.True(t, scr.IsPrimary)
}

func TestQueryByID_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryByID{ID: "2"})
	require.NoError(t, err)
	assert.True(t, handled)
	scr := result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

func TestQueryByID_Bad(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryByID{ID: "99"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
}

func TestQueryAtPoint_Good(t *testing.T) {
	_, c := newTestService(t)

	// Point on primary screen
	result, handled, err := c.QUERY(QueryAtPoint{X: 100, Y: 100})
	require.NoError(t, err)
	assert.True(t, handled)
	scr := result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "Built-in", scr.Name)

	// Point on external screen
	result, _, _ = c.QUERY(QueryAtPoint{X: 3000, Y: 500})
	scr = result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

func TestQueryAtPoint_Bad(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryAtPoint{X: -1000, Y: -1000})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
}

func TestQueryWorkAreas_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryWorkAreas{})
	require.NoError(t, err)
	assert.True(t, handled)
	areas := result.([]Rect)
	assert.Len(t, areas, 2)
	assert.Equal(t, 38, areas[0].Y) // primary has menu bar offset
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/screen/ -v`
Expected: FAIL — `Service` type not defined

- [ ] **Step 5: Create service.go**

```go
// pkg/screen/service.go
package screen

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the screen service.
type Options struct{}

// Service is a core.Service providing screen/display queries via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryAll:
		return s.platform.GetAll(), true, nil
	case QueryPrimary:
		return s.platform.GetPrimary(), true, nil
	case QueryByID:
		return s.queryByID(q.ID), true, nil
	case QueryAtPoint:
		return s.queryAtPoint(q.X, q.Y), true, nil
	case QueryWorkAreas:
		return s.queryWorkAreas(), true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) queryByID(id string) *Screen {
	for _, scr := range s.platform.GetAll() {
		if scr.ID == id {
			return &scr
		}
	}
	return nil
}

func (s *Service) queryAtPoint(x, y int) *Screen {
	for _, scr := range s.platform.GetAll() {
		b := scr.Bounds
		if x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height {
			return &scr
		}
	}
	return nil
}

func (s *Service) queryWorkAreas() []Rect {
	screens := s.platform.GetAll()
	areas := make([]Rect, len(screens))
	for i, scr := range screens {
		areas[i] = scr.WorkArea
	}
	return areas
}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/screen/ -v`
Expected: PASS (8 tests)

- [ ] **Step 7: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/screen/
git commit -m "feat(screen): add screen core.Service with computed queries via IPC"
```

---

### Task 6: Update display orchestrator — add new IPC bridge cases

**Files:**
- Modify: `pkg/display/display.go` — add import aliases, add HandleIPCEvents cases for new Actions
- Modify: `pkg/display/events.go` — add `EventNotificationClick` constant

- [ ] **Step 1: Add imports for new packages to display.go**

Add to the import block in `pkg/display/display.go`:

```go
"forge.lthn.ai/core/gui/pkg/clipboard"
"forge.lthn.ai/core/gui/pkg/dialog"
"forge.lthn.ai/core/gui/pkg/environment"
"forge.lthn.ai/core/gui/pkg/notification"
"forge.lthn.ai/core/gui/pkg/screen"
```

- [ ] **Step 2: Add EventNotificationClick constant to events.go**

Add after existing event constants in `pkg/display/events.go`:

```go
EventNotificationClick EventType = "notification.click"
```

- [ ] **Step 3: Add HandleIPCEvents cases for new Action types**

Add to the `HandleIPCEvents` switch in `pkg/display/display.go`, after the existing systray cases:

```go
case environment.ActionThemeChanged:
    if s.events != nil {
        theme := "light"
        if m.IsDark {
            theme = "dark"
        }
        s.events.Emit(Event{Type: EventThemeChange,
            Data: map[string]any{"isDark": m.IsDark, "theme": theme}})
    }
case notification.ActionNotificationClicked:
    if s.events != nil {
        s.events.Emit(Event{Type: EventNotificationClick,
            Data: map[string]any{"id": m.ID}})
    }
case screen.ActionScreensChanged:
    if s.events != nil {
        s.events.Emit(Event{Type: EventScreenChange,
            Data: map[string]any{"screens": m.Screens}})
    }
```

- [ ] **Step 4: Run all tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/display/ ./pkg/clipboard/ ./pkg/dialog/ ./pkg/notification/ ./pkg/environment/ ./pkg/screen/ -v`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/display/display.go pkg/display/events.go
git commit -m "feat(display): bridge new service Actions to WSEventManager"
```

---

### Task 7: Remove extracted code from display

**Files:**
- Delete: `pkg/display/clipboard.go`
- Delete: `pkg/display/dialog.go`
- Delete: `pkg/display/notification.go`
- Delete: `pkg/display/theme.go`
- Modify: `pkg/display/display.go` — remove ScreenInfo, WorkArea, screen methods, ShowEnvironmentDialog, GetScreenForWindow, notifier field, replace handleTrayAction "env-info" with IPC
- Modify: `pkg/display/events.go` — remove SetupWindowEventListeners, change NewWSEventManager signature
- Modify: `pkg/display/interfaces.go` — remove wailsEventSource (keep wailsDialogManager — still used by handleOpenFile)
- Delete: `pkg/display/types.go` — EventSource interface (only content, file deleted)
- Modify: `pkg/display/mocks_test.go` — remove mockEventSource
- Modify: `pkg/display/display_test.go` — update NewWSEventManager calls, remove SetupWindowEventListeners test

- [ ] **Step 1: Delete the 4 extracted files**

```bash
cd /Users/snider/Code/core/gui
rm pkg/display/clipboard.go pkg/display/dialog.go pkg/display/notification.go pkg/display/theme.go
```

- [ ] **Step 2: Remove ScreenInfo, WorkArea types and screen methods from display.go**

Remove the `ScreenInfo` struct (lines 602-611), `WorkArea` struct (lines 613-620), and ALL screen methods: `GetScreens`, `GetWorkAreas`, `GetPrimaryScreen`, `GetScreen`, `GetScreenAtPoint`, `GetScreenForWindow`, `ShowEnvironmentDialog` (lines 622-763 approximately). Also remove the `notifier` field from the Service struct and the `notifications` import.

- [ ] **Step 3: Replace handleTrayAction "env-info" with IPC**

The `handleTrayAction` method calls `s.ShowEnvironmentDialog()` which is being removed. Replace with IPC calls:

```go
case "env-info":
    // Query environment info via IPC and show as dialog
    result, handled, _ := s.Core().QUERY(environment.QueryInfo{})
    if handled {
        info := result.(environment.EnvironmentInfo)
        details := fmt.Sprintf("OS: %s\nArch: %s\nPlatform: %s %s",
            info.OS, info.Arch, info.Platform.Name, info.Platform.Version)
        _, _, _ = s.Core().PERFORM(dialog.TaskMessageDialog{
            Opts: dialog.MessageDialogOptions{
                Type: dialog.DialogInfo, Title: "Environment",
                Message: details, Buttons: []string{"OK"},
            },
        })
    }
```

- [ ] **Step 4: Delete types.go**

Delete `pkg/display/types.go` — it only contains the `EventSource` interface which is replaced by `environment.Platform.OnThemeChange()`.

```bash
rm pkg/display/types.go
```

- [ ] **Step 5: Remove wailsEventSource from interfaces.go**

Remove `wailsEventSource` struct and its methods (lines 80-94) and `newWailsEventSource()` factory. **Keep** `wailsDialogManager` — it is still used by `handleOpenFile` in display.go. Keep `wailsApp`, `wailsEnvManager`, `wailsEventManager`.

- [ ] **Step 6: Update NewWSEventManager signature in events.go**

Change `NewWSEventManager(es EventSource)` to `NewWSEventManager()`. Remove the `eventSource` field from WSEventManager. Remove the `SetupWindowEventListeners()` method (theme change now comes via IPC ActionThemeChanged).

- [ ] **Step 7: Update display.go OnStartup**

Change the line:
```go
s.events = NewWSEventManager(newWailsEventSource(s.wailsApp))
s.events.SetupWindowEventListeners()
```
To:
```go
s.events = NewWSEventManager()
```

- [ ] **Step 8: Update display test files**

In `pkg/display/mocks_test.go`: remove `mockEventSource` struct and `newMockEventSource()`.

In `pkg/display/display_test.go`:
- Update `TestWSEventManager_Good`: change `NewWSEventManager(es)` to `NewWSEventManager()`, remove `es` variable
- Remove `TestWSEventManager_SetupWindowEventListeners_Good` entirely (method deleted)

- [ ] **Step 9: Run all tests**

Run: `cd /Users/snider/Code/core/gui && go test ./... -v`
Expected: ALL PASS across all packages

- [ ] **Step 10: Commit**

```bash
cd /Users/snider/Code/core/gui
git add -A pkg/display/
git commit -m "refactor(display): remove extracted clipboard/dialog/notification/theme/screen code"
```

---

### Task 8: Final verification and push

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/snider/Code/core/gui && go test ./... -count=1`
Expected: ALL PASS

- [ ] **Step 2: Verify build**

Run: `cd /Users/snider/Code/core/gui && go build ./...`
Expected: Clean build, no errors

- [ ] **Step 3: Push to forge**

```bash
cd /Users/snider/Code/core/gui && git push origin main
```
