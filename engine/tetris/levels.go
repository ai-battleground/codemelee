package tetris

import "math/rand"

func tutorialPieces() []Piece {
	return []Piece{
		Pieces.O, Pieces.I, Pieces.T,
	}
}

var levels = [...]Level{
	{
		number:   1,
		speed:    1,
		maxLines: 10,
		NextPiece: func() Piece {
			return Pieces.O
		},
		Score: func(lines int) int {
			return lines * lines
		},
	},
	{
		number:   2,
		speed:    1,
		maxLines: 20,
		NextPiece: func() Piece {
			index := rand.Intn(2)
			return tutorialPieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 1
			return base * base
		},
	},
	{
		number:   3,
		speed:    1,
		maxLines: 30,
		NextPiece: func() Piece {
			index := rand.Intn(3)
			return tutorialPieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 2
			return base * base
		},
	}}

func getLevel(i int) Level {
	var l Level
	if i > len(levels)-1 {
		l = levels[len(levels)-1]
	} else {
		l = levels[i]
	}
	l.Next = func() Level {
		return getLevel(i + 1)
	}
	return l
}
