package room

import (
	"errors"
	"log"
	"sync"
)

// memory rooms
var _ IRooms = (*MemoryRooms)(nil)

// MemoryRooms is a rooms in memory implement IRooms interface
// You can use room id to route connections with the same room number to the same node, such as using nginx
type MemoryRooms struct {
	NewRoom    func(id string) IRoom
	rooms      map[string]IRoom
	roomsMutex sync.RWMutex

	openCh  chan string
	closeCh chan string
}

func NewMemoryRooms(newRoom func(id string) IRoom) *MemoryRooms {
	return &MemoryRooms{NewRoom: newRoom, roomsMutex: sync.RWMutex{}}
}

func (r *MemoryRooms) Init() error {
	if r.NewRoom == nil {
		return errors.New("new room is nil")
	}

	r.rooms = make(map[string]IRoom)
	r.openCh = make(chan string, 10)
	r.closeCh = make(chan string, 10)
	go r.init()
	return nil
}

func (r *MemoryRooms) init() {
	for {
		select {
		case id := <-r.openCh:
			r.roomsMutex.Lock()
			r.rooms[id] = r.NewRoom(id)
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

func (r *MemoryRooms) Room(id string) (IRoom, error) {
	r.roomsMutex.RLock()
	defer r.roomsMutex.RUnlock()
	if room, ok := r.rooms[id]; ok {
		return room, nil
	}
	return nil, errors.New("room is not exist")
}

func (r *MemoryRooms) OpenRoom(id string) (IRoom, error) {
	r.openCh <- id
	for {
		room, err := r.Room(id)
		if err != nil {
			continue
		}
		return room, nil
	}
}

func (r *MemoryRooms) CloseRoom(id string) error {
	r.closeCh <- id
	for {
		_, err := r.Room(id)
		if err == nil {
			continue
		}
		return nil
	}
}

func (r *MemoryRooms) Close() error {
	close(r.openCh)
	close(r.closeCh)
	return nil
}
