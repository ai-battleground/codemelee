package server

import (
    "net/http"
    "github.com/gorilla/websocket"
    "log"
)

type TetrisGameServer struct {
    port int
}

func (s *TetrisGameServer) Stop() {

}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if (err != nil) {
        log.Println("Error setting up websocket: %v", err)
    }
    if err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, who's there?")); err != nil {
        log.Printf("Error sending hello: %v", err)
    }
}

func (s *TetrisGameServer) Listen() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/tetris", gameHandler)
    http.ListenAndServe(":8000", nil)
}

func NewTetrisServer(port int) *TetrisGameServer {
    s := TetrisGameServer{port: port}
    return &s
}