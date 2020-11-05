package tictactoe

import "time"

type ChallengeState struct {
	Active    bool
	Match     string
	Confirmed bool
	Started   bool
}

type Observation struct {
	Boards        [][]byte
	MyTurn        bool
	State         string
	Error         error
	Score         int
	OpponentScore int
	Bot           string
	Opponent      string
	Round         int
	MoveExpires   time.Time
	MatchExpires  time.Time
}

type GameTime struct {
	Ticks int
	Time  time.Time
}
