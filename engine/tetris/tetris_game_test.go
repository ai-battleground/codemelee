package tetris

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestTetrisGame(t *testing.T) {
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
}
