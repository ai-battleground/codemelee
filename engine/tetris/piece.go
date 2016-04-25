package tetris

import (
	"fmt"
)

type TetrisPiece struct {
	Name         string
	Orientations [][4]Point
}

func (piece TetrisPiece) Height() int {
	var low, high int
	for _, point := range piece.Orientations[0] {
		if point.y < low {
			low = point.y
		}
		if point.y > high {
			high = point.y
		}
	}
	return high - low + 1
}

type ActivePiece struct {
	TetrisPiece
	Position    Point
	Orientation int
}

type Point struct {
	x, y int
}

func (point Point) String() string {
	return fmt.Sprintf("(%d, %d)", point.x, point.y)
}

type Space struct {
	empty bool
}

func (space Space) String() string {
	if space.empty {
		return " "
	} else {
		return "*"
	}
}

func (active ActivePiece) Points() [4]Point {
	return active.Orientations[active.Orientation]
}

var Pieces = struct {
	O TetrisPiece
	I TetrisPiece
	T TetrisPiece
}{
	TetrisPiece{
		Name:         "O",
		Orientations: [][4]Point{}},
	TetrisPiece{
		Name:         "I",
		Orientations: [][4]Point{}},
	TetrisPiece{
		Name:         "T",
		Orientations: [][4]Point{}}}

func init() {
	Pieces.O.Orientations = append(Pieces.O.Orientations,
		[4]Point{Point{0, 0}, Point{0, 1}, Point{1, 0}, Point{1, 1}})

	Pieces.I.Orientations = append(Pieces.I.Orientations,
		[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{0, 3}})
	Pieces.I.Orientations = append(Pieces.I.Orientations,
		[4]Point{Point{0, 0}, Point{1, 0}, Point{2, 0}, Point{3, 0}})

	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 0}, Point{1, 0}, Point{1, 1}, Point{2, 0}})
	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{1, 1}})
	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{2, 1}})
	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{1, 2}})
}
