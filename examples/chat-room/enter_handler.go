package main

import (
	"github.com/burybell/room"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

func EnterHandler(w http.ResponseWriter, r *http.Request) {

	var roomID = mux.Vars(r)["room_id"]
	rm, err := rooms.Room(roomID)
	if err != nil {
		rm, err = rooms.OpenRoom(roomID)
		if err != nil {
			ResponseErr(w, 400, err)
			return
		}
	}

	upgrade := websocket.Upgrader{}
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		ResponseErr(w, 400, err)
		return
	}

	var userID = mux.Vars(r)["user_id"]
	err = rm.Enter(room.NewUser(userID, room.NewConnWS(conn)))
	if err != nil {
		ResponseErr(w, 500, err)
		return
	}

	err = rm.Broadcast(RoomEvent{EventType: UserEnter, EventData: map[string]any{"userID": userID}})
	if err != nil {
		ResponseErr(w, 500, err)
		return
	}
}
