import React from 'react';
import ReactDOM from 'react-dom';

import Board from './Board';

export default class App extends React.Component {
    render(){
        return (
            <div>
                <div>Hello Codemelee</div>
                <Board />
            </div>
        );
    }
}



ReactDOM.render(<Board />, document.getElementById('board'));