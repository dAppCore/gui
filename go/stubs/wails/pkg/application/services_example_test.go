//go:build compliance

package application

import core "dappco.re/go"

func ExampleNewService() {
	core.Println("NewService")
	// Output:
	// NewService
}

func ExampleNewServiceWithOptions() {
	core.Println("NewServiceWithOptions")
	// Output:
	// NewServiceWithOptions
}

func ExampleService_Instance() {
	core.Println("Service_Instance")
	// Output:
	// Service_Instance
}

func ExampleService_Options() {
	core.Println("Service_Options")
	// Output:
	// Service_Options
}
