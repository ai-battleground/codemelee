package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    "encoding/json"
    "log"
    "fmt"
    "github.com/rhoegg/codemelee/engine/tetris"
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
    game := tetris.NewTetrisGame()
    game.Start()
    log.Printf("Game started")
    go func() {
        for {
            piece := <-game.PieceState
            ps := struct {
                Piece string
                Position tetris.Point
                Shape [4]tetris.Point
            }{piece.Name, piece.Position, piece.Points()}
            psJson, err := json.Marshal(ps)
            if err != nil {
                log.Printf("Error making json: %v", err)
            } else {
                if err = conn.WriteMessage(websocket.TextMessage, psJson); err != nil {
                    log.Printf("Error sending piece state: %v", err)
                }
            }
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