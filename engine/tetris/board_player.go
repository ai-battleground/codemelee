package tetris

func (board *Board) MoveRight() {
	board.move(Point{1, 0})
}

func (board *Board) MoveLeft() {
	board.move(Point{-1, 0})
}

func (board *Board) Rotate() {
	targetOrientation := (board.Active.Orientation + 1) % len(board.Active.Orientations)
	targetPosition := board.Active.Position

	for board.anyPointsCollide(targetPosition, board.Active.Orientations[targetOrientation]) &&
		targetPosition.x >= 0 {
		targetPosition.x--
	}
	if targetPosition.x >= 0 {
		board.Active.Orientation = targetOrientation
		board.Active.Position = targetPosition
	}
}

func (board *Board) Drop() {
	for !board.shouldAnchor() {
		board.move(Point{0, -1})
	}
	board.anchor()
}
