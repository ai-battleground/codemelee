package tetris

import (
	"fmt"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestTetrisBoard(t *testing.T) {

	Convey("Given a tetris board", t, func() {
		board := NewBoard()

		Convey("when time is advanced", func() {
			Convey("and the piece is clear", func() {
				board.Active.Piece = Pieces.O
				board.Active.Position.Y = 15
				board.Advance()

				Convey("the piece should descend", func() {
					So(board.Active.Position.Y, ShouldEqual, 14)
				})
			})

			Convey("and the piece is at the bottom", func() {
				board.Active.Piece = Pieces.O
				board.Active.Position = Point{3, 0}
				board.Advance()

				Convey("the piece should be anchored to the board", func() {
					So(board.plane[1], ShouldResemble, row("   OO     "))
					So(board.plane[0], ShouldResemble, row("   OO     "))
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
				board.Active.Piece = Pieces.O
				board.Active.Position = Point{7, 2}
				//                           xx
				//                           xx
				board.plane[1] = row("      **  ")
				board.plane[0] = row("     **   ")
				board.Advance()

				Convey("the piece should be anchored to the board", func() {
					So(board.plane[3], ShouldResemble, row("       OO "))
					So(board.plane[2], ShouldResemble, row("       OO "))
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
				So(board.Active.Position.Y, ShouldEqual, board.height-board.Active.Height())
			})

			Convey("the piece should be centered", func() {
				So(board.Active.Position.X, ShouldEqual, 4)
			})
		})

		Convey("when the player scores", func() {
			board.plane[3] = row("********  ")
			board.plane[2] = row("********* ")
			board.plane[1] = row("* ********")
			board.plane[0] = row("* ********")

			board.Active.Piece = Pieces.I
			board.Active.Position = Point{9, 15}

			Convey("a single line", func() {
				var cleared []int
				board.OnClearedLines(func(lines []int) {
					cleared = lines
				})
				board.Drop()

				Convey("the line should be sent to observers", func() {
					So(len(cleared), should.Equal, 1)
					So(cleared[0], should.Equal, 2)
				})

				Convey("the line should be removed and higher lines dropped", func() {
					So(board.plane[5], ShouldResemble, row("          "))
					So(board.plane[4], ShouldResemble, row("         I"))
					So(board.plane[3], ShouldResemble, row("         I"))
					So(board.plane[2], ShouldResemble, row("******** I"))
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

				var cleared []int
				board.OnClearedLines(func(lines []int) {
					cleared = lines
				})
				board.Drop()

				Convey("the lines should be sent to observers", func() {
					So(len(cleared), should.Equal, 3)
					So(cleared, should.Resemble, []int{2, 3, 4})
				})

				Convey("the lines should be removed and higher lines dropped", func() {
					So(board.plane[5], ShouldResemble, row("          "))
					So(board.plane[4], ShouldResemble, row("          "))
					So(board.plane[3], ShouldResemble, row("          "))
					So(board.plane[2], ShouldResemble, row("         I"))
					So(board.plane[1], ShouldResemble, row("* ********"))
					So(board.plane[0], ShouldResemble, row("* ********"))
				})
			})
		})

		Convey("when a piece is staged that would collide", func() {
			collisionDetected := false
			board.OnCollision(func() {
				collisionDetected = true
			})
			for i := 0; i < 20; i++ {
				board.plane[i] = row("    **    ")
			}

			board.Stage(Pieces.O)

			Convey("a collision event should be published", func() {
				So(collisionDetected, ShouldBeTrue)
			})
		})

		Convey("when a snapshot is taken", func() {
			board.plane[3] = row("LLL       ")
			board.plane[2] = row("L         ")
			board.plane[1] = row("OO   T    ")
			board.plane[0] = row("OO  TTT   ")

			Convey("it should have the same height as the board", func() {
				snapshot := board.TakeSnapshot()
				snapshotLines := strings.Split(snapshot, "\n")
				So(len(snapshotLines), ShouldEqual, board.height)
				for lineNum, line := range snapshotLines {
					Convey(fmt.Sprintf("Line %d should have the same width as the board", lineNum), func() {
						So(len(line), ShouldEqual, board.width)
					})
				}
			})

			Convey("it should reflect the board state", func() {
				snapshot := board.TakeSnapshot()
				snapshotLines := strings.Split(snapshot, "\n")
				expectedLastFourLines := "LLL       \n" +
					"L         \n" +
					"OO   T    \n" +
					"OO  TTT   "
				snapshotLastFourLines := strings.Join(snapshotLines[16:], "\n")
				So(snapshotLastFourLines, ShouldEqual, expectedLastFourLines)
			})

			Convey("it should reflect the active piece", func() {
				board.Stage(Pieces.T)
				board.Advance()
				board.Advance()
				board.RotateRight()

				snapshot := board.TakeSnapshot()
				snapshotLines := strings.Split(snapshot, "\n")

				expectedPieceLines := "    T     \n" +
					"    TT    \n" +
					"    T     "
				snapshotRelevantLines := strings.Join(snapshotLines[1:4], "\n")
				So(snapshotRelevantLines, ShouldEqual, expectedPieceLines)
			})
		})
	})
}

func row(s string) []Space {
	row := make([]Space, 10)
	for i := 0; i < 10; i++ {
		row[i] = Space{contents: s[i]}
	}
	return row
}
