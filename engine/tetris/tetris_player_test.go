package tetris

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTetrisPlayer(t *testing.T) {

	Convey("Given a tetris board", t, func() {
		board := NewTetrisBoard()

		Convey("when the player moves right", func() {
			board.Active.TetrisPiece = Pieces.O

			Convey("the piece should move to the right", func() {
				board.Active.Position = Point{3, 10}
				board.MoveRight()
				So(board.Active.Position.X, ShouldEqual, 4)
			})

			Convey("and the piece is against the wall, the piece should not move", func() {
				board.Active.Position = Point{9, 10} // box width is 2
				board.MoveRight()
				So(board.Active.Position.X, ShouldEqual, 9)
			})

			Convey("with wider piece against the wall, the piece should not move", func() {
				board.Active.TetrisPiece = TetrisPiece{Name: "TestPiece", Orientations: [][4]Point{}}
				board.Active.Orientations = append(board.Active.Orientations,
					[4]Point{Point{0, 0}, Point{1, 0}, Point{2, 0}, Point{3, 0}})
				board.Active.Position = Point{7, 10}
				board.MoveRight()
				So(board.Active.Position.X, ShouldEqual, 7)
			})

			Convey("and the piece is adjacent to a filled space, the piece should not move", func() {
				board.Active.Position = Point{0, 10}
				board.plane[11] = row("  **      ")
				board.MoveRight()
				So(board.Active.Position.X, ShouldEqual, 0)
			})
		})

		Convey("when the player moves left", func() {
			board.Active.TetrisPiece = Pieces.O

			Convey("the piece should move to the left", func() {
				board.Active.Position = Point{8, 10}
				board.MoveLeft()
				So(board.Active.Position.X, ShouldEqual, 7)
			})

			Convey("and the piece is against the wall, the piece should not move", func() {
				board.Active.Position = Point{0, 10}
				board.MoveLeft()
				So(board.Active.Position.X, ShouldEqual, 0)
			})

			Convey("and the piece is adjacent to a filled space, the piece should not move", func() {
				board.Active.Position = Point{4, 10}
				board.plane[11] = row("  **      ")
				board.MoveLeft()
				So(board.Active.Position.X, ShouldEqual, 4)
			})
		})

		Convey("when the player rotates", func() {
			Convey("an I piece", func() {
				board.Active.TetrisPiece = Pieces.I
				Convey("once", func() {
					board.RotateRight()

					Convey("the piece should be horizontal", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{1, 0}, Point{2, 0}, Point{3, 0}})
					})
				})

				Convey("twice", func() {
					board.RotateRight()
					board.RotateRight()

					Convey("the piece should be vertical", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{0, 3}})
					})
				})

				Convey("once, too close to the right wall", func() {
					board.Active.Position = Point{9, 10}
					board.RotateRight()

					Convey("the piece should rotate, and move left", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{1, 0}, Point{2, 0}, Point{3, 0}})
						So(board.Active.Position.X, ShouldEqual, 6)
					})
				})

				Convey("once, when there's not enough room", func() {
					board.Active.Position = Point{1, 10}
					board.plane[11] = row("   **     ")
					board.plane[10] = row("*  ***    ")
					board.plane[9] = row("*  **     ")
					board.plane[8] = row("*****     ")
					// and so on
					board.RotateRight()

					Convey("the piece should not rotate", func() {
						So(board.Active.Orientation, ShouldEqual, 0)
					})

					Convey("the piece should not move", func() {
						So(board.Active.Position.X, ShouldEqual, 1)
					})
				})
			})

			Convey("an O piece", func() {
				board.Active.TetrisPiece = Pieces.O

				Convey("once", func() {
					board.RotateRight()

					Convey("the piece should be the same", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{0, 1}, Point{1, 0}, Point{1, 1}})
					})
				})

				Convey("twice", func() {
					board.RotateRight()
					board.RotateRight()

					Convey("the piece should be the same", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{0, 1}, Point{1, 0}, Point{1, 1}})
					})
				})
			})

			Convey("a T piece", func() {
				board.Active.TetrisPiece = Pieces.T
				Convey("to the right", func() {
					Convey("once", func() {
						board.RotateRight()

						Convey("the piece should point right", func() {
							So(board.Active.Points(), ShouldResemble,
								[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{1, 1}})
						})
					})

					Convey("twice", func() {
						board.RotateRight()
						board.RotateRight()

						Convey("the piece should point down", func() {
							So(board.Active.Points(), ShouldResemble,
								[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{2, 1}})
						})
					})

					Convey("three times", func() {
						for i := 0; i < 3; i++ {
							board.RotateRight()
						}

						Convey("the piece should point left", func() {
							So(board.Active.Points(), ShouldResemble,
								[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{1, 2}})
						})
					})

					Convey("four times", func() {
						for i := 0; i < 4; i++ {
							board.RotateRight()
						}

						Convey("the piece should point up", func() {
							So(board.Active.Points(), ShouldResemble,
								[4]Point{Point{0, 0}, Point{1, 0}, Point{1, 1}, Point{2, 0}})
						})
					})
				})

				Convey("to the left", func() {
					Convey("once", func() {
						board.RotateLeft()

						Convey("the piece should point left", func() {
							So(board.Active.Points(), ShouldResemble,
								[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{1, 2}})
						})
					})
				})
			})
		})

		Convey("when the player drops the piece", func() {
			board.Active.TetrisPiece = Pieces.O
			board.Active.Position = Point{6, 17}

			Convey("with no filled spaces below, the piece should anchor to the floor", func() {
				board.Drop()
				select {
				case anchoredPiece := <-board.Anchored:
					So(anchoredPiece, ShouldResemble, Pieces.O)
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}
				So(board.plane[1], ShouldResemble, row("      OO  "))
				So(board.plane[0], ShouldResemble, row("      OO  "))
			})

			Convey("with a filled space below, the piece should anchor directly above it", func() {
				board.plane[1] = row("    ***   ")
				board.plane[0] = row("    *     ")
				board.Drop()

				select {
				case anchoredPiece := <-board.Anchored:
					So(anchoredPiece, ShouldResemble, Pieces.O)
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}

				So(board.plane[3], ShouldResemble, row("      OO  "))
				So(board.plane[2], ShouldResemble, row("      OO  "))
				So(board.plane[1], ShouldResemble, row("    ***   "))
				So(board.plane[0], ShouldResemble, row("    *     "))
			})
		})
	})
}
