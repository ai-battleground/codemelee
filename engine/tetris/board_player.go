package tetris

func (b *Board) MoveRight() {
	b.move(Point{1, 0})
}

func (b *Board) MoveLeft() {
	b.move(Point{-1, 0})
}

func (b *Board) Rotate() {
	targetOrientation := (b.Active.Orientation + 1) % len(b.Active.Orientations)
	targetPosition := b.Active.Position

	for b.anyPointsCollide(targetPosition, b.Active.Orientations[targetOrientation]) &&
		targetPosition.X >= 0 {
		targetPosition.X--
	}
	if targetPosition.X >= 0 {
		b.Active.Orientation = targetOrientation
		b.Active.Position = targetPosition
	}
}

func (b *Board) Drop() {
	for !b.shouldAnchor() {
		b.move(Point{0, -1})
	}
	b.anchor()
}
