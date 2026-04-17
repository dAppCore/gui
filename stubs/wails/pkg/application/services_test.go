package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type namedService struct{}

func (namedService) ServiceName() string { return "named" }

type plainService struct{}

func TestServices_NewService_Good(t *testing.T) {
	instance := &plainService{}
	service := NewService(instance)

	assert.Same(t, instance, service.Instance())
	assert.Equal(t, DefaultServiceOptions, service.Options())
	assert.Equal(t, "application.plainService", getServiceName(service))
}

func TestServices_NewService_Bad(t *testing.T) {
	instance := &namedService{}
	service := NewServiceWithOptions(instance, ServiceOptions{Name: "explicit"})

	assert.Same(t, instance, service.Instance())
	assert.Equal(t, "explicit", getServiceName(service))
}

func TestServices_NewService_Ugly(t *testing.T) {
	instance := &namedService{}
	service := NewService(instance)

	assert.Equal(t, "named", getServiceName(service))
}

func TestServices_ServiceInterfaces_Good(t *testing.T) {
	var _ ServiceName = namedService{}
	var _ ServiceStartup = (*startupService)(nil)
	var _ ServiceShutdown = (*shutdownService)(nil)
}

type startupService struct{}

func (*startupService) ServiceStartup(context.Context, ServiceOptions) error { return nil }

type shutdownService struct{}

func (*shutdownService) ServiceShutdown() error { return nil }
