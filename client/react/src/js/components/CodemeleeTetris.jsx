import React from 'react';
import ReactDOM from 'react-dom';

import Board from './Board';

export default class App extends React.Component {
    constructor() {
        super();
        this.state = {
            name: ""
        }
    }

    componentDidMount() {
        this.connection = new GameConnection("ws://localhost:54545/tetris");
        this.connection.connect();
    }

    render(){
        return (
            <div>
                <div>Hello {this.state.name}</div>
            </div>
        );
    }
}


ReactDOM.render(<App />, document.getElementById('app'));
ReactDOM.render(<Board />, document.getElementById('board'));