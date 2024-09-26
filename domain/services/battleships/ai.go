package battleships

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	BoardSize  = 10
	EmptyCell  = 0
	ShipInCell = 1
	Hit        = 2
	Miss       = 3
)

type Point struct {
	x, y int
}

type Ship struct {
	name   string
	length int
}

type Board struct {
	grid [BoardSize][BoardSize]int
}

type AI struct {
	myBoard        *Board
	opponentBoard  *Board
	hitPoints      []Point // Tracks hit cells to switch to target mode
	mode           string  // Either "hunt" or "target"
	shotsChannel   chan Point
	turnChannel    chan bool
	attacksChannel chan Point
}

func NewAI(attacksChannel chan Point) (*AI, chan Point, chan bool) {
	myBoard := NewBoard()
	myBoard.PlaceAllShips()

	opponentBoard := NewBoard()

	ai := &AI{
		myBoard:        myBoard,
		opponentBoard:  opponentBoard,
		mode:           "hunt",
		shotsChannel:   make(chan Point),
		turnChannel:    make(chan bool),
		attacksChannel: attacksChannel,
	}

	return ai, ai.shotsChannel, ai.turnChannel
}

// List of ships with their sizes
var ships = []Ship{
	{"Carrier", 5},
	{"Battleship", 4},
	{"Cruiser", 3},
	{"Submarine", 3},
	{"Destroyer", 2},
}

// Initialize an empty board
func NewBoard() *Board {
	return &Board{}
}

// Place ships randomly on the board (same as previous function)
func (b *Board) PlaceAllShips() {
	rand.Seed(time.Now().UnixNano())

	for _, ship := range ships {
		placed := false
		for !placed {
			start := Point{rand.Intn(BoardSize), rand.Intn(BoardSize)}
			orientation := rand.Intn(2) // 0: horizontal, 1: vertical

			if b.CanPlaceShip(start, ship, orientation) {
				b.PlaceShip(start, ship, orientation)
				placed = true
			}
		}
	}
}

// Check if a ship can be placed at a given position
func (b *Board) CanPlaceShip(start Point, ship Ship, orientation int) bool {
	if orientation == 0 { // Horizontal
		if start.x+ship.length > BoardSize {
			return false
		}
		for i := 0; i < ship.length; i++ {
			if b.grid[start.y][start.x+i] != EmptyCell {
				return false
			}
		}
	} else { // Vertical
		if start.y+ship.length > BoardSize {
			return false
		}
		for i := 0; i < ship.length; i++ {
			if b.grid[start.y+i][start.x] != EmptyCell {
				return false
			}
		}
	}
	return true
}

// Place ship on the board
func (b *Board) PlaceShip(start Point, ship Ship, orientation int) {
	if orientation == 0 {
		for i := 0; i < ship.length; i++ {
			b.grid[start.y][start.x+i] = ShipInCell
		}
	} else {
		for i := 0; i < ship.length; i++ {
			b.grid[start.y+i][start.x] = ShipInCell
		}
	}
}

// AI makes a move by firing at the board
func (ai *AI) MakeMove() Point {
	if ai.mode == "hunt" {
		return ai.huntMode()
	}
	return ai.targetMode()
}

func (ai *AI) cellAlreadyTarget(x int, y int) bool {
	return ai.opponentBoard.grid[y][x] == Hit || ai.opponentBoard.grid[y][x] == Miss
}

// Hunt mode: fire in a checkerboard pattern
func (ai *AI) huntMode() Point {
	for {
		x, y := rand.Intn(BoardSize), rand.Intn(BoardSize)
		if (x+y)%2 == 0 && !ai.cellAlreadyTarget(x, y) { // Checkerboard pattern
			return Point{x, y}
		}
	}
}

// Target mode: focus on cells around the last hit
func (ai *AI) targetMode() Point {
	lastHit := ai.hitPoints[len(ai.hitPoints)-1]
	neighbors := []Point{
		{lastHit.x + 1, lastHit.y}, {lastHit.x - 1, lastHit.y}, // Left, Right
		{lastHit.x, lastHit.y + 1}, {lastHit.x, lastHit.y - 1}, // Up, Down
	}

	for _, neighbor := range neighbors {
		if ai.isValid(neighbor) && !ai.cellAlreadyTarget(neighbor.x, neighbor.y) {
			return neighbor
		}
	}

	// If no neighbors are valid, return to hunt mode
	ai.mode = "hunt"
	return ai.huntMode()
}

// Check if a point is within board limits
func (ai *AI) isValid(p Point) bool {
	return p.x >= 0 && p.x < BoardSize && p.y >= 0 && p.y < BoardSize
}

// AI attacks a specific point on the board
func (ai *AI) Attack(p Point) {
	if ai.opponentBoard.grid[p.y][p.x] == ShipInCell {
		ai.opponentBoard.grid[p.y][p.x] = Hit
		fmt.Printf("Hit at (%d, %d)\n", p.x, p.y)
		ai.hitPoints = append(ai.hitPoints, p)
		ai.mode = "target"
	} else {
		ai.opponentBoard.grid[p.y][p.x] = Miss
		fmt.Printf("Miss at (%d, %d)\n", p.x, p.y)
	}
}

func (ai *AI) Shot(p Point) int {
	if ai.myBoard.grid[p.y][p.x] == ShipInCell {
		return Hit
	} else {
		return Miss
	}
}

// Display the current board state
func (b *Board) Display() {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			switch b.grid[i][j] {
			case EmptyCell:
				fmt.Print(". ")
			case ShipInCell:
				fmt.Print("S ")
			case Hit:
				fmt.Print("X ")
			case Miss:
				fmt.Print("O ")
			}
		}
		fmt.Println()
	}
}

func (ai *AI) Start(ctx context.Context) {
	for {
		select {
		case p := <-ai.shotsChannel:
			ai.Shot(p)
		case <-ai.turnChannel:
			move := ai.MakeMove()
			ai.Attack(move)
			fmt.Println("Current Board:")
			ai.opponentBoard.Display()
		case <-ctx.Done():
			return
		}
	}
}

// Main function to simulate the AI playing Battleship
/*func main() {
	attacksChannel := make(chan Point)
	ai := NewAI(attacksChannel)

	fmt.Println("Initial Board:")
	ai.myBoard.Display()
}*/
