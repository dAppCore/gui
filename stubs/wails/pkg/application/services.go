package application

import (
	"context"
	"reflect"
)

// ServiceOptions provides optional parameters when constructing a Service.
//
//	svc := NewServiceWithOptions(&myImpl, ServiceOptions{Name: "MyService", Route: "/api"})
type ServiceOptions struct {
	// Name overrides the service name used for logging and debugging.
	// When empty the name is derived from the ServiceName interface or the type name.
	Name string

	// Route mounts the service on the internal asset server at this prefix
	// if the service instance implements http.Handler.
	Route string

	// MarshalError serialises error values returned by service methods to JSON.
	// A nil return falls back to the globally configured error handler.
	MarshalError func(error) []byte
}

// DefaultServiceOptions holds the default values used when no ServiceOptions
// are passed to NewService.
var DefaultServiceOptions = ServiceOptions{}

// Service wraps a bound type instance for registration with the application.
// The zero value is invalid; obtain valid values via NewService or NewServiceWithOptions.
//
//	svc := NewService(&myStruct)
//	app.RegisterService(svc)
type Service struct {
	instance any
	options  ServiceOptions
}

// NewService wraps instance in a Service using DefaultServiceOptions.
//
//	svc := NewService(&MyController{})
func NewService[T any](instance *T) Service {
	return Service{instance: instance, options: DefaultServiceOptions}
}

// NewServiceWithOptions wraps instance in a Service with explicit options.
//
//	svc := NewServiceWithOptions(&MyController{}, ServiceOptions{Name: "ctrl", Route: "/ctrl"})
func NewServiceWithOptions[T any](instance *T, options ServiceOptions) Service {
	service := NewService(instance)
	service.options = options
	return service
}

// Instance returns the underlying pointer originally passed to NewService.
//
//	ctrl := svc.Instance().(*MyController)
func (s Service) Instance() any {
	return s.instance
}

// Options returns the ServiceOptions associated with this service.
func (s Service) Options() ServiceOptions {
	return s.options
}

// ServiceName is an optional interface that service instances may implement
// to provide a human-readable name for logging and debugging.
//
//	func (s *MyService) ServiceName() string { return "MyService" }
type ServiceName interface {
	ServiceName() string
}

// ServiceStartup is an optional interface for service initialisation.
// Called during application startup in registration order.
// A non-nil return aborts startup and surfaces the error to the caller.
//
//	func (s *MyService) ServiceStartup(ctx context.Context, opts ServiceOptions) error {
//	    return s.connect(ctx)
//	}
type ServiceStartup interface {
	ServiceStartup(ctx context.Context, options ServiceOptions) error
}

// ServiceShutdown is an optional interface for service cleanup.
// Called during application shutdown in reverse registration order,
// after all user-provided shutdown hooks have run.
//
//	func (s *MyService) ServiceShutdown() error { return s.db.Close() }
type ServiceShutdown interface {
	ServiceShutdown() error
}

// getServiceName resolves the display name for a service using, in order:
// the explicit options name, the ServiceName interface, and the reflect type name.
func getServiceName(service Service) string {
	if service.options.Name != "" {
		return service.options.Name
	}
	if named, ok := service.Instance().(ServiceName); ok {
		return named.ServiceName()
	}
	return reflect.TypeOf(service.Instance()).Elem().String()
}
