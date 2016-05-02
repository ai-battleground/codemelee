export class GameConnection {
    constructor(url) {
        this.socket = new WebSocket(url)
    }

    connect() {
        socket.onopen = () => {
            socket.send("Ping");
        };
        socket.onmessage = (message) => {
            console.log(`[TETRIS] ${message}`)
        }
        socket.onerror = (error) => console.log(error);
    }
}