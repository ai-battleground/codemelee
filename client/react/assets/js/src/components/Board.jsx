import React from 'react';
import Piece from './Piece';
import { Point } from '../graphics'

export default class Board extends React.Component {
    constructor() {
        super();
        this.state = {
          screen: {
            width: 180,
            height: 360,
            ratio: window.devicePixelRatio || 1,
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

        context.fillStyle = '#000';
        context.fillRect(0, 0, this.state.screen.width, this.state.screen.height);

        context.restore();

        // Render Piece
        this.state.piece.render(context)

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

