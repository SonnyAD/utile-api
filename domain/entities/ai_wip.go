package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	BoardSize  = 10
	Horizontal = 0
	Vertical   = 1
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

// Check if a ship can be placed on the board at a given position
func (b *Board) CanPlaceShip(start Point, ship Ship, orientation int) bool {
	if orientation == Horizontal {
		// Ensure ship fits horizontally
		if start.x+ship.length > BoardSize {
			return false
		}
		// Check if the ship overlaps with existing ships
		for i := 0; i < ship.length; i++ {
			if b.grid[start.y][start.x+i] != 0 {
				return false
			}
		}
	} else {
		// Ensure ship fits vertically
		if start.y+ship.length > BoardSize {
			return false
		}
		// Check if the ship overlaps with existing ships
		for i := 0; i < ship.length; i++ {
			if b.grid[start.y+i][start.x] != 0 {
				return false
			}
		}
	}
	return true
}

// Place a ship on the board at a valid position
func (b *Board) PlaceShip(start Point, ship Ship, orientation int) {
	if orientation == Horizontal {
		for i := 0; i < ship.length; i++ {
			b.grid[start.y][start.x+i] = 1
		}
	} else {
		for i := 0; i < ship.length; i++ {
			b.grid[start.y+i][start.x] = 1
		}
	}
}

// Randomly place all ships on the board
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

// Display the board
func (b *Board) Display() {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if b.grid[i][j] == 1 {
				fmt.Print("S ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

func main() {
	board := NewBoard()
	board.PlaceAllShips()
	board.Display()
}
