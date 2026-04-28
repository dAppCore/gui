package application

import (
	"context"
	core "dappco.re/go"
)

type namedService struct{}

func (namedService) ServiceName() string { return "named" }

type plainService struct{}

func TestServices_NewService_Good(t *core.T) {
	instance := &plainService{}
	service := NewService(instance)

	core.AssertSame(t, instance, service.Instance())
	core.AssertEqual(t, DefaultServiceOptions, service.Options())
	core.AssertEqual(t, "application.plainService", getServiceName(service))
}

func TestServices_NewService_Bad(t *core.T) {
	instance := &namedService{}
	service := NewServiceWithOptions(instance, ServiceOptions{Name: "explicit"})

	core.AssertSame(t, instance, service.Instance())
	core.AssertEqual(t, "explicit", getServiceName(service))
}

func TestServices_NewService_Ugly(t *core.T) {
	instance := &namedService{}
	service := NewService(instance)

	core.AssertEqual(t, "named", getServiceName(service))
}

func TestServices_ServiceInterfaces_Good(t *core.T) {
	var _ ServiceName = namedService{}
	var _ ServiceStartup = (*startupService)(nil)
	var _ ServiceShutdown = (*shutdownService)(nil)
}

type startupService struct{}

func (*startupService) ServiceStartup(context.Context, ServiceOptions) error { return nil }

type shutdownService struct{}

func (*shutdownService) ServiceShutdown() error { return nil }
