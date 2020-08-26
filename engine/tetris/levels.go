package tetris

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
		maxLines: 10,
		NextPiece: func() Piece {
			return Pieces.O
		},
		Score: func(lines int) int {
			return lines * lines
		},
	},
}

func getLevel(i int) Level {
	var l Level
	if i > len(levels) {
		l = levels[len(levels)-1]
	} else {
		l = levels[i]
	}
	l.Next = func() Level {
		return getLevel(i + 1)
	}
	return l
}
