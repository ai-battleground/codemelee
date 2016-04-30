import React from 'react';


export default class Tile {
  constructor(props) {
    this.position = props.position;
    this.colors = props.colors;
    this.size = props.size;
    this.parentProjection = props.projection;
  }

  projection(context) {
    this.parentProjection(context);
    context.scale(this.size, this.size);
    context.translate(this.position.x, this.position.y);
  }

  render(state) {
    const context = state.context;
    
    context.save();
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