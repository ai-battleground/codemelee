import React from 'react';
import Piece from './Piece';
import { Point } from '../graphics'
import { Layer, Rect, Stage } from 'react-konva'

const KEY = {
  LEFT:  37,
  RIGHT: 39,
  DOWN:  40,
  UP:    38,
  A:     65,
  D:     68,
  S:     83,
  W:     87,
  SPACE: 32
};

export default class Board extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
          screen: {
            width: 180,
            height: 360,
            ratio: window.devicePixelRatio || 1,
          },
          keys : {
            left  : 0,
            right : 0,
            up    : 0,
            down  : 0,
            space : 0,
          },
          position: {
            x: 10,
            y: 20
          },
          piece: {
            position: new Point(1, 0)
          },
          running: false,
          context: null,
          lastAnimated: 0
        };
        this.animate = this.animate.bind(this);
        this.tick = this.tick.bind(this);
        this.update = this.update.bind(this);
        this.start = this.start.bind(this);
    }

    componentDidMount() {
        window.addEventListener('keyup',   this.handleKeys.bind(this, false));
        window.addEventListener('keydown', this.handleKeys.bind(this, true));

        this.start();
        requestAnimationFrame(() => this.update());
    }

    start() {
        this.state.running = true;
    }

    update() {
        this.animate(new Date().getTime());
        // Next frame
        requestAnimationFrame(() => this.update());
    }

    tick() {
        if (this.state.running) {
          this.setState({piece: {position: {x: 2, y: this.state.piece.position.y + 1}}})
        }
    }

    animate(timestamp) {
        if (timestamp - this.state.lastAnimated > 1000) {
            this.tick();
            this.state.lastAnimated = timestamp;
        }
    }

    handleKeys(value, e){
      let keys = this.state.keys;
      if(e.keyCode === KEY.LEFT   || e.keyCode === KEY.A) keys.left  = value;
      if(e.keyCode === KEY.RIGHT  || e.keyCode === KEY.D) keys.right = value;
      if(e.keyCode === KEY.UP     || e.keyCode === KEY.W) keys.up    = value;
      if(e.keyCode === KEY.DOWN   || e.keyCode === KEY.S) keys.down  = value;
      if(e.keyCode === KEY.SPACE) keys.space = value;
      this.setState({
        keys : keys
      });
    }

    render() {
      return (
          <div className="board">
            <Stage width={480} height={480}>
              <Layer x={10} y={20}>
                <Rect
                  x={0} y={0} width={180} height={360} fill={'black'} />
                <Piece position={this.state.piece.position}
                     tetromino={Piece.I} cellSize={18} />
              </Layer>
            </Stage>
          </div>
      )
    }
}

