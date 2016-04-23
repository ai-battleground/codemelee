package tetris

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTetrisGame(t *testing.T) {
	Convey("Given a tetris game", t, func() {
		game := NewTetrisGame()

		Convey("the initial level should be 1", func() {
			So(game.Level.number, ShouldEqual, 1)
		})

		Convey("the initial speed should be 1", func() {
			So(game.Level.speed, ShouldEqual, 1)
		})

		Convey("the initial score should be 0", func() {
			So(game.score, ShouldEqual, 0)
		})

		Convey("the board", func() {
			board := game.Board

			Convey("dimensions should be 10x20", func() {
				So(board.width, ShouldEqual, 10)
				So(board.height, ShouldEqual, 20)
			})

			Convey("should be empty", func() {
				for _, row := range board.plane {
					for _, space := range row {
						So(space.empty, ShouldBeTrue)
					}
				}
			})
		})

		Convey("when the game is started", func() {
			game.Start()

			Convey("the board", func() {
				board := game.Board

				Convey("should have an active piece", func() {
					So(board.Active.Piece, ShouldEqual, Pieces.O)
				})

				Convey("should position the piece at the top", func() {
					So(board.Active.Position.y, ShouldEqual, board.height-board.Active.Piece.Height())
				})
			})

		})

		Convey("when a piece is anchored", func() {
			game.Start()

			game.Board.Anchored <- Pieces.O

			Convey("a new piece should be staged", func() {
				So(game.Board.Active.Piece, ShouldEqual, Pieces.O)
				So(game.Board.Active.Position.y, ShouldEqual, game.Board.height-game.Board.Active.Piece.Height())
			})
		})
	})
}
