package room

import "encoding/json"

// IMessage is a message interface
// You can implement this interface to define your own message type
// Broadcast will call Bytes() to get the message bytes
type IMessage interface {
	Bytes() ([]byte, error)
}

var (
	DefaultMarshal = json.Marshal // default marshal function, you can change it
)

// HMessage is a map[string]interface{} message
type HMessage map[string]interface{}

func (r HMessage) Bytes() ([]byte, error) {
	return DefaultMarshal(r)
}

// AMessage is a []interface{} message
type AMessage []interface{}

func (r AMessage) Bytes() ([]byte, error) {
	return DefaultMarshal(r)
}
