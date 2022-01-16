// Package worldmap knows how to handle world maps
package worldmap

import (
	"fmt"
)

// AxisType defines different kinds of axis types
type AxisType int

// Axis contains the values for different kinds of axis types
var Axis = struct { //nolint:gochecknoglobals
	X AxisType
	Y AxisType
}{
	X: 1, //nolint:gomnd
	Y: 2, //nolint:gomnd
}

// WorldMap is a world map
type WorldMap struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

// IsInMap returns whether a point on the axis of a given type, is within the map
func (m *WorldMap) IsInMap(axisValue int, axis AxisType) (bool, error) {
	switch axis {
	case Axis.X:
		return m.xIsInMap(axisValue), nil
	case Axis.Y:
		return m.yIsInMap(axisValue), nil
	default:
		return false, fmt.Errorf("not a valid Axis type: %d", axis)
	}
}

func (m *WorldMap) xIsInMap(x int) bool {
	return x >= m.MinX && x <= m.MaxX
}

func (m *WorldMap) yIsInMap(y int) bool {
	return y >= m.MinY && y <= m.MaxY
}

// New returns a new WorldMap
func New(maxX int, maxY int) *WorldMap {
	return &WorldMap{
		MinX: 0,
		MaxX: maxX,
		MinY: 0,
		MaxY: maxY,
	}
}
