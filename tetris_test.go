package tetris

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type TetrisTestSuite struct {
    suite.Suite
    NewTetrisGame *TetrisGame
}

func (suite *TetrisTestSuite) SetupTest() {
    suite.NewTetrisGame = NewTetrisGame()
}

func TestStub(t *testing.T) {
    assert.True(t, true, "Canary test passing")
}

func (suite *TetrisTestSuite) TestInitialDifficulty() {
    assert.Equal(suite.T(), 1, suite.NewTetrisGame.difficulty, "Difficulty should default to 1")
}

func (suite *TetrisTestSuite) TestInitialBoardDimenstions() {
    assert.Equal(suite.T(), 10, suite.NewTetrisGame.Board.width, "Width should be 10")
    assert.Equal(suite.T(), 20, suite.NewTetrisGame.Board.height, "Height should be 10")
}

func (suite *TetrisTestSuite) TestInitialSpeed() {
    assert.Equal(suite.T(), 1, suite.NewTetrisGame.speed, "Speed should default to 1")
}

func TestTetrisTestSuite(t *testing.T) {
    suite.Run(t, new(TetrisTestSuite))
}