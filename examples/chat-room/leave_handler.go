package main

import (
	"github.com/burybell/room"
	"github.com/gorilla/mux"
	"net/http"
)

func LeaveHandler(w http.ResponseWriter, r *http.Request) {

	var roomID = mux.Vars(r)["room_id"]
	rm, err := rooms.Room(roomID)
	if err != nil {
		ResponseErr(w, 400, err)
		return
	}

	var userID = mux.Vars(r)["user_id"]
	err = rm.Leave(room.NewUser(userID, nil))
	if err != nil {
		ResponseErr(w, 500, err)
		return
	}
}
