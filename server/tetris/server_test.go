package server

import (
    . "github.com/smartystreets/goconvey/convey"
    "github.com/gorilla/websocket"
    "net/http"
    "testing"
    "log"
)

func TestTetrisGame(t *testing.T) {
    Convey("Given a tetris game server", t, func() {
        go func() {
            server := NewTetrisServer(8000)
            server.Listen()
            defer server.Stop()
        }()

        Convey("explore two clients' behavior", func() {
            dialer := websocket.Dialer{}
            conn, _, err := dialer.Dial("ws://localhost:8000/tetris", http.Header{})
            So(err, ShouldBeNil)
            
            messageType, message, messageErr := conn.ReadMessage()
            So(messageErr, ShouldBeNil)
            log.Printf("C1 <- S: (%d) %v", messageType, string(message))

            messageErr = conn.WriteMessage(0, []byte("First response from client"))

        })
    })}
        