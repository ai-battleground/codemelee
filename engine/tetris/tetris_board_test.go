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