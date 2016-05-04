package main

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
            defer conn.Close()
            
            go func() {
                for i:=0; i<3; i++ {
                    if messageType, message, messageErr := conn.ReadMessage(); messageErr == nil {
                        log.Printf("C1 <- S: (%d) %v", messageType, string(message))
                    } else {
                        log.Printf("C1: Error reading message %v", messageErr)
                    }
                }
            }()

            conn2, _, err := dialer.Dial("ws://localhost:8000/tetris", http.Header{})
            So(err, ShouldBeNil)

            for i := 0; i < 3; i++ {
                messageType, message, messageErr := conn2.ReadMessage()
                So(messageErr, ShouldBeNil)
                log.Printf("C2 <- S: (%d) %v", messageType, string(message))
            }
        })
    })
}

        