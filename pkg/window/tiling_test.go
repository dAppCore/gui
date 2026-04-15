// pkg/window/tiling_test.go
package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTiling_SnapPosition_String_Good(t *testing.T) {
	tests := []struct {
		name string
		pos  SnapPosition
		want string
	}{
		{name: "left", pos: SnapLeft, want: "left"},
		{name: "right", pos: SnapRight, want: "right"},
		{name: "center", pos: SnapCenter, want: "center"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.pos.String())
		})
	}
}

func TestTiling_SnapPosition_String_Bad(t *testing.T) {
	assert.Empty(t, SnapPosition(123).String())
}

func TestTiling_SnapPosition_String_Ugly(t *testing.T) {
	assert.Empty(t, SnapPosition(-1).String())
}
