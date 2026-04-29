//go:build compliance

package application

import core "dappco.re/go"

func ExampleNewMenuItem() {
	core.Println("NewMenuItem")
	// Output:
	// NewMenuItem
}

func ExampleNewMenuItemSeparator() {
	core.Println("NewMenuItemSeparator")
	// Output:
	// NewMenuItemSeparator
}

func ExampleNewMenuItemCheckbox() {
	core.Println("NewMenuItemCheckbox")
	// Output:
	// NewMenuItemCheckbox
}

func ExampleNewMenuItemRadio() {
	core.Println("NewMenuItemRadio")
	// Output:
	// NewMenuItemRadio
}

func ExampleNewSubMenuItem() {
	core.Println("NewSubMenuItem")
	// Output:
	// NewSubMenuItem
}

func ExampleNewRole() {
	core.Println("NewRole")
	// Output:
	// NewRole
}

func ExampleNewServicesMenu() {
	core.Println("NewServicesMenu")
	// Output:
	// NewServicesMenu
}

func ExampleMenuItem_GetAccelerator() {
	core.Println("MenuItem_GetAccelerator")
	// Output:
	// MenuItem_GetAccelerator
}

func ExampleMenuItem_GetSubmenu() {
	core.Println("MenuItem_GetSubmenu")
	// Output:
	// MenuItem_GetSubmenu
}
