package tetris

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
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
					So(board.Active.TetrisPiece, ShouldResemble, Pieces.O)
				})

				Convey("should position the piece at the top", func() {
					So(board.Active.Position.y, ShouldEqual, board.height-board.Active.Height())
				})
			})

			Convey("the shelf", func() {
				shelf := game.Shelf()

				Convey("should load 4 O pieces", func() {
					for i := 0; i < 4; i++ {
						So(shelf[i], ShouldResemble, Pieces.O)
					}
				})
			})
		})

		Convey("when a piece is anchored", func() {
			game.Start()

			game.Board.Anchored <- Pieces.O

			Convey("a new piece should be staged", func() {
				So(game.Board.Active.TetrisPiece, ShouldResemble, Pieces.O)
				So(game.Board.Active.Position.y, ShouldEqual, game.Board.height-game.Board.Active.Height())
			})
		})

		Convey("the shelf", func() {

			Convey("should load 4 pieces from the Level when the game starts", func() {
				piece := 0
				game.Level = Level{0, 0, (func() TetrisPiece {
					defer func() { piece++ }()
					return TetrisPiece{Name: fmt.Sprintf("Piece %d", piece), Orientations: Pieces.O.Orientations}
				})}

				game.Start()
				shelf := game.Shelf()

				for i := range shelf {
					So(shelf[i].Name, ShouldEqual, fmt.Sprintf("Piece %d", i + 1)) // Piece 0 was already staged
				}
			})

			Convey("should load another piece from the Level when a piece is anchored", func() {
				game.Start()

				game.shelf.pieces = [4]TetrisPiece{
					TetrisPiece{Name: "A"},
					TetrisPiece{Name: "B"},
					TetrisPiece{Name: "C"},
					TetrisPiece{Name: "D"}}

				game.Level = Level{0, 0, func() TetrisPiece {
					return TetrisPiece{Name: "X"}
				}}

				go func() { game.Board.Anchored <- TetrisPiece{Name: "A", Orientations: Pieces.O.Orientations} }()

				select {
				case s := <-game.ShelfUpdated:
					So(s, ShouldResemble, [4]TetrisPiece{
						TetrisPiece{Name: "B"},
						TetrisPiece{Name: "C"},
						TetrisPiece{Name: "D"},
						TetrisPiece{Name: "X"}})
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}

				So(game.Shelf(), ShouldResemble, [4]TetrisPiece{
					TetrisPiece{Name: "B"},
					TetrisPiece{Name: "C"},
					TetrisPiece{Name: "D"},
					TetrisPiece{Name: "X"}})
			})
		})
	})
}
