package tetris

import (
	"strings"
	"sync"
)

type Board struct {
	width, height      int
	plane              [][]Space
	Active             *ActivePiece
	Anchored           chan Piece
	clearedObservers   []func([]int)
	collisionObservers []func()
}

func NewBoard() *Board {
	b := &Board{
		width:    10,
		height:   20,
		plane:    newPlane(10, 20),
		Active:   &ActivePiece{},
		Anchored: make(chan Piece)}
	return b
}

func (b *Board) Advance() {
	if b.shouldAnchor() {
		b.anchor()
	} else {
		b.move(Point{0, -1})
	}
}

func (b *Board) Stage(piece Piece) {
	stagePosition := Point{4, b.height - piece.Height()}

	b.Active = &ActivePiece{piece, stagePosition, 0}

	if b.anyPointsCollide(stagePosition, piece.Orientations[0]) {
		var wg sync.WaitGroup
		wg.Add(len(b.collisionObservers))
		for _, o := range b.collisionObservers {
			observer := o
			go func() {
				observer()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func (b *Board) TakeSnapshot() string {
	var snapshotLines []string
	for i, _ := range b.plane {
		row := b.plane[b.height-1-i]
		var line []byte
		for _, space := range row {
			line = append(line, space.contents)
		}
		snapshotLines = append(snapshotLines, string(line))
	}
	if len(b.Active.Name) > 0 {
		for _, point := range b.Active.Points() {
			boardPoint := translate(b.Active.Position, point)
			line := []byte(snapshotLines[b.height-1-boardPoint.Y])
			line[boardPoint.X] = byte(b.Active.Name[0])
			snapshotLines[b.height-1-boardPoint.Y] = string(line)
		}
	}
	return strings.Join(snapshotLines, "\n")
}

func (b *Board) OnClearedLines(observer func([]int)) {
	b.clearedObservers = append(b.clearedObservers, observer)
}

func (b *Board) OnCollision(observer func()) {
	b.collisionObservers = append(b.collisionObservers, observer)
}

func (b *Board) anchor() {
	for _, p := range b.Active.Points() {
		b.space(translate(b.Active.Position, p)).contents = b.Active.Name[0]
	}
	go func() { b.Anchored <- b.Active.Piece }()
	b.clearLines()
}

func (b *Board) clearLines() {
	var completedLines []int
	for i, row := range b.plane {
		if isComplete(row) {
			completedLines = append(completedLines, i)
		}
	}
	if len(completedLines) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(b.clearedObservers))
		for _, o := range b.clearedObservers {
			observer := o
			go func() {
				observer(completedLines)
				wg.Done()
			}()
		}
		wg.Wait()
		// iterate in reverse so that the row indices stay accurate
		for i := len(completedLines) - 1; i >= 0; i-- {
			b.clearLine(completedLines[i])
		}
	}
}

func (b *Board) move(vector Point) {
	if !b.wouldCollide(vector) {
		destination := translate(b.Active.Position, vector)
		b.Active.Position = destination
	}
}

func (b Board) shouldAnchor() bool {
	return b.wouldCollide(Point{0, -1})
}

func (b Board) wouldCollide(vector Point) bool {
	position := translate(b.Active.Position, vector)
	return b.anyPointsCollide(position, b.Active.Points())
}

func (b Board) anyPointsCollide(position Point, points [4]Point) bool {
	for _, p := range points {
		testPoint := translate(position, p)
		if testPoint.Y < 20 {
			if testPoint.Y < 0 || testPoint.X < 0 || testPoint.X >= 10 {
				return true
			}
			if !b.space(testPoint).Empty() {
				return true
			}
		}
	}
	return false
}

func (b *Board) space(point Point) *Space {
	return &b.plane[point.Y][point.X]
}

func (b *Board) clearLine(y int) {
	for i := y; i < 19; i++ {
		for j := range b.plane[i] {
			b.plane[i][j].contents = b.plane[i+1][j].contents
		}
	}
	for x := range b.plane[19] {
		b.plane[19][x].contents = ' '
	}
}

func translate(origin Point, vector Point) Point {
	return Point{origin.X + vector.X, origin.Y + vector.Y}
}

func isComplete(row []Space) bool {
	for _, space := range row {
		if space.Empty() {
			return false
		}
	}
	return true
}

func newPlane(width int, height int) [][]Space {
	plane := make([][]Space, height)
	for i := range plane {
		plane[i] = newRow(width)
	}
	return plane
}

func newRow(width int) []Space {
	row := make([]Space, width)
	for i := range row {
		row[i] = Space{contents: ' '}
	}
	return row
}
