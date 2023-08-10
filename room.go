package room

// IRoom is a room interface
type IRoom interface {

	// ID returns room id
	ID() string

	// Init initializes room
	Init() error

	// Enter adds user to room
	Enter(user IUser) error

	// Leave removes user from room
	Leave(user IUser) error

	// Broadcast broadcasts message to all users in room
	Broadcast(data IMessage) error

	// Close closes room
	Close() error
}
