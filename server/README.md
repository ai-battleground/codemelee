# Notes

## Transport
use WebSockets from gorilla on the server, and send small JSON documents representing information updates in both directions.

## Encoding

We need to transmit player moves and control commands to the server, and transmit the board, current piece, shelf, speed, level, and score to the client

### The Board
The tetris board needs 3 bits (7 colors + 1 empty) per space, so a 20x10 board should be 600 bits, or 75 bytes minimum. 3 bit words are inconvenient, so wasting a 4th bit gives us 800 bits or 100 bytes.  This method gives us 2 spaces per byte, and 5 bytes per line. We will order the bytes left to right, starting on the bottom row and going up the board.

### Other game state
Each piece (tetromino) is represented by a letter (O, I, T, S, Z, L, J). The current piece additionally has position (x, y) and orientation (0-3).  The shelf contains 4 piece identifiers. The level is just a number, as is the speed and score. Additionally we will notify which lines are cleared in case the client wants to highlight them in some way. This will be an array of integers representing the rows, numbered zero based starting from the bottom.

### Player
Moves are left, right, rotate, and drop. Control commands are start, pause, quit. These will be represented by letters: S, P, Q for Start, Pause, Quit, and L, R, T, D for Left, Right, Rotate, Drop.
