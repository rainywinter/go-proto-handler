package main

import (
	"conn-proto-handler/cmd/server/service"
	"conn-proto-handler/server"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const host = ":8000"

func main() {
	up := websocket.Upgrader{}
	svc := service.NewService(up)

	r := mux.NewRouter()
	r.Handle("/ws", svc)

	go func() {
		for {
			var cmd string
			var id uint32
			var msg string
			fmt.Println("Please enter your cmd 、conn id 、 msg: ")
			fmt.Scanln(&cmd, &id, &msg)
			if cmd == "hello" {
				svc.Hello(context.Background(), server.ConnID(id), msg)
			} else if cmd == "exit" {
				svc.Shutdown()
				break
			} else {
				fmt.Println("unsupport cmd, try again.")
			}
		}
	}()

	fmt.Println("start ", host)
	http.ListenAndServe(host, r)
}
