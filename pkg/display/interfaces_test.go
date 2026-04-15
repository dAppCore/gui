package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestInterfaces_newWailsApp_Good(t *testing.T) {
	app := &application.App{Logger: application.Logger{}}
	wrapped := newWailsApp(app)

	require.NotNil(t, wrapped)
	require.NotNil(t, wrapped.Logger())
	assert.NotPanics(t, func() {
		wrapped.Quit()
		wrapped.Logger().Info("ready")
	})
}

func TestInterfaces_newWailsApp_Bad(t *testing.T) {
	wrapped := newWailsApp(&application.App{})
	require.NotNil(t, wrapped)
	assert.NotNil(t, wrapped.Logger())
}

func TestInterfaces_newWailsApp_Ugly(t *testing.T) {
	wrapped := newWailsApp(nil)
	require.NotNil(t, wrapped)
	assert.Panics(t, func() {
		_ = wrapped.Logger()
	})
}
