export class GameConnection {
    constructor(url) {
        this.socket = new WebSocket(url);
    }

    connect() {
        this.socket.onopen = () => {
            this.socket.send("Ping");
        };
        this.socket.onmessage = (message) => {
            console.log(`[TETRIS] ${message}`);
        }
        this.socket.onerror = (error) => console.log(error);
    }
}