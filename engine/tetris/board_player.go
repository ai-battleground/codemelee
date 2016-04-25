package tetris

func (board *Board) MoveRight() {
	board.move(Point{1, 0})
}

func (board *Board) MoveLeft() {
	board.move(Point{-1, 0})
}

func (board *Board) Rotate() {
	board.Active.Orientation = (board.Active.Orientation + 1) % len(board.Active.Orientations)
}

func (board *Board) Drop() {
	for !board.shouldAnchor() {
		board.move(Point{0, -1})
	}
	board.anchor()
}
