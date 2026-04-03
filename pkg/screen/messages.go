// pkg/screen/messages.go
package screen

// QueryAll returns all screens. Result: []Screen
// Use: result, _, err := c.QUERY(screen.QueryAll{})
type QueryAll struct{}

// QueryPrimary returns the primary screen. Result: *Screen (nil if not found)
// Use: result, _, err := c.QUERY(screen.QueryPrimary{})
type QueryPrimary struct{}

// QueryByID returns a screen by ID. Result: *Screen (nil if not found)
// Use: result, _, err := c.QUERY(screen.QueryByID{ID: "display-1"})
type QueryByID struct{ ID string }

// QueryAtPoint returns the screen containing a point. Result: *Screen (nil if none)
// Use: result, _, err := c.QUERY(screen.QueryAtPoint{X: 100, Y: 100})
type QueryAtPoint struct{ X, Y int }

// QueryWorkAreas returns work areas for all screens. Result: []Rect
// Use: result, _, err := c.QUERY(screen.QueryWorkAreas{})
type QueryWorkAreas struct{}

// ActionScreensChanged is broadcast when displays change.
// Use: _ = c.ACTION(screen.ActionScreensChanged{Screens: screens})
type ActionScreensChanged struct{ Screens []Screen }
