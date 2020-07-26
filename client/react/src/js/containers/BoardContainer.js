import { connect } from 'react-redux'
import Board from '../components/Board'

const getPieceState = (pieceState, serverEvent) => {
    if (serverEvent.Piece) {
        return serverEvent;
    } else {
        return pieceState;
    }
}

const mapStateToProps = (state) => {
    return {
        piece: getPieceState(state.piece, state.serverEvent)
    }
}

const BoardContainer = connect(mapStateToProps)(Board)

export default BoardContainer