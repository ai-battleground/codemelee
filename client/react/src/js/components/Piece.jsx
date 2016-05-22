import React from 'react';
import Tile from './Tile';
import { Point, point } from '../graphics'
import { Group, Rect } from 'react-konva'

export default class Piece extends React.Component {

  static get O() {
    return {
        points: [point(0, 0), point(0, 1), point(1, 0), point(1, 1)],
        colors: {
            fg: "#662200",
            bg: "#FF6600"
        }
    }
  }

  static get I() {
    return {
        points: [point(0, 0), point(0, 1), point(0, 2), point(0, 3)],
        colors: {
            fg: "#227700",
            bg: "#44FF00"
        }
    }
  }

  moveLeft() {
    if (! this.leftLocked) {
      this.position.x -= 1;
      this.leftLocked = true;
    }
  }

  moveRight() {
    if (! this.rightLocked) {
        this.position.x += 1;
        this.rightLocked = true;
    }
  }

  unlockLeft() {
    this.leftLocked = false;
  }

  unlockRight() {
    this.rightLocked = false;
  }

  render() {
    return (
      <Group 
        x={this.props.position.x * this.props.cellSize} 
        y={this.props.position.y * this.props.cellSize}>
        {this.props.tetromino.points.map( (p, i) => 
          <Tile 
            key={i}
            size={this.props.cellSize}
            position={p}
            colors={this.props.tetromino.colors}/>
        )}
      </Group>
    )
  }
}