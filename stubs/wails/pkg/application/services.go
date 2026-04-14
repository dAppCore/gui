package application

import "context"

// Service wraps a bound type instance for registration with the application.
// svc := NewService(&MyService{})
type Service struct {
	instance any
	options  ServiceOptions
}

// ServiceOptions provides optional configuration for a Service.
// opts := ServiceOptions{Name: "my-service", Route: "/api/my"}
type ServiceOptions struct {
	// Name overrides the service name used in logging and debugging.
	// Defaults to the value from ServiceName interface, or the type name.
	Name string

	// Route mounts the service on the internal asset server at this path
	// when the instance implements http.Handler.
	Route string

	// MarshalError serialises errors returned by service methods to JSON.
	// Return nil to fall back to the globally configured error handler.
	MarshalError func(error) []byte
}

// DefaultServiceOptions is the default service options used by NewService.
var DefaultServiceOptions = ServiceOptions{}

// NewService wraps instance as a Service.
// svc := NewService(&Calculator{})
func NewService[T any](instance *T) Service {
	return Service{instance: instance, options: DefaultServiceOptions}
}

// NewServiceWithOptions wraps instance as a Service with explicit options.
// svc := NewServiceWithOptions(&Calculator{}, ServiceOptions{Name: "calculator"})
func NewServiceWithOptions[T any](instance *T, options ServiceOptions) Service {
	service := NewService(instance)
	service.options = options
	return service
}

// Instance returns the underlying pointer passed to NewService.
// raw := svc.Instance().(*Calculator)
func (s Service) Instance() any {
	return s.instance
}

// Options returns the ServiceOptions for this service.
// opts := svc.Options()
func (s Service) Options() ServiceOptions {
	return s.options
}

// ServiceName is an optional interface that service instances may implement
// to provide a human-readable name for logging and debugging.
// func (s *MyService) ServiceName() string { return "my-service" }
type ServiceName interface {
	ServiceName() string
}

// ServiceStartup is an optional interface for services that need initialisation.
// Called in registration order during application startup.
// func (s *MyService) ServiceStartup(ctx context.Context, opts ServiceOptions) error { ... }
type ServiceStartup interface {
	ServiceStartup(ctx context.Context, options ServiceOptions) error
}

// ServiceShutdown is an optional interface for services that need cleanup.
// Called in reverse registration order during application shutdown.
// func (s *MyService) ServiceShutdown() error { ... }
type ServiceShutdown interface {
	ServiceShutdown() error
}
