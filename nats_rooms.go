package room

import (
	"errors"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
)

var _ IRooms = (*NatsRooms)(nil)

// NatsRooms is a rooms by nats implement IRooms interface
// nats is a message queue, so we don't need to use channel to send message
type NatsRooms struct {
	NewRoom    func(id string, conn *nats.Conn) IRoom
	rooms      map[string]IRoom
	roomsMutex sync.RWMutex

	openCh  chan string
	closeCh chan string

	natsConn          *nats.Conn
	openSubscription  *nats.Subscription
	closeSubscription *nats.Subscription
}

func NewNatsRooms(natsConn *nats.Conn, newRoom func(id string, conn *nats.Conn) IRoom) *NatsRooms {
	return &NatsRooms{NewRoom: newRoom, natsConn: natsConn, roomsMutex: sync.RWMutex{}}
}

func (r *NatsRooms) Init() (err error) {
	if r.NewRoom == nil {
		return errors.New("new room is nil")
	}

	r.openSubscription, err = r.natsConn.Subscribe("room_open", func(msg *nats.Msg) {
		r.openCh <- string(msg.Data)
	})
	if err != nil {
		return err
	}

	r.closeSubscription, err = r.natsConn.Subscribe("room_close", func(msg *nats.Msg) {
		r.openCh <- string(msg.Data)
	})
	if err != nil {
		return err
	}

	r.rooms = make(map[string]IRoom)
	r.openCh = make(chan string, 1000)
	r.closeCh = make(chan string, 1000)
	go r.init()
	return nil
}

func (r *NatsRooms) init() {
	for {
		select {
		case id := <-r.openCh:
			r.roomsMutex.Lock()
			r.rooms[id] = r.NewRoom(id, r.natsConn)
			err := r.rooms[id].Init()
			if err != nil {
				log.Printf("init room error: %v", err)
				delete(r.rooms, id)
			}
			r.roomsMutex.Unlock()
		case id := <-r.closeCh:
			r.roomsMutex.Lock()
			err := r.rooms[id].Close()
			if err != nil {
				log.Printf("close room error: %v", err)
			}
			delete(r.rooms, id)
			r.roomsMutex.Unlock()
		}
	}
}

func (r *NatsRooms) Room(id string) (IRoom, error) {
	r.roomsMutex.RLock()
	defer r.roomsMutex.RUnlock()
	if room, ok := r.rooms[id]; ok {
		return room, nil
	}
	return nil, errors.New("room is not exist")
}

func (r *NatsRooms) OpenRoom(id string) (IRoom, error) {
	err := r.natsConn.Publish("room_open", []byte(id))
	if err != nil {
		return nil, err
	}
	for {
		room, err := r.Room(id)
		if err != nil {
			continue
		}
		return room, nil
	}
}

func (r *NatsRooms) CloseRoom(id string) error {
	err := r.natsConn.Publish("room_close", []byte(id))
	if err != nil {
		return err
	}
	for {
		_, err := r.Room(id)
		if err == nil {
			continue
		}
		return nil
	}
}

func (r *NatsRooms) Close() error {

	err := r.openSubscription.Unsubscribe()
	if err != nil {
		return err
	}
	err = r.closeSubscription.Unsubscribe()
	if err != nil {
		return err
	}
	close(r.openCh)
	close(r.closeCh)
	return nil
}
