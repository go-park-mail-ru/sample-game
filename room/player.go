package room

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerData struct {
	Username string `json:"username"`
	HP       string
	Position Position `json:"position"`
}

type Player struct {
	ID   string
	Room *Room
	Conn *websocket.Conn
	Data PlayerData
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Player  *Player         `json:"-"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (p *Player) Listen() {
	log.Printf("start listening messages from player %s", p.ID)

	for {
		m := &IncomingMessage{}

		err := p.Conn.ReadJSON(m)
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("player %s was disconnected", p.ID)
			p.Room.Unregister <- p
			return
		}

		m.Player = p
		p.Room.Message <- m
	}
}

func (p *Player) Send(s *Message) {
	err := p.Conn.WriteJSON(s)
	if err != nil {
		log.Printf("cannot send state to client: %s", err)
	}
}
