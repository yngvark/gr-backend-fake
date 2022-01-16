package worldmap_test

import (
	"fmt"
	"testing"

	"github.com/yngvark/gr-zombie/pkg/worldmap"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestWorldMap(t *testing.T) {
	m := worldmap.New(1, 1)

	testCases := []struct {
		Val      int
		Axis     worldmap.AxisType
		Expected bool
	}{
		{
			Val:      -1,
			Axis:     worldmap.Axis.X,
			Expected: false,
		},
		{
			Val:      0,
			Axis:     worldmap.Axis.X,
			Expected: true,
		},
		{
			Val:      1,
			Axis:     worldmap.Axis.X,
			Expected: true,
		},
		{
			Val:      2,
			Axis:     worldmap.Axis.X,
			Expected: false,
		},

		{
			Val:      -1,
			Axis:     worldmap.Axis.Y,
			Expected: false,
		},
		{
			Val:      0,
			Axis:     worldmap.Axis.Y,
			Expected: true,
		},
		{
			Val:      1,
			Axis:     worldmap.Axis.Y,
			Expected: true,
		},
		{
			Val:      2,
			Axis:     worldmap.Axis.Y,
			Expected: false,
		},
	}
	for _, tc := range testCases {
		tc := tc

		name := "Should return correct withinmap val=%d, axis=%d -> %t"
		t.Run(fmt.Sprintf(name, tc.Val, tc.Axis, tc.Expected), func(t *testing.T) {
			actual, err := m.IsInMap(tc.Val, tc.Axis)
			assert.Nil(t, err)
			assert.Equal(t, tc.Expected, actual,
				fmt.Sprintf("Expected %t when checking %d inside [%d, %d]", tc.Expected, tc.Val, m.MaxX, m.MaxY))
		})
	}
}
