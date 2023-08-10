package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type ReqPublish struct {
	Message string `json:"message,omitempty"`
}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	var req ReqPublish
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ResponseErr(w, 400, err)
		return
	}

	var roomID = mux.Vars(r)["room_id"]
	rm, err := rooms.Room(roomID)
	if err != nil {
		ResponseErr(w, 500, err)
		return
	}

	var userID = mux.Vars(r)["user_id"]
	err = rm.Broadcast(RoomEvent{
		EventType: UserMsg,
		EventData: map[string]any{
			"message": req.Message,
			"userID":  userID,
		},
	})
	if err != nil {
		ResponseErr(w, 500, err)
		return
	}
}
