package zombie_test

import (
	"math/rand"
	"testing"

	"github.com/yngvark/gr-zombie/pkg/gamelogic"

	"github.com/yngvark/gr-zombie/pkg/worldmap"
	"github.com/yngvark/gr-zombie/pkg/zombie"

	"github.com/stretchr/testify/assert"
)

func TestZombie(t *testing.T) {
	t.Run("Should move randomly", func(t *testing.T) {
		// Given
		m := worldmap.New(20, 10)                                          //nolint:gomnd
		z := zombie.NewZombie("1", 10, 5, m, rand.New(rand.NewSource(45))) //nolint:gosec,gomnd
		generator := gamelogic.NewGenerator(z)

		// When+Then
		assertNextPosition(t, generator, 9, 5)
		assertNextPosition(t, generator, 8, 4)
		assertNextPosition(t, generator, 9, 4)
		assertNextPosition(t, generator, 8, 5)
		assertNextPosition(t, generator, 8, 5)
	})
}

func assertNextPosition(t *testing.T, generator *gamelogic.Generator, x int, y int) {
	move, err := generator.Next()
	assert.Nil(t, err)
	assert.Equal(t, zombie.NewZombieMove("1", x, y), move)
}
