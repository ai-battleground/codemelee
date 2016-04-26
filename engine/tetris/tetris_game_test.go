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

			game.shelf.pieces = makeShelf("A", "B", "C", "D")

			game.Level = Level{0, 0, func() TetrisPiece {
				return TetrisPiece{Name: "X", Orientations: Pieces.O.Orientations}
			}, func(int) int { return 0 }}

			game.Board.Anchored <- Pieces.I

			Convey("a new piece should be staged from the shelf", func() {
				So(game.Board.Active.Name, ShouldEqual, "A")
				So(game.Board.Active.Position.y, ShouldEqual, game.Board.height-game.Board.Active.Height())
			})

			Convey("the shelf should load another piece from the level", func() {
				select {
				case s := <-game.ShelfUpdated:
					So(s, ShouldResemble, makeShelf("B", "C", "D", "X"))
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}

				So(game.Shelf(), ShouldResemble, makeShelf("B", "C", "D", "X"))
			})

			Convey("twice", func() {
				game.Board.Anchored <- Pieces.I // second time
				select {
				case _ = <-game.ShelfUpdated:
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}

				Convey("a new piece should be staged from the shelf", func() {
					select {
					case _ = <-game.ShelfUpdated:
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}
					So(game.Board.Active.Name, ShouldEqual, "B")
					So(game.Board.Active.Position.y, ShouldEqual, game.Board.height-game.Board.Active.Height())
				})

				Convey("the shelf should load another piece from the level", func() {
					select {
					case s := <-game.ShelfUpdated:
						So(s, ShouldResemble, makeShelf("C", "D", "X", "X"))
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}

					So(game.Shelf(), ShouldResemble, makeShelf("C", "D", "X", "X"))
				})
			})
		})

		Convey("when a piece collides", func() {
			game.Start()

			game.shelf.pieces = makeShelf("A", "B", "C", "D")

			game.Level = Level{0, 0, func() TetrisPiece {
				return TetrisPiece{Name: "X", Orientations: Pieces.O.Orientations}
			}, func(int) int { return 0 }}

			game.Board.Collision <- Pieces.I

			Convey("a new piece should be staged from the shelf", func() {
				So(game.Board.Active.Name, ShouldEqual, "A")
				So(game.Board.Active.Position.y, ShouldEqual, game.Board.height-game.Board.Active.Height())
			})

			Convey("the shelf should load another piece from the level", func() {
				select {
				case s := <-game.ShelfUpdated:
					So(s, ShouldResemble, makeShelf("B", "C", "D", "X"))
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}

				So(game.Shelf(), ShouldResemble, makeShelf("B", "C", "D", "X"))
			})

			Convey("twice", func() {
				game.Board.Collision <- Pieces.I // second time
				select {
				case _ = <-game.ShelfUpdated:
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}

				Convey("a new piece should be staged from the shelf", func() {
					select {
					case _ = <-game.ShelfUpdated:
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}
					So(game.Board.Active.Name, ShouldEqual, "B")
					So(game.Board.Active.Position.y, ShouldEqual, game.Board.height-game.Board.Active.Height())
				})

				Convey("the shelf should load another piece from the level", func() {
					select {
					case s := <-game.ShelfUpdated:
						So(s, ShouldResemble, makeShelf("C", "D", "X", "X"))
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}

					So(game.Shelf(), ShouldResemble, makeShelf("C", "D", "X", "X"))
				})
			})
		})

		Convey("the shelf", func() {

			Convey("should load 4 pieces from the Level when the game starts", func() {
				piece := 0
				game.Level = Level{0, 0, func() TetrisPiece {
					defer func() { piece++ }()
					return TetrisPiece{Name: fmt.Sprintf("Piece %d", piece), Orientations: Pieces.O.Orientations}
				}, func(int) int { return 0 }}

				game.Start()
				shelf := game.Shelf()

				for i := range shelf {
					So(shelf[i].Name, ShouldEqual, fmt.Sprintf("Piece %d", i+1)) // Piece 0 was already staged
				}
			})
		})

		Convey("when lines are cleared", func() {
			game.Start()

			Convey("the score increases according to the level", func() {
				game.Level = Level{0, 0, func() TetrisPiece { return Pieces.O }, func(lines int) int {
					return lines // just so we get a different number for different input
				}}
				Convey("for one line", func() {
					game.Board.Cleared <- []int{0}

					select {
					case s := <-game.ScoreChange:
						So(s, ShouldEqual, 1)
						So(game.score, ShouldEqual, 1)
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}
				})
			})
		})
	})
}

func makeShelf(p1, p2, p3, p4 string) [4]TetrisPiece {
	return [4]TetrisPiece{
		TetrisPiece{Name: p1, Orientations: Pieces.O.Orientations},
		TetrisPiece{Name: p2, Orientations: Pieces.O.Orientations},
		TetrisPiece{Name: p3, Orientations: Pieces.O.Orientations},
		TetrisPiece{Name: p4, Orientations: Pieces.O.Orientations}}
}
