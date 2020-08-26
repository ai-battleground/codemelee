package tetris

import (
	"fmt"
	"time"
)

type GameState int

const (
	PreStart GameState = iota
	Running
	Paused
	GameOver
)

func (gs GameState) String() string {
	switch gs {
	case PreStart:
		return "Not Started"
	case Running:
		return "Running"
	case Paused:
		return "Paused"
	case GameOver:
		return "Game Over"
	default:
		return ""
	}
}

type Game struct {
	score int
	state GameState
	lines int
	shelf
	Level
	*Board
	ShelfUpdated chan [4]Piece
	PieceState   chan ActivePiece
}

func NewGame() *Game {
	g := &Game{
		score:        0,
		lines:        0,
		Level:        getLevel(0),
		Board:        NewBoard(),
		ShelfUpdated: make(chan [4]Piece),
		PieceState:   make(chan ActivePiece)}
	go g.handleAnchored()
	g.Board.OnClearedLines(g.handleClearedLines)

	return g
}

type Level struct {
	number    int
	speed     int64
	maxLines  int
	Next      func() Level
	NextPiece func() Piece
	Score     func(lines int) int
}

type shelf struct {
	pieces [4]Piece
	head   int
}

func (g *Game) Start() {
	if g.state == GameOver {
		return
	}
	if g.state == PreStart {
		first := g.NextPiece()
		for i := 0; i < 4; i++ {
			g.shelf.push(g.NextPiece())
		}
		g.Stage(first)
	}
	if g.state == Paused || g.state == PreStart {
		time.AfterFunc(time.Second/time.Duration(g.Level.speed), g.advance)
	}
	g.state = Running
}

func (g *Game) Pause() {
	g.state = Paused
}

func (g Game) Score() int {
	return g.score
}

func (g Game) State() GameState {
	return g.state
}

func (s *shelf) Shelf() [4]Piece {
	return [4]Piece{
		s.pieces[s.head],
		s.pieces[(s.head+1)%4],
		s.pieces[(s.head+2)%4],
		s.pieces[(s.head+3)%4]}
}

func (s *shelf) push(p Piece) {
	s.pieces[s.head] = p
	s.head = (s.head + 1) % 4
}

func (s *shelf) next() Piece {
	return s.pieces[s.head]
}

func (g *Game) advance() {
	if g.state == Running {
		g.Board.Advance()
		go func() { g.PieceState <- *g.Active }()
		time.AfterFunc(time.Second/time.Duration(g.Level.speed), g.advance)
	}
}

func (g *Game) handleAnchored() {
	for {
		_ = <-g.Board.Anchored
		g.advancePiece()
	}
}

func (g *Game) handleClearedLines(lines []int) {
	g.lines += len(lines)
	g.score += g.Level.Score(len(lines))
	fmt.Printf("Total lines cleared %d", g.lines)
	if g.lines >= g.Level.maxLines {
		g.Level = g.Level.Next()
	}
}

func (g *Game) advancePiece() {
	g.Stage(g.shelf.next())
	g.shelf.push(g.Level.NextPiece())
	go func() { g.ShelfUpdated <- g.Shelf() }()
}
