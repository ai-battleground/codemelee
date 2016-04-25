package tetris

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestTetrisBoard(t *testing.T) {

	Convey("Given a tetris board", t, func() {
		board := NewTetrisBoard()

		Convey("when time is advanced", func() {
			Convey("and the piece is clear", func() {
				board.Active.TetrisPiece = Pieces.O
				board.Active.Position.y = 15
				board.Advance()

				Convey("the piece should descend", func() {
					So(board.Active.Position.y, ShouldEqual, 14)
				})
			})

			Convey("and the piece is at the bottom", func() {
				board.Active.TetrisPiece = Pieces.O
				board.Active.Position = Point{3, 0}
				board.Advance()

				Convey("the piece should be anchored to the board", func() {
					So(board.plane[1], ShouldResemble, row("   **     "))
					So(board.plane[0], ShouldResemble, row("   **     "))
				})

				Convey("the piece should be sent to the anchor channel", func() {
					select {
					case anchoredPiece := <-board.Anchored:
						So(anchoredPiece, ShouldResemble, Pieces.O)
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}
				})
			})

			Convey("and the piece is directly above at least one filled space", func() {
				board.Active.TetrisPiece = Pieces.O
				board.Active.Position = Point{7, 2}
				//                           xx
				//                           xx
				board.plane[1] = row("      **  ")
				board.plane[0] = row("     **   ")
				board.Advance()

				Convey("the piece should be anchored to the board", func() {
					So(board.plane[3], ShouldResemble, row("       ** "))
					So(board.plane[2], ShouldResemble, row("       ** "))
				})

				Convey("the piece should be sent to the anchor channel", func() {
					select {
					case anchoredPiece := <-board.Anchored:
						So(anchoredPiece, ShouldResemble, Pieces.O)
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}
				})
			})
		})

		Convey("when a piece is staged", func() {
			board.Stage(Pieces.O)

			Convey("the piece should be positioned at the top", func() {
				So(board.Active.Position.y, ShouldEqual, board.height-board.Active.Height())
			})

			Convey("the piece should be centered", func() {
				So(board.Active.Position.x, ShouldEqual, 4)
			})
		})

		Convey("when the player moves right", func() {
			board.Active.TetrisPiece = Pieces.O

			Convey("the piece should move to the right", func() {
				board.Active.Position = Point{3, 10}
				board.MoveRight()
				So(board.Active.Position.x, ShouldEqual, 4)
			})

			Convey("and the piece is against the wall, the piece should not move", func() {
				board.Active.Position = Point{9, 10} // box width is 2
				board.MoveRight()
				So(board.Active.Position.x, ShouldEqual, 9)
			})

			Convey("with wider piece against the wall, the piece should not move", func() {
				board.Active.TetrisPiece = TetrisPiece{Name: "TestPiece", Orientations: [][4]Point{}}
				board.Active.Orientations = append(board.Active.Orientations,
					[4]Point{Point{0, 0}, Point{1, 0}, Point{2, 0}, Point{3, 0}})
				board.Active.Position = Point{7, 10}
				board.MoveRight()
				So(board.Active.Position.x, ShouldEqual, 7)
			})

			Convey("and the piece is adjacent to a filled space, the piece should not move", func() {
				board.Active.Position = Point{0, 10}
				board.plane[11] = row("  **      ")
				board.MoveRight()
				So(board.Active.Position.x, ShouldEqual, 0)
			})
		})

		Convey("when the player moves left", func() {
			board.Active.TetrisPiece = Pieces.O

			Convey("the piece should move to the left", func() {
				board.Active.Position = Point{8, 10}
				board.MoveLeft()
				So(board.Active.Position.x, ShouldEqual, 7)
			})

			Convey("and the piece is against the wall, the piece should not move", func() {
				board.Active.Position = Point{0, 10}
				board.MoveLeft()
				So(board.Active.Position.x, ShouldEqual, 0)
			})

			Convey("and the piece is adjacent to a filled space, the piece should not move", func() {
				board.Active.Position = Point{4, 10}
				board.plane[11] = row("  **      ")
				board.MoveLeft()
				So(board.Active.Position.x, ShouldEqual, 4)
			})
		})

		Convey("when the player rotates", func() {
			Convey("an I piece", func() {
				board.Active.TetrisPiece = Pieces.I
				Convey("once", func() {
					board.Rotate()

					Convey("the piece should be horizontal", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{1, 0}, Point{2, 0}, Point{3, 0}})
					})
				})

				Convey("twice", func() {
					board.Rotate()
					board.Rotate()

					Convey("the piece should be vertical", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{0, 3}})
					})
				})
			})

			Convey("an O piece", func() {
				board.Active.TetrisPiece = Pieces.O

				Convey("once", func() {
					board.Rotate()

					Convey("the piece should be the same", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{0, 1}, Point{1, 0}, Point{1, 1}})
					})
				})

				Convey("twice", func() {
					board.Rotate()
					board.Rotate()

					Convey("the piece should be the same", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{0, 1}, Point{1, 0}, Point{1, 1}})
					})
				})
			})

			Convey("a T piece", func() {
				board.Active.TetrisPiece = Pieces.T

				Convey("once", func() {
					board.Rotate()

					Convey("the piece should point right", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{1, 1}})
					})
				})

				Convey("twice", func() {
					board.Rotate()
					board.Rotate()

					Convey("the piece should point down", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{2, 1}})
					})
				})

				Convey("three times", func() {
					for i := 0; i < 3; i++ {
						board.Rotate()
					}

					Convey("the piece should point left", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{1, 2}})
					})
				})

				Convey("four times", func() {
					for i := 0; i < 4; i++ {
						board.Rotate()
					}

					Convey("the piece should point up", func() {
						So(board.Active.Points(), ShouldResemble,
							[4]Point{Point{0, 0}, Point{1, 0}, Point{1, 1}, Point{2, 0}})
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
				So(board.plane[1], ShouldResemble, row("      **  "))
				So(board.plane[0], ShouldResemble, row("      **  "))
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

				So(board.plane[3], ShouldResemble, row("      **  "))
				So(board.plane[2], ShouldResemble, row("      **  "))
				So(board.plane[1], ShouldResemble, row("    ***   "))
				So(board.plane[0], ShouldResemble, row("    *     "))
			})
		})

		Convey("when the player scores", func() {
			board.plane[3] = row("********  ")
			board.plane[2] = row("********* ")
			board.plane[1] = row("* ********")
			board.plane[0] = row("* ********")

			board.Active.TetrisPiece = Pieces.I
			board.Active.Position = Point{9, 15}

			Convey("a single line", func() {
				board.Drop()

				Convey("the line should be sent to the cleared channel", func() {
					select {
					case clearedLines := <-board.Cleared:
						So(len(clearedLines), ShouldEqual, 1)
						So(clearedLines[0], ShouldEqual, 2)
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}
				})

				Convey("the line should be removed and higher lines dropped", func() {
					select {
					case _ = <-board.Cleared:
					case <-time.After(time.Second * 1):
					}
					So(board.plane[5], ShouldResemble, row("          "))
					So(board.plane[4], ShouldResemble, row("         *"))
					So(board.plane[3], ShouldResemble, row("         *"))
					So(board.plane[2], ShouldResemble, row("******** *"))
					So(board.plane[1], ShouldResemble, row("* ********"))
					So(board.plane[0], ShouldResemble, row("* ********"))
				})
			})

			Convey("three lines", func() {
				board.plane[4] = row("********* ")
				board.plane[3] = row("********* ")
				// already there:     *********
				// already there:     * ********
				// already there:     * ********

				board.Drop()

				Convey("the lines should be sent to the cleared channel", func() {
					select {
					case clearedLines := <-board.Cleared:
						So(len(clearedLines), ShouldEqual, 3)
						So(clearedLines[0], ShouldEqual, 2)
						So(clearedLines[1], ShouldEqual, 3)
						So(clearedLines[2], ShouldEqual, 4)
					case <-time.After(time.Second * 1):
						So(nil, ShouldNotBeNil)
					}
				})

				Convey("the lines should be removed and higher lines dropped", func() {
					select {
					case _ = <-board.Cleared:
					case <-time.After(time.Second * 1):
					}
					So(board.plane[5], ShouldResemble, row("          "))
					So(board.plane[4], ShouldResemble, row("          "))
					So(board.plane[3], ShouldResemble, row("          "))
					So(board.plane[2], ShouldResemble, row("         *"))
					So(board.plane[1], ShouldResemble, row("* ********"))
					So(board.plane[0], ShouldResemble, row("* ********"))
				})
			})
		})

		Convey("when a piece is staged that would collide", func() {
			for i := 0; i < 20; i++ {
				board.plane[i] = row("    **    ")
			}

			board.Stage(Pieces.O)

			Convey("the piece should be sent to the collision channel", func() {
				select {
				case collision := <-board.Collision:
					So(collision, ShouldResemble, Pieces.O)
				case <-time.After(time.Second * 1):
					So(nil, ShouldNotBeNil)
				}
			})

			Convey("the board is cleared", func() {
				for i := 0; i < 20; i++ {
					So(board.plane[i], ShouldResemble, row("          "))
				}
			})
		})
	})
}

func row(s string) []Space {
	row := make([]Space, 10)
	for i := 0; i < 10; i++ {
		row[i] = Space{empty: s[i] == byte(' ')}
	}
	return row
}
