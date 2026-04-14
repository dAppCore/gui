// pkg/mcp/tools_window.go
package mcp

import (
	"context"

	"dappco.re/go/core/gui/pkg/screen"
	"dappco.re/go/core/gui/pkg/window"
	coreerr "forge.lthn.ai/core/go-log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- window_list ---

type WindowListInput struct{}
type WindowListOutput struct {
	Windows []window.WindowInfo `json:"windows"`
}

func (s *Subsystem) windowList(_ context.Context, _ *mcp.CallToolRequest, _ WindowListInput) (*mcp.CallToolResult, WindowListOutput, error) {
	result, _, err := s.core.QUERY(window.QueryWindowList{})
	if err != nil {
		return nil, WindowListOutput{}, err
	}
	windows, ok := result.([]window.WindowInfo)
	if !ok {
		return nil, WindowListOutput{}, coreerr.E("mcp.windowList", "unexpected result type", nil)
	}
	return nil, WindowListOutput{Windows: windows}, nil
}

// --- window_get ---

type WindowGetInput struct {
	Name string `json:"name"`
}
type WindowGetOutput struct {
	Window *window.WindowInfo `json:"window"`
}

func (s *Subsystem) windowGet(_ context.Context, _ *mcp.CallToolRequest, input WindowGetInput) (*mcp.CallToolResult, WindowGetOutput, error) {
	result, _, err := s.core.QUERY(window.QueryWindowByName{Name: input.Name})
	if err != nil {
		return nil, WindowGetOutput{}, err
	}
	info, ok := result.(*window.WindowInfo)
	if !ok {
		return nil, WindowGetOutput{}, coreerr.E("mcp.windowGet", "unexpected result type", nil)
	}
	return nil, WindowGetOutput{Window: info}, nil
}

// --- window_focused ---

type WindowFocusedInput struct{}
type WindowFocusedOutput struct {
	Window string `json:"window"`
}

func (s *Subsystem) windowFocused(_ context.Context, _ *mcp.CallToolRequest, _ WindowFocusedInput) (*mcp.CallToolResult, WindowFocusedOutput, error) {
	result, _, err := s.core.QUERY(window.QueryWindowList{})
	if err != nil {
		return nil, WindowFocusedOutput{}, err
	}
	windows, ok := result.([]window.WindowInfo)
	if !ok {
		return nil, WindowFocusedOutput{}, coreerr.E("mcp.windowFocused", "unexpected result type", nil)
	}
	for _, w := range windows {
		if w.Focused {
			return nil, WindowFocusedOutput{Window: w.Name}, nil
		}
	}
	return nil, WindowFocusedOutput{}, nil
}

// --- window_create ---

type WindowCreateInput struct {
	Name   string `json:"name"`
	Title  string `json:"title,omitempty"`
	URL    string `json:"url,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}
type WindowCreateOutput struct {
	Window window.WindowInfo `json:"window"`
}

func (s *Subsystem) windowCreate(_ context.Context, _ *mcp.CallToolRequest, input WindowCreateInput) (*mcp.CallToolResult, WindowCreateOutput, error) {
	result, _, err := s.core.PERFORM(window.TaskOpenWindow{
		Window: &window.Window{
			Name:   input.Name,
			Title:  input.Title,
			URL:    input.URL,
			Width:  input.Width,
			Height: input.Height,
			X:      input.X,
			Y:      input.Y,
		},
	})
	if err != nil {
		return nil, WindowCreateOutput{}, err
	}
	info, ok := result.(window.WindowInfo)
	if !ok {
		return nil, WindowCreateOutput{}, coreerr.E("mcp.windowCreate", "unexpected result type", nil)
	}
	return nil, WindowCreateOutput{Window: info}, nil
}

// --- window_close ---

type WindowCloseInput struct {
	Name string `json:"name"`
}
type WindowCloseOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowClose(_ context.Context, _ *mcp.CallToolRequest, input WindowCloseInput) (*mcp.CallToolResult, WindowCloseOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskCloseWindow{Name: input.Name})
	if err != nil {
		return nil, WindowCloseOutput{}, err
	}
	return nil, WindowCloseOutput{Success: true}, nil
}

// --- window_position ---

type WindowPositionInput struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}
type WindowPositionOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowPosition(_ context.Context, _ *mcp.CallToolRequest, input WindowPositionInput) (*mcp.CallToolResult, WindowPositionOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSetPosition{Name: input.Name, X: input.X, Y: input.Y})
	if err != nil {
		return nil, WindowPositionOutput{}, err
	}
	return nil, WindowPositionOutput{Success: true}, nil
}

// --- window_size ---

type WindowSizeInput struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type WindowSizeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowSize(_ context.Context, _ *mcp.CallToolRequest, input WindowSizeInput) (*mcp.CallToolResult, WindowSizeOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSetSize{Name: input.Name, Width: input.Width, Height: input.Height})
	if err != nil {
		return nil, WindowSizeOutput{}, err
	}
	return nil, WindowSizeOutput{Success: true}, nil
}

// --- window_bounds ---

type WindowBoundsInput struct {
	Name   string `json:"name"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type WindowBoundsOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowBounds(_ context.Context, _ *mcp.CallToolRequest, input WindowBoundsInput) (*mcp.CallToolResult, WindowBoundsOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSetPosition{Name: input.Name, X: input.X, Y: input.Y})
	if err != nil {
		return nil, WindowBoundsOutput{}, err
	}
	_, _, err = s.core.PERFORM(window.TaskSetSize{Name: input.Name, Width: input.Width, Height: input.Height})
	if err != nil {
		return nil, WindowBoundsOutput{}, err
	}
	return nil, WindowBoundsOutput{Success: true}, nil
}

// --- window_maximize ---

type WindowMaximizeInput struct {
	Name string `json:"name"`
}
type WindowMaximizeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowMaximize(_ context.Context, _ *mcp.CallToolRequest, input WindowMaximizeInput) (*mcp.CallToolResult, WindowMaximizeOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskMaximise{Name: input.Name})
	if err != nil {
		return nil, WindowMaximizeOutput{}, err
	}
	return nil, WindowMaximizeOutput{Success: true}, nil
}

