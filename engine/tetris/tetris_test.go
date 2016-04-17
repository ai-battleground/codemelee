package tetris

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestTetris(t *testing.T) {
    Convey("Given a new tetris game", t, func() {
        game := NewTetrisGame()

        Convey("the difficulty should be 1", func() {
            So(game.difficulty, ShouldEqual, 1)
        })

        Convey("the speed should be 1", func() {
            So(game.speed, ShouldEqual, 1)
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

            Convey("the player", func() {
                player := game.Player

                Convey("should have a box", func() {
                    So(player.Piece, ShouldEqual, Pieces.Box)
                })
            })
        })
    })
}