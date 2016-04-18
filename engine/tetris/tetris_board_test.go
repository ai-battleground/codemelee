package tetris

import (
    "testing"
    "time"
    . "github.com/smartystreets/goconvey/convey"
)

func TestTetrisBoard(t *testing.T) {

    Convey("Given a tetris board", t, func() {
        board := NewTetrisBoard()

        Convey("when time is advanced", func() {
            Convey("and the piece is clear", func() {
                board.Piece = Pieces.Box
                board.PiecePosition.y = 15
                board.Advance()

                Convey("the piece should descend", func() {
                    So(board.PiecePosition.y, ShouldEqual, 14)
                })
            })

            Convey("and the piece is at the bottom", func() {
                board.Piece = Pieces.Box
                board.PiecePosition = Point{3, 0}
                board.Advance()

                Convey("the piece should be anchored to the board", func() {
                    So(board.plane[0], ShouldResemble, row("   **     "))
                    So(board.plane[1], ShouldResemble, row("   **     "))
                })

                Convey("the piece should be sent to the anchor channel", func() {
                    select {
                        case anchoredPiece := <-board.Anchored:
                            So(anchoredPiece, ShouldEqual, Pieces.Box)
                        case <-time.After(time.Second * 1):
                            So(nil, ShouldNotBeNil)
                    }
                })
            })

            Convey("and the piece is directly above at least one filled space", func() {
                board.Piece = Pieces.Box
                board.PiecePosition = Point{7, 2}
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
                            So(anchoredPiece, ShouldEqual, Pieces.Box)
                        case <-time.After(time.Second * 1):
                            So(nil, ShouldNotBeNil)
                    }
                })
            })
        })

        Convey("when a piece is staged", func() {
            board.Stage(Pieces.Box)
            
            Convey("the piece should be positioned at the top", func() {
                So(board.PiecePosition.y, ShouldEqual, board.height - board.Piece.height)
            })

            Convey("the piece should be centered", func() {
                So(board.PiecePosition.x, ShouldEqual, 4)
            })
        })

        Convey("when the player moves right", func() {
            board.Piece = Pieces.Box

            Convey("the piece should move to the right", func() {
                board.PiecePosition = Point{3, 10}
                board.MoveRight()
                So(board.PiecePosition.x, ShouldEqual, 4)
            })

            Convey("and the piece is against the wall, the piece should not move", func() {
                board.PiecePosition = Point{9, 10} // box width is 2
                board.MoveRight()
                So(board.PiecePosition.x, ShouldEqual, 9)
            })

            Convey("with wider piece against the wall, the piece should not move", func() {
                board.Piece = &TetrisPiece{width:4, name: "TestPiece", Points:[]Point{Point{0,0},Point{3,0}}}
                board.PiecePosition = Point{7, 10}
                board.MoveRight()
                So(board.PiecePosition.x, ShouldEqual, 7)
            })

            Convey("and the piece is adjacent to a filled space, the piece should not move", func() {
                board.PiecePosition = Point{0, 10}
                board.plane[11] = row("  **      ")
                board.MoveRight()
                So(board.PiecePosition.x, ShouldEqual, 0)
            })
        })

        Convey("when the player moves left", func() {
            board.Piece = Pieces.Box

            Convey("the piece should move to the left", func() {
                board.PiecePosition = Point{8, 10}
                board.MoveLeft()
                So(board.PiecePosition.x, ShouldEqual, 7)
            })

            Convey("and the piece is against the wall, the piece should not move", func() {
                board.PiecePosition = Point{0, 10}
                board.MoveLeft()
                So(board.PiecePosition.x, ShouldEqual, 0)
            })

            Convey("and the piece is adjacent to a filled space, the piece should not move", func() {
                board.PiecePosition = Point{4, 10}
                board.plane[11] = row("  **      ")
                board.MoveLeft()
                So(board.PiecePosition.x, ShouldEqual, 4)
            })
        })
    })
}

func row(s string) []Space {
    row := make([]Space, 10)
    for i := 0; i<10; i++ {
        row[i] = Space{empty: s[i] == byte(' ')}
    }
    return row
}