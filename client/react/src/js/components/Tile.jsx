import React from 'react';
import { Rect } from 'react-konva'

export default class Tile extends React.Component {

  render() {
    return (
      <Rect 
        x={this.props.position.x * this.props.size}
        y={this.props.position.y * this.props.size}
        width={this.props.size} height={this.props.size}
        fill={this.props.colors.bg}
        stroke={this.props.colors.fg}
        strokeWidth={2} />
    )
  }
}