// --- window_minimize ---

type WindowMinimizeInput struct {
	Name string `json:"name"`
}
type WindowMinimizeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowMinimize(_ context.Context, _ *mcp.CallToolRequest, input WindowMinimizeInput) (*mcp.CallToolResult, WindowMinimizeOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskMinimise{Name: input.Name})
	if err != nil {
		return nil, WindowMinimizeOutput{}, err
	}
	return nil, WindowMinimizeOutput{Success: true}, nil
}

// --- window_restore ---

type WindowRestoreInput struct {
	Name string `json:"name"`
}
type WindowRestoreOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowRestore(_ context.Context, _ *mcp.CallToolRequest, input WindowRestoreInput) (*mcp.CallToolResult, WindowRestoreOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskRestore{Name: input.Name})
	if err != nil {
		return nil, WindowRestoreOutput{}, err
	}
	return nil, WindowRestoreOutput{Success: true}, nil
}

// --- window_focus ---

type WindowFocusInput struct {
	Name string `json:"name"`
}
type WindowFocusOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowFocus(_ context.Context, _ *mcp.CallToolRequest, input WindowFocusInput) (*mcp.CallToolResult, WindowFocusOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskFocus{Name: input.Name})
	if err != nil {
		return nil, WindowFocusOutput{}, err
	}
	return nil, WindowFocusOutput{Success: true}, nil
}

// --- focus_set ---

type FocusSetInput struct {
	Name string `json:"name"`
}
type FocusSetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) focusSet(ctx context.Context, req *mcp.CallToolRequest, input FocusSetInput) (*mcp.CallToolResult, FocusSetOutput, error) {
	_, out, err := s.windowFocus(ctx, req, WindowFocusInput{Name: input.Name})
	if err != nil {
		return nil, FocusSetOutput{}, err
	}
	return nil, FocusSetOutput{Success: out.Success}, nil
}

// --- window_title ---

type WindowTitleInput struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
type WindowTitleOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowTitle(_ context.Context, _ *mcp.CallToolRequest, input WindowTitleInput) (*mcp.CallToolResult, WindowTitleOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSetTitle{Name: input.Name, Title: input.Title})
	if err != nil {
		return nil, WindowTitleOutput{}, err
	}
	return nil, WindowTitleOutput{Success: true}, nil
}

// --- window_title_get ---

type WindowTitleGetInput struct {
	Name string `json:"name"`
}
type WindowTitleGetOutput struct {
	Title string `json:"title"`
}

func (s *Subsystem) windowTitleGet(_ context.Context, _ *mcp.CallToolRequest, input WindowTitleGetInput) (*mcp.CallToolResult, WindowTitleGetOutput, error) {
	result, _, err := s.core.QUERY(window.QueryWindowByName{Name: input.Name})
	if err != nil {
		return nil, WindowTitleGetOutput{}, err
	}
	info, _ := result.(*window.WindowInfo)
	if info == nil {
		return nil, WindowTitleGetOutput{}, nil
	}
	return nil, WindowTitleGetOutput{Title: info.Title}, nil
}

