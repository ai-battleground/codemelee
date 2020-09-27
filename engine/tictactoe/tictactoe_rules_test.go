package tictactoe

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"testing"
)

func TestTicTacToeRulesFixture(t *testing.T) {
	gunit.Run(new(TicTacToeRulesFixture), t)
}

type TicTacToeRulesFixture struct {
	*gunit.Fixture
	normalBoard, giantBoard *Board
}

func (this *TicTacToeRulesFixture) SetupBoards() {
	this.normalBoard = EmptyBoard(3, 3)
	this.giantBoard = EmptyBoard(13, 17)
}

func (this *TicTacToeRulesFixture) TestInvalidMoveGivesError() {
	normalBoardNonsense := [][2]int{{-1, 1}, {1, -1}, {0, 3}, {3, 0}, {1000, 1000}}
	giantBoardNonsense := [][2]int{{-1, 0}, {0, -1}, {17, 0}, {0, 13}, {1000, 1000}}
	for _, testdata := range normalBoardNonsense {
		err := this.normalBoard.X(testdata[0], testdata[1])
		this.So(err, should.BeError, impossibleMove(3, 3))
	}
	for _, testdata := range giantBoardNonsense {
		err := this.giantBoard.O(testdata[0], testdata[1])
		this.So(err, should.BeError, impossibleMove(13, 17))
	}
}

func (this *TicTacToeRulesFixture) TestMoveMustNotBOccupied() {
	this.normalBoard.X(1, 1)
	err := this.normalBoard.O(1, 1)
	this.So(err, should.BeError, spaceIsOccupied(1, 1))
	this.normalBoard.O(0, 0)
	err = this.normalBoard.X(0, 0)
	this.So(err, should.BeError, spaceIsOccupied(0, 0))
}
