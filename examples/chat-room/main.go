package main

import (
	"github.com/burybell/room"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

const (
	UserEnter = "user_enter"
	UserLeave = "user_leave"
	UserMsg   = "user_msg"
)

type RoomEvent struct {
	EventType string         `json:"eventType"`
	EventData map[string]any `json:"eventData"`
}

func (u RoomEvent) Bytes() ([]byte, error) {
	return room.DefaultMarshal(u)
}

var rooms room.IRooms

func init() {

	// connect to nats server
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	// create rooms
	rooms = room.NewNatsRooms(conn, func(id string, conn *nats.Conn) room.IRoom {
		return room.NewNatsRoom[RoomEvent](id, conn)
	})

	// init rooms
	err = rooms.Init()
	if err != nil {
		panic(err)
	}
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/rooms/{room_id}/enter/{user_id}", EnterHandler)
	r.HandleFunc("/rooms/{room_id}/leave/{user_id}", LeaveHandler)
	r.HandleFunc("/rooms/{room_id}/publish/{user_id}", PublishHandler)
	r.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./examples/chat-room/index.html")
	})
	log.Fatalln(http.ListenAndServe(":8080", r))
}