// --- window_visibility ---

type WindowVisibilityInput struct {
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
}
type WindowVisibilityOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowVisibility(_ context.Context, _ *mcp.CallToolRequest, input WindowVisibilityInput) (*mcp.CallToolResult, WindowVisibilityOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSetVisibility{Name: input.Name, Visible: input.Visible})
	if err != nil {
		return nil, WindowVisibilityOutput{}, err
	}
	return nil, WindowVisibilityOutput{Success: true}, nil
}

// --- window_always_on_top ---

type WindowAlwaysOnTopInput struct {
	Name        string `json:"name"`
	AlwaysOnTop bool   `json:"alwaysOnTop"`
}
type WindowAlwaysOnTopOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowAlwaysOnTop(_ context.Context, _ *mcp.CallToolRequest, input WindowAlwaysOnTopInput) (*mcp.CallToolResult, WindowAlwaysOnTopOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSetAlwaysOnTop{Name: input.Name, AlwaysOnTop: input.AlwaysOnTop})
	if err != nil {
		return nil, WindowAlwaysOnTopOutput{}, err
	}
	return nil, WindowAlwaysOnTopOutput{Success: true}, nil
}

// --- window_background_colour ---

type WindowBackgroundColourInput struct {
	Name  string `json:"name"`
	Red   uint8  `json:"red"`
	Green uint8  `json:"green"`
	Blue  uint8  `json:"blue"`
	Alpha uint8  `json:"alpha"`
}
type WindowBackgroundColourOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowBackgroundColour(_ context.Context, _ *mcp.CallToolRequest, input WindowBackgroundColourInput) (*mcp.CallToolResult, WindowBackgroundColourOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSetBackgroundColour{
		Name: input.Name, Red: input.Red, Green: input.Green, Blue: input.Blue, Alpha: input.Alpha,
	})
	if err != nil {
		return nil, WindowBackgroundColourOutput{}, err
	}
	return nil, WindowBackgroundColourOutput{Success: true}, nil
}

// --- window_fullscreen ---

type WindowFullscreenInput struct {
	Name       string `json:"name"`
	Fullscreen bool   `json:"fullscreen"`
}
type WindowFullscreenOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowFullscreen(_ context.Context, _ *mcp.CallToolRequest, input WindowFullscreenInput) (*mcp.CallToolResult, WindowFullscreenOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskFullscreen{Name: input.Name, Fullscreen: input.Fullscreen})
	if err != nil {
		return nil, WindowFullscreenOutput{}, err
	}
	return nil, WindowFullscreenOutput{Success: true}, nil
}

type arrangedPair struct {
	First  screen.Rect
	Second screen.Rect
}

func (s *Subsystem) arrangePairOnScreen(firstName, secondName string, scr *screen.Screen, orientation string) (arrangedPair, error) {
	workArea := workAreaRect(scr)
	if workArea.Width == 0 || workArea.Height == 0 {
		return arrangedPair{}, coreerr.E("mcp.arrangePairOnScreen", "screen work area is empty", nil)
	}
	if orientation == "" {
		if workArea.Width >= workArea.Height {
			orientation = "horizontal"
		} else {
			orientation = "vertical"
		}
	}

	var firstRect screen.Rect
	var secondRect screen.Rect
	switch orientation {
	case "vertical", "stacked":
		firstHeight := workArea.Height / 2
		firstRect = screen.Rect{X: workArea.X, Y: workArea.Y, Width: workArea.Width, Height: firstHeight}
		secondRect = screen.Rect{X: workArea.X, Y: workArea.Y + firstHeight, Width: workArea.Width, Height: workArea.Height - firstHeight}
	default:
		firstWidth := workArea.Width / 2
		firstRect = screen.Rect{X: workArea.X, Y: workArea.Y, Width: firstWidth, Height: workArea.Height}
		secondRect = screen.Rect{X: workArea.X + firstWidth, Y: workArea.Y, Width: workArea.Width - firstWidth, Height: workArea.Height}
	}

	if err := applyRect(s.core, firstName, firstRect); err != nil {
		return arrangedPair{}, err
	}
	if err := applyRect(s.core, secondName, secondRect); err != nil {
		return arrangedPair{}, err
	}
	return arrangedPair{First: firstRect, Second: secondRect}, nil
}

// --- window_arrange_pair ---

