package tictactoe_test

import (
	"github.com/ai-battleground/codemelee/engine/tictactoe"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"testing"
)

func TestTicTacToeFixture(t *testing.T) {
	gunit.Run(new(TicTacToeFixture), t)
}

type TicTacToeFixture struct {
	*gunit.Fixture
	normalBoard, giantBoard *tictactoe.Board
}

func (this *TicTacToeFixture) SetupBoards() {
	this.normalBoard = tictactoe.EmptyBoard(3, 3)
	this.giantBoard = tictactoe.EmptyBoard(13, 17)
}

func (this *TicTacToeFixture) TestBoardRows() {
	this.So(this.normalBoard.Height(), should.Equal, 3)
	this.So(this.giantBoard.Height(), should.Equal, 13)
}

func (this *TicTacToeFixture) TestBoardCols() {
	this.So(this.normalBoard.Width(), should.Equal, 3)
	this.So(this.giantBoard.Width(), should.Equal, 17)
}

func (this *TicTacToeFixture) TestEmptyBoardCellsAreEmpty() {
	for r := 0; r < this.normalBoard.Height(); r++ {
		for c := 0; c < this.normalBoard.Width(); c++ {
			this.So(this.normalBoard.Cell(r, c), should.Equal, tictactoe.CellStateEmpty)
		}
	}
}

func (this *TicTacToeFixture) TestNextMoveIsXWhenEmpty() {
	this.So(this.normalBoard.NextMove(), should.Equal, tictactoe.CellStateX)
	this.So(this.giantBoard.NextMove(), should.Equal, tictactoe.CellStateX)
}

func (this *TicTacToeFixture) TestNextMoveChangesProperlyAfterMoves() {
	this.normalBoard.X(0, 0)
	this.giantBoard.X(0, 0)
	this.So(this.normalBoard.NextMove(), should.Equal, tictactoe.CellStateO)
	this.So(this.giantBoard.NextMove(), should.Equal, tictactoe.CellStateO)
	this.normalBoard.O(0, 1)
	this.giantBoard.O(0, 1)
	this.So(this.normalBoard.NextMove(), should.Equal, tictactoe.CellStateX)
	this.So(this.giantBoard.NextMove(), should.Equal, tictactoe.CellStateX)
}

func (this *TicTacToeFixture) TestMoveXWorks() {
	errNormal := this.normalBoard.X(1, 2)
	errGiant := this.giantBoard.X(5, 6)
	this.So(errNormal, should.BeNil)
	this.So(this.normalBoard.Cell(1, 2), should.Equal, tictactoe.CellStateX)
	this.So(errGiant, should.BeNil)
	this.So(this.giantBoard.Cell(5, 6), should.Equal, tictactoe.CellStateX)
}

func (this *TicTacToeFixture) TestMoveOWorks() {
	// X goes first
	this.normalBoard.X(0, 0)
	this.giantBoard.X(0, 0)
	errNormal := this.normalBoard.O(2, 1)
	errGiant := this.giantBoard.O(6, 5)
	this.So(errNormal, should.BeNil)
	this.So(this.normalBoard.Cell(2, 1), should.Equal, tictactoe.CellStateO)
	this.So(errGiant, should.BeNil)
	this.So(this.giantBoard.Cell(6, 5), should.Equal, tictactoe.CellStateO)
}

func (this *TicTacToeFixture) TestXGoesFirst() {
	err := this.normalBoard.O(0, 0)
	this.So(err, should.BeError, tictactoe.MoveOutOfTurn)
	err = this.normalBoard.X(2, 2)
	this.So(err, should.BeNil)
	err = this.normalBoard.O(0, 0)
	this.So(err, should.BeNil)
}

func (this *TicTacToeFixture) TestOneTurnAtATime() {
	err := this.normalBoard.X(0, 0)
	err = this.normalBoard.X(1, 1)
	this.So(err, should.BeError, tictactoe.MoveOutOfTurn)
	err = this.normalBoard.O(1, 1)
	err = this.normalBoard.O(0, 2)
	this.So(err, should.BeError, tictactoe.MoveOutOfTurn)
}

func (this *TicTacToeFixture) TestGameOver_3x3_NextMoveIsEmpty() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(0, 1)
	this.normalBoard.X(1, 0)
	this.normalBoard.O(0, 2)
	this.normalBoard.X(2, 0) // three in a row
	this.So(this.normalBoard.NextMove(), should.Equal, tictactoe.CellStateEmpty)
}

func (this *TicTacToeFixture) TestXWins_3x3_Horizontally() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(0, 1)
	this.normalBoard.X(1, 0)
	this.normalBoard.O(0, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Undetermined)
	this.normalBoard.X(2, 0) // three in a row
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.WonByX)
}

func (this *TicTacToeFixture) TestOWins_3x3_Horizontally() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(0, 1)
	this.normalBoard.X(1, 0)
	this.normalBoard.O(1, 1)
	this.normalBoard.X(0, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Undetermined)
	this.normalBoard.O(2, 1) // three in a row
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.WonByO)
}

func (this *TicTacToeFixture) TestXWins_3x3_Vertically() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(2, 2)
	this.normalBoard.X(0, 1)
	this.normalBoard.O(2, 1)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Undetermined)
	this.normalBoard.X(0, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.WonByX)
}

func (this *TicTacToeFixture) TestOWins_3x3_Vertically() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(2, 2)
	this.normalBoard.X(0, 1)
	this.normalBoard.O(2, 1)
	this.normalBoard.X(1, 1)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Undetermined)
	this.normalBoard.O(2, 0)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.WonByO)
}

func (this *TicTacToeFixture) TestXWins_3x3_DescendingDiagonal() {
	this.normalBoard.X(0, 2)
	this.normalBoard.O(0, 0)
	this.normalBoard.X(1, 1)
	this.normalBoard.O(2, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Undetermined)
	this.normalBoard.X(2, 0)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.WonByX)
}

func (this *TicTacToeFixture) TestOWins_3x3_AscendingDiagonal() {
	this.normalBoard.X(0, 2)
	this.normalBoard.O(0, 0)
	this.normalBoard.X(2, 0)
	this.normalBoard.O(1, 1)
	this.normalBoard.X(1, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Undetermined)
	this.normalBoard.O(2, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.WonByO)
}

func (this *TicTacToeFixture) TestDraw_3x3() {
	this.normalBoard.X(0, 0)
	this.normalBoard.O(1, 1)
	this.normalBoard.X(2, 2)
	this.normalBoard.O(0, 1)
	this.normalBoard.X(2, 1) // block forced from here on
	this.normalBoard.O(2, 0)
	this.normalBoard.X(0, 2)
	this.normalBoard.O(1, 2)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Undetermined)
	this.normalBoard.X(1, 0)
	this.So(this.normalBoard.GameOutcome(), should.Equal, tictactoe.Draw)
}
