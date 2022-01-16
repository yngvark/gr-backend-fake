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

// MapCreate is a world map
type MapCreate struct {
	Type string `json:"type,omitempty"`
	MinX int    `json:"minX,omitempty"`
	MaxX int    `json:"maxX,omitempty"`
	MinY int    `json:"minY,omitempty"`
	MaxY int    `json:"maxY,omitempty"`
}

// IsInMap returns whether a point on the axis of a given type, is within the map
func (m *MapCreate) IsInMap(axisValue int, axis AxisType) (bool, error) {
	switch axis {
	case Axis.X:
		return m.xIsInMap(axisValue), nil
	case Axis.Y:
		return m.yIsInMap(axisValue), nil
	default:
		return false, fmt.Errorf("not a valid Axis type: %d", axis)
	}
}

func (m *MapCreate) xIsInMap(x int) bool {
	return x >= m.MinX && x <= m.MaxX
}

func (m *MapCreate) yIsInMap(y int) bool {
	return y >= m.MinY && y <= m.MaxY
}

// New returns a new MapCreate
func New(maxX int, maxY int) *MapCreate {
	return &MapCreate{
		Type: "mapCreate",
		MinX: 0,
		MaxX: maxX,
		MinY: 0,
		MaxY: maxY,
	}
}
