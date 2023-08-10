package room

import (
	"errors"
	"log"
	"time"
)

// memory room
var _ IRoom = (*MemoryRoom)(nil)

// MemoryRoom is a room in memory implement IRoom interface
type MemoryRoom struct {
	id    string
	users map[string]IUser

	msgCh   chan IMessage
	enterCh chan IUser
	leaveCh chan IUser
	closeCh chan bool
}

func NewMemoryRoom(id string) *MemoryRoom {
	return &MemoryRoom{id: id}
}

func (r *MemoryRoom) ID() string {
	return r.id
}

func (r *MemoryRoom) Init() error {
	if r.id == "" {
		return errors.New("room id is empty")
	}
	r.users = make(map[string]IUser)
	r.msgCh = make(chan IMessage, 1000)
	r.enterCh = make(chan IUser, 1)
	r.leaveCh = make(chan IUser, 1)
	go r.init()
	return nil
}

func (r *MemoryRoom) init() {
	heartbeat := time.NewTimer(time.Second)
	for {
		select {
		case msg := <-r.msgCh:
			for _, user := range r.users {
				err := user.Conn().PushMessage(msg)
				if err != nil {
					log.Printf("write message error: %v", err)
				}
			}
		case user := <-r.enterCh:
			r.users[user.ID()] = user
		case user := <-r.leaveCh:
			if r.users[user.ID()] != nil {
				err := r.users[user.ID()].Close()
				if err != nil {
					log.Printf("close connection error: %v", err)
				}
				delete(r.users, user.ID())
			}
		case <-heartbeat.C:
			for _, user := range r.users {
				err := user.Conn().Heartbeat()
				if err != nil {
					log.Printf("write message error: %v", err)
				}
			}
		case <-r.closeCh:
			close(r.msgCh)
			close(r.enterCh)
			close(r.leaveCh)
			close(r.closeCh)
			for _, user := range r.users {
				err := user.Close()
				if err != nil {
					log.Printf("close connection error: %v", err)
				}
			}
			break
		}
	}
}

func (r *MemoryRoom) Enter(user IUser) error {
	if r.users[user.ID()] != nil {
		err := r.users[user.ID()].Close()
		if err != nil {
			return err
		}
	}
	r.enterCh <- user
	return nil
}

func (r *MemoryRoom) Leave(user IUser) error {
	r.leaveCh <- user
	return nil
}

func (r *MemoryRoom) Broadcast(data IMessage) error {
	r.msgCh <- data
	return nil
}

func (r *MemoryRoom) Close() error {
	r.closeCh <- true
	return nil
}
