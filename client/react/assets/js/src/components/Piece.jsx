import React from 'react';
import Tile from './Tile';
import { Point, point } from '../graphics'

const p = point;
const cellSize = 18;

export default class Piece {
  constructor(props) {
    this.position = props.position;
    this.tetromino = props.tetromino;
    this.parentProjection = props.projection;
    this.tiles = this.tetromino.points.map(p => new Tile({
        position: p,
        colors: this.tetromino.colors,
        size: cellSize,
        projection: this.projection.bind(this)
    }));
  }

  static get O() {
    return {
        points: [p(0, 0), p(0, 1), p(1, 0), p(1, 1)],
        colors: {
            fg: "#662200",
            bg: "#FF6600"
        }
    }
  }

  static get I() {
    return {
        points: [p(0, 0), p(0, 1), p(0, 2), p(0, 3)],
        colors: {
            fg: "#227700",
            bg: "#44FF00"
        }
    }
  }

  // Graphics

  projection(context) {
    this.parentProjection(context);
    context.translate(this.position.x * cellSize, this.position.y * cellSize);
  }

  render(context) {
    for (let tile of this.tiles) {
        tile.render(context);
    }
  }
}