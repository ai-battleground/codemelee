package tetris

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestStub(t *testing.T) {
    assert.True(t, true, "Canary test passing")
}

func TestInitialDifficulty(t *testing.T) {
    game := NewTetrisGame()
    assert.Equal(t, 1, game.difficulty, "Difficulty should default to 1")
}

func TestInitialBoardWidth(t *testing.T) {
    game := NewTetrisGame()
    assert.Equal(t, 10, game.Board.width, "Width should be 10")
}