package tetris

func (b *Board) MoveRight() {
	b.move(Point{1, 0})
}

func (b *Board) MoveLeft() {
	b.move(Point{-1, 0})
}

func (b *Board) RotateLeft() {
	b.rotate(b.rotateCounterclockwise())
}

func (b *Board) RotateRight() {
	b.rotate(b.rotateClockwise())
}

func (b *Board) rotate(targetOrientation int) {
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

func (b *Board) rotateClockwise() int {
	return (b.Active.Orientation + 1) % len(b.Active.Orientations)
}

func (b *Board) rotateCounterclockwise() int {
	if b.Active.Orientation == 0 {
		return len(b.Active.Orientations) - 1
	}
	return b.Active.Orientation - 1
}

func (b *Board) HardDrop() {
	for !b.shouldAnchor() {
		b.move(Point{0, -1})
	}
	b.anchor()
}

func (b *Board) SoftDrop() {
	b.Advance()
}
