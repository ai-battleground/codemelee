package tetris

import (
	"math/rand"
)

type Level struct {
	Number    int
	speed     float32
	maxLines  int
	Next      func() Level
	NextPiece func() Piece
	Score     func(lines int) int
}

func pieces() []Piece {
	return []Piece{
		Pieces.O, Pieces.I, Pieces.T, Pieces.J, Pieces.L, Pieces.S, Pieces.Z,
	}
}

var levels = [...]Level{
	{
		Number:   1,
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
		Number:   2,
		speed:    1,
		maxLines: 20,
		NextPiece: func() Piece {
			index := rand.Intn(2)
			return pieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 1
			return base * base
		},
	},
	{
		Number:   3,
		speed:    1,
		maxLines: 30,
		NextPiece: func() Piece {
			index := rand.Intn(3)
			return pieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 2
			return base * base
		},
	},
	{
		Number:   4,
		speed:    1,
		maxLines: 40,
		NextPiece: func() Piece {
			index := rand.Intn(5)
			return pieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 3
			return base * base
		},
	},
	{
		Number:   5,
		speed:    1.05,
		maxLines: 50,
		NextPiece: func() Piece {
			index := rand.Intn(5)
			return pieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 4
			return base * base
		},
	},
	{
		Number:   6,
		speed:    1.1,
		maxLines: 60,
		NextPiece: func() Piece {
			index := rand.Intn(7)
			return pieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 5
			return base * base
		},
	},
	{
		Number:   7,
		speed:    1.15,
		maxLines: 70,
		NextPiece: func() Piece {
			index := rand.Intn(7)
			return pieces()[index]
		},
		Score: func(lines int) int {
			base := lines + 6
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
