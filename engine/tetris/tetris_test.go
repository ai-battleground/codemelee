package tetris

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestTetris(t *testing.T) {
    Convey("Given a new tetris game", t, func() {
        game := NewTetrisGame()

        Convey("the level should be 1", func() {
            So(game.Level.number, ShouldEqual, 1)
        })

        Convey("the speed should be 1", func() {
            So(game.Level.speed, ShouldEqual, 1)
        })

        Convey("the score should be 0", func() {
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
                    So(board.Piece, ShouldEqual, Pieces.Box)
                })

                Convey("should position the piece at the top", func() {
                    So(board.PiecePosition.y, ShouldEqual, board.height - board.Piece.height)
                })
            })
        })
    })

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

                Convey("a new piece should be staged", func() {
                    So(board.PiecePosition.y, ShouldEqual, board.height - board.Piece.height)
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

                Convey("a new piece should be staged", func() {
                    So(board.PiecePosition.y, ShouldEqual, board.height - board.Piece.height)
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