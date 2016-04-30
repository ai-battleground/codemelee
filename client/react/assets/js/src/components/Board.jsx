import React from 'react';
import Tile from './Tile';

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
          tiles: [],
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
        for (var i = 0; i<4; i++) {
            this.state.tiles.push(new Tile({
                position: {
                    x: Math.floor(i / 2),
                    y: i % 2
                },
                colors: {
                    fg: "#660000",
                    bg: "#FF0000"
                },
                size: 18,
                boardProjection: this.projection.bind(this)
            }));
        }
    }

    update() {
        const context = this.state.context;
        context.save();

        this.projection(context);

        context.fillStyle = '#000';
        context.globalAlpha = 0.4;
        context.fillRect(0, 0, this.state.screen.width, this.state.screen.height);
        context.globalAlpha = 1;

        context.restore();

        this.animate(new Date().getTime());

        // Render Tiles
        for (let tile of this.state.tiles) {
            tile.render(this.state);
        }
        // Next frame
        requestAnimationFrame(() => this.update());
    }

    tick() {
        if (this.state.running) {
            for (let tile of this.state.tiles) {
                tile.position.y += 1;
            }
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

