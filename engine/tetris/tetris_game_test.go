package tetris

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
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
                    So(board.Piece, ShouldEqual, Pieces.Box)
                })

                Convey("should position the piece at the top", func() {
                    So(board.PiecePosition.y, ShouldEqual, board.height - board.Piece.height)
                })
            })

        })

        Convey("when a piece is anchored", func() {
            game.Start()

            game.Board.Anchored <- Pieces.Box

            Convey("a new piece should be staged", func() {
                So(game.Board.Piece, ShouldEqual, Pieces.Box)
                So(game.Board.PiecePosition.y, ShouldEqual, game.Board.height - game.Board.Piece.height)
            })
        })
    })
}

