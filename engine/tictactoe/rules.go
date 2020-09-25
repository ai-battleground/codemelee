package tictactoe

import "fmt"

const (
	NoProblem     = RuleViolation("OK")
	MoveOutOfTurn = RuleViolation("Moved out of turn")
)

type RuleViolation string
type Rule func() (RuleViolation, bool)

func (v RuleViolation) Error() string {
	return string(v)
}

func ImpossibleMove(rows, cols int) RuleViolation {
	return RuleViolation(fmt.Sprintf("Move not possible: Board is %dx%d", rows, cols))
}

func SpaceIsOccupied(x, y int) RuleViolation {
	return RuleViolation(fmt.Sprintf("Space is already taken (%d, %d)", x, y))
}

func OK() (RuleViolation, bool) {
	return NoProblem, true
}

func (b Board) NextMove() CellState {
	moves := 0
	for _, row := range b.cells {
		for _, cell := range row {
			if cell != CellStateEmpty {
				moves++
			}
		}
	}
	if moves%2 == 1 {
		return CellStateO
	} else {
		return CellStateX
	}
}

func (b Board) CheckMoveIsPossible(x, y int) Rule {
	return func() (RuleViolation, bool) {
		if x < 0 || y < 0 {
			return ImpossibleMove(b.Height(), b.Width()), false
		}
		if y >= b.Height() || x >= b.Width() {
			return ImpossibleMove(b.Height(), b.Width()), false
		}
		return OK()
	}
}

func (b Board) CheckXGoesFirst(state CellState) Rule {
	return func() (RuleViolation, bool) {
		if !b.inProgress && state == CellStateO {
			return MoveOutOfTurn, false
		}
		return OK()
	}
}

func (b Board) CheckTakingTurns(state CellState) Rule {
	var turnIsO bool = false
	for i := 0; i < b.Width(); i++ {
		for j := 0; j < b.Height(); j++ {
			if b.Cell(i, j) != CellStateEmpty {
				turnIsO = !turnIsO
			}
		}
	}
	return func() (RuleViolation, bool) {
		if turnIsO {
			if state == CellStateX {
				return MoveOutOfTurn, false
			}
		} else {
			if state == CellStateO {
				return MoveOutOfTurn, false
			}
		}
		return OK()
	}
}

func (b Board) CheckUnoccupied(x, y int) Rule {
	return func() (RuleViolation, bool) {
		if b.Cell(x, y) != CellStateEmpty {
			return SpaceIsOccupied(x, y), false
		}
		return OK()
	}
}
