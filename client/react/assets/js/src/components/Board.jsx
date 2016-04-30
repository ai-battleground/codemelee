import React from 'react';
import Piece from './Piece';
import { Point } from '../graphics'

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
    constructor() {
        super();
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
          running: false,
          context: null,
          lastAnimated: 0
        };
    }

    componentDidMount() {
        window.addEventListener('keyup',   this.handleKeys.bind(this, false));
        window.addEventListener('keydown', this.handleKeys.bind(this, true));

        const context = this.refs.canvas.getContext('2d');
        this.setState({ context: context });
        this.start();
        requestAnimationFrame(() => this.update());
    }

    start() {
        this.state.running = true;
        this.state.piece = new Piece({
            position: new Point(1, 0),
            tetromino: Piece.I,
            projection: this.projection.bind(this)
        });
    }

    update() {
        const context = this.state.context;
        context.save();

        this.projection(context);
        this.drawBoard(context, this.state.screen.width, this.state.screen.height);

        context.restore();

        // Render Piece
        this.state.piece.render(this.state)

        this.animate(new Date().getTime());
        // Next frame
        requestAnimationFrame(() => this.update());
    }

    tick() {
        if (this.state.running) {
            this.state.piece.position.y += 1;
        }
    }

    projection(context) {
        context.translate(this.state.position.x, this.state.position.y);
    }

    animate(timestamp) {
        if (timestamp - this.state.lastAnimated > 1000) {
            this.tick();
            this.state.lastAnimated = timestamp;
        }
    }

    drawBoard(context, width, height) {
        // background
        context.fillStyle = "#000";
        context.fillRect(0, 0, width, height);
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
                <canvas ref="canvas"
                  width={this.state.screen.width * this.state.screen.ratio}
                  height={this.state.screen.height * this.state.screen.ratio}
                />
            </div>
        );
    }
}

