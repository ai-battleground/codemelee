import React from 'react';


export default class Tile extends React.Component {
  constructor(args) {
    super();
    this.position = args.position;
    this.colors = args.colors;
    this.size = args.size;
    this.boardProjection = args.boardProjection;
  }

  projection(context) {
    context.scale(this.size, this.size);
    context.translate(this.position.x, this.position.y);
  }

  render(state) {
    const context = state.context;
    context.save();
    this.boardProjection(context);
    this.projection(context);

    context.strokeStyle = this.colors.fg;
    context.fillStyle = this.colors.bg;
    context.lineWidth = 2 / this.size;
    context.strokeRect(0, 0, 1, 1);
    context.fillRect(0, 0, 1, 1);

    context.fill();
    context.stroke();
    context.restore();
  }
}