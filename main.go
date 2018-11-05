package main

import (
	"game/game"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {

	g := game.New()
	go g.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("open connection")

		// check auth

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("cannot upgrade connection: %s", err)
		}

		g.Register <- conn
	})

	http.ListenAndServe(":9090", nil)
}