type WindowArrangePairInput struct {
	First       string `json:"first"`
	Second      string `json:"second"`
	ScreenID    string `json:"screenId,omitempty"`
	Orientation string `json:"orientation,omitempty"`
}
type WindowArrangePairOutput struct {
	FirstBounds  screen.Rect `json:"firstBounds"`
	SecondBounds screen.Rect `json:"secondBounds"`
}

func (s *Subsystem) windowArrangePair(_ context.Context, _ *mcp.CallToolRequest, input WindowArrangePairInput) (*mcp.CallToolResult, WindowArrangePairOutput, error) {
	screens, err := s.allScreens()
	if err != nil {
		return nil, WindowArrangePairOutput{}, err
	}
	windows, err := s.allWindows()
	if err != nil {
		return nil, WindowArrangePairOutput{}, err
	}

	var targetScreen *screen.Screen
	if input.ScreenID != "" {
		targetScreen = chooseScreenByIDOrPrimary(screens, input.ScreenID)
	} else {
		for _, info := range windows {
			if info.Name == input.First {
				targetScreen = screenForWindowInfo(screens, info)
				break
			}
		}
		if targetScreen == nil {
			targetScreen = chooseScreenByIDOrPrimary(screens, "")
		}
	}
	if targetScreen == nil {
		return nil, WindowArrangePairOutput{}, coreerr.E("mcp.windowArrangePair", "no screen available", nil)
	}

	arranged, err := s.arrangePairOnScreen(input.First, input.Second, targetScreen, input.Orientation)
	if err != nil {
		return nil, WindowArrangePairOutput{}, err
	}
	return nil, WindowArrangePairOutput{FirstBounds: arranged.First, SecondBounds: arranged.Second}, nil
}

// --- Registration ---

func (s *Subsystem) registerWindowTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "window_list", Description: "List all application windows"}, s.windowList)
	mcp.AddTool(server, &mcp.Tool{Name: "window_get", Description: "Get information about a specific window"}, s.windowGet)
	mcp.AddTool(server, &mcp.Tool{Name: "window_focused", Description: "Get the currently focused window"}, s.windowFocused)
	mcp.AddTool(server, &mcp.Tool{Name: "window_create", Description: "Create a new application window"}, s.windowCreate)
	mcp.AddTool(server, &mcp.Tool{Name: "window_close", Description: "Close an application window"}, s.windowClose)
	mcp.AddTool(server, &mcp.Tool{Name: "window_position", Description: "Set the position of a window"}, s.windowPosition)
	mcp.AddTool(server, &mcp.Tool{Name: "window_size", Description: "Set the size of a window"}, s.windowSize)
	mcp.AddTool(server, &mcp.Tool{Name: "window_bounds", Description: "Set both position and size of a window"}, s.windowBounds)
	mcp.AddTool(server, &mcp.Tool{Name: "window_maximize", Description: "Maximise a window"}, s.windowMaximize)
	mcp.AddTool(server, &mcp.Tool{Name: "window_minimize", Description: "Minimise a window"}, s.windowMinimize)
	mcp.AddTool(server, &mcp.Tool{Name: "window_restore", Description: "Restore a maximised or minimised window"}, s.windowRestore)
	mcp.AddTool(server, &mcp.Tool{Name: "window_focus", Description: "Bring a window to the front"}, s.windowFocus)
	mcp.AddTool(server, &mcp.Tool{Name: "focus_set", Description: "Set focus to a specific window"}, s.focusSet)
	mcp.AddTool(server, &mcp.Tool{Name: "window_title", Description: "Set the title of a window"}, s.windowTitle)
	mcp.AddTool(server, &mcp.Tool{Name: "window_title_get", Description: "Get the title of a window"}, s.windowTitleGet)
	mcp.AddTool(server, &mcp.Tool{Name: "window_visibility", Description: "Show or hide a window"}, s.windowVisibility)
	mcp.AddTool(server, &mcp.Tool{Name: "window_always_on_top", Description: "Pin a window above others"}, s.windowAlwaysOnTop)
	mcp.AddTool(server, &mcp.Tool{Name: "window_background_colour", Description: "Set a window background colour"}, s.windowBackgroundColour)
	mcp.AddTool(server, &mcp.Tool{Name: "window_fullscreen", Description: "Set a window to fullscreen mode"}, s.windowFullscreen)
	mcp.AddTool(server, &mcp.Tool{Name: "window_arrange_pair", Description: "Arrange two windows side-by-side or stacked on a screen"}, s.windowArrangePair)
}
