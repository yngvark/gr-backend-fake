// Package zombie knows how to manage zombies
package zombie

import (
	"fmt"
	"math/rand"

	"github.com/yngvark/gr-zombie/pkg/worldmap"
)

// Zombie is a horrible monster
type Zombie struct {
	ID       string
	X        int
	Y        int
	WorldMap *worldmap.MapCreate
	Rand     *rand.Rand
}

// Move moves the Zombie
func (z *Zombie) Move() (*Zombie, *Move, error) {
	newX, err := z.getNewCoordPart(z.X, worldmap.Axis.X)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get new x coordinate: %w", err)
	}

	newY, err := z.getNewCoordPart(z.Y, worldmap.Axis.Y)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get new x coordinate: %w", err)
	}

	newZ := NewZombie(z.ID, newX, newY, z.WorldMap, z.Rand)
	move := NewZombieMove(z.ID, newX, newY)

	return newZ, move, nil
}

func (z *Zombie) getNewCoordPart(currentValue int, axisType worldmap.AxisType) (int, error) {
	direction := z.Rand.Intn(3) - 1 //nolint:gomnd,gosec    // [-1, 1]
	suggestion := currentValue + direction

	isInMap, err := z.WorldMap.IsInMap(suggestion, axisType)
	if err != nil {
		return -1, fmt.Errorf("could not detect if value is within map: %w", err)
	}

	if isInMap {
		return suggestion, nil
	}

	return currentValue, nil
}

// NewZombie returns a new Zombie
func NewZombie(id string, x int, y int, worldMap *worldmap.MapCreate, rnd *rand.Rand) *Zombie {
	return &Zombie{
		ID:       id,
		X:        x,
		Y:        y,
		WorldMap: worldMap,
		Rand:     rnd,
	}
}
