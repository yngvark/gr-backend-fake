package zombie

// Move is a move in the world
type Move struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// NewZombieMove returns a new Move
func NewZombieMove(id string, x int, y int) *Move {
	return &Move{
		Type: "zombieMove",
		ID:   id,
		X:    x,
		Y:    y,
	}
}
