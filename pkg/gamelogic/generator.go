// Package gamelogic contains the core game logic
package gamelogic

import (
	"fmt"

	zombiePkg "github.com/yngvark/gr-zombie/pkg/zombie"
)

// Generator knows how to generate zombie actions
type Generator struct {
	zombie *zombiePkg.Zombie
}

// Next returns the next zombie action
func (g *Generator) Next() (*zombiePkg.Move, error) {
	z, move, err := g.zombie.Move()
	if err != nil {
		return nil, fmt.Errorf("could not move zombie: %w", err)
	}

	g.zombie = z

	return move, nil
}

// NewGenerator returns a new Generator
func NewGenerator(initialZombie *zombiePkg.Zombie) *Generator {
	return &Generator{
		zombie: initialZombie,
	}
}
