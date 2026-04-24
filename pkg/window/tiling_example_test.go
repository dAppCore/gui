package window

import "fmt"

func ExampleBesideEditor() {
	rect := BesideEditor(
		Rect{X: 240, Y: 40, Width: 960, Height: 720},
		Size{Width: 1600, Height: 900},
	)

	fmt.Printf("%d,%d %dx%d\n", rect.X, rect.Y, rect.Width, rect.Height)
	// Output:
	// 1200,40 400x720
}

func ExampleSuggestLayout() {
	placements := SuggestLayout(
		[]Window{{Name: "editor"}, {Name: "preview"}},
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
	)

	for _, placement := range placements {
		fmt.Printf("%s:%d,%d %dx%d\n", placement.Name, placement.Bounds.X, placement.Bounds.Y, placement.Bounds.Width, placement.Bounds.Height)
	}
	// Output:
	// editor:0,0 989x900
	// preview:989,0 611x900
}

func ExampleFindEmptySpace() {
	rect, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
		[]Window{{Name: "editor", X: 0, Y: 0, Width: 1000, Height: 900}},
		Size{Width: 300, Height: 300},
	)

	fmt.Printf("%t %d,%d %dx%d\n", ok, rect.X, rect.Y, rect.Width, rect.Height)
	// Output:
	// true 1000,0 600x900
}

func ExampleArrangePair() {
	left, right := ArrangePair(
		Window{Name: "editor", Width: 1600, Height: 900},
		Window{Name: "chat", Width: 900, Height: 1600},
		Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	)

	fmt.Printf("%dx%d | %dx%d\n", left.Width, left.Height, right.Width, right.Height)
	// Output:
	// 1200x1000 | 800x1000
}
