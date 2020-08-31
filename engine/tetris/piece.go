package tetris

import (
	"fmt"
)

type Piece struct {
	Name         string
	Orientations [][4]Point
}

func (piece Piece) Height() int {
	var low, high int
	for _, point := range piece.Orientations[0] {
		if point.Y < low {
			low = point.Y
		}
		if point.Y > high {
			high = point.Y
		}
	}
	return high - low + 1
}

type ActivePiece struct {
	Piece
	Position    Point
	Orientation int
}

type Point struct {
	X, Y int
}

func (point Point) String() string {
	return fmt.Sprintf("(%d, %d)", point.X, point.Y)
}

type Space struct {
	contents byte
}

func (space Space) String() string {
	return string(space.contents)
}

func (space Space) Empty() bool {
	return space.contents == ' '
}

func (active ActivePiece) Points() [4]Point {
	return active.Orientations[active.Orientation]
}

var Pieces = struct {
	O Piece
	I Piece
	T Piece
	L Piece
	J Piece
}{
	O: Piece{
		Name: "O",
		Orientations: [][4]Point{
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		}},
	I: Piece{
		Name: "I",
		Orientations: [][4]Point{
			{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		}},
	T: Piece{
		Name: "T",
		Orientations: [][4]Point{
			{{0, 0}, {1, 0}, {1, 1}, {2, 0}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 1}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 1}},
			{{0, 1}, {1, 0}, {1, 1}, {1, 2}},
		}},
	L: Piece{
		Name: "L",
		Orientations: [][4]Point{
			{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 0}},
			{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
			{{0, 2}, {1, 0}, {1, 1}, {1, 2}},
		}},
	J: Piece{
		Name: "J",
		Orientations: [][4]Point{
			{{0, 0}, {0, 1}, {1, 0}, {2, 0}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
			{{0, 1}, {1, 1}, {2, 0}, {2, 1}},
			{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
		}},
}
