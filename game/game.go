package game

import (
	"game/room"
	"log"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func New() *Game {
	return &Game{
		Rooms:    make(map[string]*room.Room),
		MaxRooms: 2,
		Register: make(chan *websocket.Conn),
	}
}

type Game struct {
	Rooms    map[string]*room.Room
	MaxRooms int
	Register chan *websocket.Conn
}

func (g *Game) Run() {
	for {
		conn := <-g.Register
		log.Printf("got new connection")

		g.ProcessConn(conn)
	}
}

func (g *Game) FindRoom() *room.Room {
	for _, r := range g.Rooms {
		if len(r.Players) < r.MaxPlayers {
			return r
		}
	}

	if len(g.Rooms) >= g.MaxRooms {
		return nil
	}

	r := room.New()
	go r.ListenToPlayers()
	g.Rooms[r.ID] = r
	log.Printf("room %s created", r.ID)

	return r
}

func (g *Game) ProcessConn(conn *websocket.Conn) {
	id := uuid.NewV4().String()
	p := &room.Player{
		Conn: conn,
		ID:   id,
	}

	r := g.FindRoom()
	if r == nil {
		return
	}
	r.Players[p.ID] = p
	p.Room = r
	log.Printf("player %s joined room %s", p.ID, r.ID)
	go p.Listen()

	if len(r.Players) == r.MaxPlayers {
		go r.Run()
	}

}
