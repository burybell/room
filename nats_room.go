package room

import (
	"errors"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

// memory room
var _ IRoom = (*NatsRoom[IMessage])(nil)

// NatsRoom is a room by nats implement IRoom interface
type NatsRoom[T IMessage] struct {
	id    string
	users map[string]IUser

	pushCh  chan IMessage
	msgCh   chan T
	enterCh chan IUser
	leaveCh chan IUser
	closeCh chan bool

	natConn      *nats.Conn
	encodedConn  *nats.EncodedConn
	subscription *nats.Subscription
}

func NewNatsRoom[T IMessage](id string, natConn *nats.Conn) IRoom {
	return &NatsRoom[T]{id: id, natConn: natConn}
}

func (r *NatsRoom[T]) ID() string {
	return r.id
}

func (r *NatsRoom[T]) Init() (err error) {
	if r.id == "" {
		return errors.New("room id is empty")
	}
	r.users = make(map[string]IUser)
	r.pushCh = make(chan IMessage, 1000)
	r.msgCh = make(chan T)
	r.enterCh = make(chan IUser, 1)
	r.leaveCh = make(chan IUser, 1)

	r.encodedConn, err = nats.NewEncodedConn(r.natConn, nats.JSON_ENCODER)
	if err != nil {
		return err
	}

	r.subscription, err = r.encodedConn.BindRecvChan(r.id, r.msgCh)
	if err != nil {
		return err
	}

	go r.init()
	return nil
}

func (r *NatsRoom[T]) init() {
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
		case msg := <-r.pushCh:
			err := r.encodedConn.Publish(r.id, msg)
			if err != nil {
				log.Printf("publish message error: %v", err)
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
			err := r.subscription.Unsubscribe()
			if err != nil {
				log.Printf("unsubscribe error: %v", err)
			}
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

func (r *NatsRoom[T]) Enter(user IUser) error {
	if r.users[user.ID()] != nil {
		err := r.users[user.ID()].Close()
		if err != nil {
			return err
		}
	}
	r.enterCh <- user
	return nil
}

func (r *NatsRoom[T]) Leave(user IUser) error {
	r.leaveCh <- user
	return nil
}

func (r *NatsRoom[T]) Broadcast(data IMessage) error {
	r.pushCh <- data
	return nil
}

func (r *NatsRoom[T]) Close() error {
	r.closeCh <- true
	return nil
}
