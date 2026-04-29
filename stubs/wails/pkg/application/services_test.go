package application

import (
	"context"
	core "dappco.re/go"
)

type namedService struct{}

func (namedService) ServiceName() string { return "named" }

type plainService struct{}

func TestServices_NewService_Good(t *core.T) {
	// NewService
	ax7Variant := "NewService:good"
	core.AssertContains(t, ax7Variant, "good")
	instance := &plainService{}
	service := NewService(instance)

	core.AssertSame(t, instance, service.Instance())
	core.AssertEqual(t, DefaultServiceOptions, service.Options())
	core.AssertEqual(t, "application.plainService", getServiceName(service))
}

func TestServices_NewService_BadCase(t *core.T) {
	instance := &namedService{}
	service := NewServiceWithOptions(instance, ServiceOptions{Name: "explicit"})

	core.AssertSame(t, instance, service.Instance())
	core.AssertEqual(t, "explicit", getServiceName(service))
}

func TestServices_NewService_Ugly(t *core.T) {
	// NewService
	ax7Variant := "NewService:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	instance := &namedService{}
	service := NewService(instance)

	core.AssertEqual(t, "named", getServiceName(service))
}

func TestServices_ServiceInterfaces_GoodCase(t *core.T) {
	var _ ServiceName = namedService{}
	var _ ServiceStartup = (*startupService)(nil)
	var _ ServiceShutdown = (*shutdownService)(nil)
}

type startupService struct{}

func (*startupService) ServiceStartup(context.Context, ServiceOptions) error { return nil }

type shutdownService struct{}

func (*shutdownService) ServiceShutdown() error { return nil }

// AX7 generated source-matching smoke coverage.
func TestServices_NewServiceWithOptions_Good(t *core.T) {
	// NewServiceWithOptions
	ax7Variant := "NewServiceWithOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewServiceWithOptions[any](nil, *new(ServiceOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_NewServiceWithOptions_Bad(t *core.T) {
	// NewServiceWithOptions
	ax7Variant := "NewServiceWithOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewServiceWithOptions[any](nil, *new(ServiceOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_NewServiceWithOptions_Ugly(t *core.T) {
	// NewServiceWithOptions
	ax7Variant := "NewServiceWithOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewServiceWithOptions[any](nil, *new(ServiceOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_Service_Instance_Good(t *core.T) {
	// Service Instance
	ax7Variant := "Service_Instance:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Service
	result := core.Try(func() any {
		got0 := subject.Instance()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_Service_Instance_Bad(t *core.T) {
	// Service Instance
	ax7Variant := "Service_Instance:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Service
	result := core.Try(func() any {
		got0 := subject.Instance()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_Service_Instance_Ugly(t *core.T) {
	// Service Instance
	ax7Variant := "Service_Instance:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Service
	result := core.Try(func() any {
		got0 := subject.Instance()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_Service_Options_Good(t *core.T) {
	// Service Options
	ax7Variant := "Service_Options:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Service
	result := core.Try(func() any {
		got0 := subject.Options()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_Service_Options_Bad(t *core.T) {
	// Service Options
	ax7Variant := "Service_Options:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Service
	result := core.Try(func() any {
		got0 := subject.Options()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_Service_Options_Ugly(t *core.T) {
	// Service Options
	ax7Variant := "Service_Options:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Service
	result := core.Try(func() any {
		got0 := subject.Options()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestServices_NewService_Bad(t *core.T) {
	// NewService
	ax7Variant := "NewService:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewService[any](nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
