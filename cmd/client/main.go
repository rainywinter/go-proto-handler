package main

import (
	"conn-proto-handler/cmd/client/service"
	"conn-proto-handler/codec"
	"conn-proto-handler/delivery"
	"conn-proto-handler/server"
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var id int

func init() {
	flag.IntVar(&id, "id", 1, "client id")
}

func main() {
	flag.Parse()

	wsConn, _, err := websocket.DefaultDialer.DialContext(context.Background(), "ws://localhost:8000/ws", http.Header{})
	if err != nil {
		panic(err)
	}
	conn := server.NewConn(codec.NewProtoCodec(), delivery.NewWebsocketConn(wsConn), nil)
	svc := service.NewService(int32(id), conn)

	go func() {
		for {
			var cmd string
			fmt.Println("Please enter your cmd: ")
			fmt.Scanln(&cmd)
			if cmd == "login" {
				svc.Login(context.Background())
			} else if cmd == "exit" {
				break
			} else {
				fmt.Println("unsupport cmd, try again.")
			}

		}
	}()

	conn.Run()

	fmt.Println("conn exit")
}
