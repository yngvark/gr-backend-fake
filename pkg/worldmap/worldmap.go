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
	Type  string  `json:"type,omitempty"`
	MinX  int     `json:"minX,omitempty"`
	MaxX  int     `json:"maxX,omitempty"`
	MinY  int     `json:"minY,omitempty"`
	MaxY  int     `json:"maxY,omitempty"`
	Tiles [][]int `json:"tiles,omitempty"`
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
	tiles := generateTiles(maxX, maxY)

	return &WorldMap{
		Type:  "mapCreate",
		MinX:  0,
		MaxX:  maxX,
		MinY:  0,
		MaxY:  maxY,
		Tiles: tiles,
	}
}

func generateTiles(maxX, maxY int) [][]int {
	ys := make([][]int, maxY)

	for y := 0; y < maxY; y++ {
		xs := make([]int, maxX)
		ys[y] = xs

		for x := 0; x < maxX; x++ {
			xs[x] = y*maxX + x
		}
	}

	return ys
}
