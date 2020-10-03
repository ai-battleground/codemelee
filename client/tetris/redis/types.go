package tetris

import "time"

type Observation struct {
	State   string
	Error   error
	Score   int
	Lines   int
	Level   int
	Tick    int
	Elapsed time.Duration
	Board   [][]byte
	Shelf   []byte
}

type GameTime struct {
	Ticks int
	Time  time.Time
}
