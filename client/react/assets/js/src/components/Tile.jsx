import React from 'react';


export default class Tile extends React.Component {
  constructor(args) {
    super();
    this.position = args.position;
    this.colors = args.colors;
    this.size = args.size;
  }

  render(state) {
    const context = state.context;
    let origin = {
        x: this.position.x * this.size,
        y: this.position.y * this.size
    };
    context.save();
    context.translate(origin.x, origin.y);

    context.strokeStyle = this.colors.fg;
    context.fillStyle = this.colors.bg;
    context.lineWidth = 2;

    context.beginPath();
    context.lineTo(origin.x, origin.y + this.size);
    context.lineTo(origin.x + this.size, origin.y + this.size);
    context.lineTo(origin.x + this.size, origin.y);
    context.lineTo(origin.x, origin.y);
    context.closePath();

    context.fill();
    context.stroke();
    context.restore();
  }
}