package tetris

import (
	"fmt"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTetrisGame(t *testing.T) {
	Convey("Given a tetris game", t, func() {
		game := NewGame()

		Convey("the initial level should be 1", func() {
			So(game.Level.number, ShouldEqual, 1)
		})

		Convey("the initial speed should be 1", func() {
			So(game.Level.speed, ShouldEqual, 1)
		})

		Convey("the initial score should be 0", func() {
			So(game.Score(), ShouldEqual, 0)
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
						So(space.Empty(), ShouldBeTrue)
					}
				}
			})
		})

		Convey("when the game is started", func() {

			Convey("the board", func() {
				board := game.Board

				Convey("should have an active piece", func() {
					game.Start()
					So(board.Active.Piece, ShouldResemble, Pieces.O)
				})

				Convey("should position the piece at the top", func() {
					game.Start()
					So(board.Active.Position.Y, ShouldEqual, board.height-board.Active.Height())
				})

				Convey("should advance the board according to the level's speed", func() {
					game.Level.speed = 12
					game.Start()
					time.Sleep(time.Second/3 + 50*time.Millisecond)

					So(board.Active.Position.Y, ShouldEqual, board.height-board.Active.Height()-4)
				})
			})

			Convey("the shelf", func() {
				game.Start()
				shelf := game.Shelf()

				Convey("should load 4 O pieces", func() {
					for i := 0; i < 4; i++ {
						So(shelf[i], ShouldResemble, Pieces.O)
					}
				})
			})
		})

		Convey("when a piece is anchored", func() {
			game.shelf.pieces = makeShelf("A", "B", "C", "D")

			game.Level.NextPiece = func() Piece {
				return Piece{Name: "X", Orientations: Pieces.O.Orientations}
			}

			game.Board.Anchored <- Pieces.I

			Convey("a new piece should be staged from the shelf", func() {
				select {
				case _ = <-game.ShelfUpdated:
					So(game.Board.Active.Name, ShouldEqual, "A")
					So(game.Board.Active.Position.Y, ShouldEqual, game.Board.height-game.Board.Active.Height())
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}
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
					So(game.Board.Active.Position.Y, ShouldEqual, game.Board.height-game.Board.Active.Height())
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
				game.Level.NextPiece = func() Piece {
					defer func() { piece++ }()
					return Piece{Name: fmt.Sprintf("Piece %d", piece), Orientations: Pieces.O.Orientations}
				}

				game.Start()
				shelf := game.Shelf()

				for i := range shelf {
					So(shelf[i].Name, ShouldEqual, fmt.Sprintf("Piece %d", i+1)) // Piece 0 was already staged
				}
			})
		})

		Convey("when a piece collides", func() {
			for i := 0; i < 20; i++ {
				game.Board.plane[i] = row("    **    ")
			}
			game.Start()
			So(game.state, ShouldEqual, GameOver)
		})

		Convey("when lines are cleared", func() {
			pieces := []Piece{Pieces.T, Pieces.I, Pieces.I}
			game.Level.NextPiece = func() Piece {
				p := pieces[0]
				pieces = append(pieces[1:], p)
				return p
			}
			game.Board.plane[3] = row("* ********")
			game.Board.plane[2] = row("********* ")
			game.Board.plane[1] = row("***** ****")
			game.Board.plane[0] = row("* ********")

			Convey("the score increases according to the level", func() {
				Convey("for one line", func() {
					game.Start()
					time.Sleep(200 * time.Millisecond)
					game.RotateRight()
					game.MoveLeft()
					game.MoveLeft()
					game.MoveLeft()
					game.Drop()

					So(game.score, ShouldEqual, 1)

					Convey("several times", func() {
						time.Sleep(200 * time.Millisecond)
						for i := 0; i < 5; i++ {
							game.MoveRight()
						}
						game.Drop()
						time.Sleep(200 * time.Millisecond)
						game.MoveRight()
						game.Drop()

						So(game.score, ShouldEqual, 3)
					})
				})

				Convey("for two lines", func() {
					game.Board.plane[4] = row("* ********")
					game.Start()
					game.MoveRight()
					game.MoveRight()
					game.Drop()
					time.Sleep(200 * time.Millisecond)
					game.MoveLeft()
					game.MoveLeft()
					game.MoveLeft()
					game.Drop()

					So(game.score, ShouldEqual, 4)
				})
			})
		})

		Convey("level 1", func() {

			Convey("only produces O pieces", func() {
				game.Start()
				for i := 0; i < 8; i++ {
					time.Sleep(100 * time.Millisecond)
					So(game.Active.Piece, should.Resemble, Pieces.O)
					game.Drop()
				}
			})

			Convey("speed is about 1 tick per second (SLOW)", func() {
				game.Start()
				time.Sleep(750 * time.Millisecond)
				So(game.Board.Active.Position.Y, should.Equal, 18) // O piece
				time.Sleep(500 * time.Millisecond)
				So(game.Board.Active.Position.Y, should.Equal, 17) // O piece
			})

			Convey("ends after 10 lines", func() {
				oMoves := []func(){
					func() { // move left 4
						for i := 0; i < 4; i++ {
							game.MoveLeft()
						}
					},
					func() { // move left 2
						for i := 0; i < 2; i++ {
							game.MoveLeft()
						}
					},
					func() {}, // no moves or rotations
					func() { // move right 2
						for i := 0; i < 2; i++ {
							game.MoveRight()
						}
					},
					func() { // move right 4
						for i := 0; i < 4; i++ {
							game.MoveRight()
						}
					},
				}
				game.Start()
				So(game.Level.number, should.Equal, 1)
				for j := 0; j < 5; j++ { // five O pieces make 2 lines, so do this 5 times to make 10 lines
					for _, move := range oMoves {
						move()
						game.Drop()
						time.Sleep(20 * time.Millisecond)
					}
				}
				time.Sleep(50 * time.Millisecond)
				So(game.Level.number, should.Equal, 2)
			})

		})

		Convey("level 2", func() {
			game.Level = getLevel(1)

			Convey("produces only O and I pieces", func() {
				game.Start()
				oCount := 0
				iCount := 0
				for i := 0; i < 20; i++ {
					time.Sleep(10 * time.Millisecond)
					So(game.Active.Piece.Name, should.BeIn, []string{"O", "I"})
					if game.Active.Piece.Name == "O" {
						oCount++
						for i := 0; i < 4; i++ {
							game.MoveLeft()
						}
					}
					if game.Active.Piece.Name == "I" {
						iCount++
						game.RotateRight()
						game.MoveRight()
						game.MoveRight()
					}
					game.Drop()
				}
				So(oCount, should.BeGreaterThan, 0)
				So(iCount, should.BeGreaterThan, 0)
			})
		})
		Convey("level 3", func() {
			game.Level = getLevel(2)

			Convey("produces only O, I, and T pieces", func() {
				pieceCount := map[string]int{}
				for i := 0; i < 30; i++ {
					p := game.Level.NextPiece()
					pieceCount[p.Name] += 1
				}
				So(len(pieceCount), ShouldEqual, 3)
				for name, count := range pieceCount {
					So(name, ShouldBeIn, []string{"O", "I", "T"})
					So(count, ShouldBeGreaterThan, 0)
				}
			})
		})
		//Convey("level 4", func() {
		//	game.Level = getLevel(3)
		//	Convey("produces only O, I, T, L, J pieces", func() {
		//		pieceCount := map[string]int{}
		//		for i :=0; i<30; i++ {
		//			p := game.Level.NextPiece()
		//			pieceCount[p.Name] += 1
		//		}
		//		So(len(pieceCount), ShouldEqual, 3)
		//		for name, count := range pieceCount {
		//			So(name, ShouldBeIn, []string{"O", "I", "T", "L", "J"})
		//			So(count, ShouldBeGreaterThan, 0)
		//		}
		//	})
		//})
	})
}

func makeShelf(p1, p2, p3, p4 string) [4]Piece {
	return [4]Piece{
		Piece{Name: p1, Orientations: Pieces.O.Orientations},
		Piece{Name: p2, Orientations: Pieces.O.Orientations},
		Piece{Name: p3, Orientations: Pieces.O.Orientations},
		Piece{Name: p4, Orientations: Pieces.O.Orientations}}
}
