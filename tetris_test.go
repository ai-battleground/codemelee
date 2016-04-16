package tetris

import (
    "testing"
    "fmt"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type TetrisTestSuite struct {
    suite.Suite
    Game *TetrisGame
}

func (suite *TetrisTestSuite) SetupTest() {
    suite.Game = NewTetrisGame()
}

func TestStub(t *testing.T) {
    assert.True(t, true, "Canary test passing")
}

func (suite *TetrisTestSuite) TestInitialDifficulty() {
    assert.Equal(suite.T(), 1, suite.Game.difficulty, "Difficulty should default to 1")
}

func (suite *TetrisTestSuite) TestInitialSpeed() {
    assert.Equal(suite.T(), 1, suite.Game.speed, "Speed should default to 1")
}

func (suite *TetrisTestSuite) TestInitialScore() {
    assert.Equal(suite.T(), 0, suite.Game.score, "Score should default to 0")
}

func (suite *TetrisTestSuite) TestInitialBoardDimensions() {
    assert.Equal(suite.T(), 10, suite.Game.Board.width, "Width should be 10")
    assert.Equal(suite.T(), 20, suite.Game.Board.height, "Height should be 10")
}

func (suite *TetrisTestSuite) TestInitialBoardEmpty() {
    for y, row := range suite.Game.Board.plane {
        for x, space := range row {
            assert.True(suite.T(), space.empty, fmt.Sprintf("Space %d, %d should be empty", x, y))
        }
    }
}

func TestTetrisTestSuite(t *testing.T) {
    suite.Run(t, new(TetrisTestSuite))
}