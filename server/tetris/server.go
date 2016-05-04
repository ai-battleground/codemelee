package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    "log"
    "fmt"
    "time"
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
    CheckOrigin: func(r *http.Request) bool { return true },
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if (err != nil) {
        log.Println("Error setting up websocket: %v", err)
    }
    go func() {
        if err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, who's there?")); err != nil {
            log.Printf("Error sending hello: %v", err)
        }
        time.Sleep(500 * time.Millisecond)
        if err = conn.WriteMessage(websocket.TextMessage, []byte("Would you like to play a game?")); err != nil {
            log.Printf("Error sending GTW proposition: %v", err)
        }
        time.Sleep(3000 * time.Millisecond)
        if err = conn.WriteMessage(websocket.TextMessage, []byte("How about a nice game of chess?")); err != nil {
            log.Printf("Error sending chess proposition: %v", err)
        }
    }()
}

func (s *TetrisGameServer) Listen() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/tetris", gameHandler)
    http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func NewTetrisServer(port int) *TetrisGameServer {
    s := TetrisGameServer{port: port}
    return &s
}