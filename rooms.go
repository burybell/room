package room

// IRooms is a rooms interface
// You can implement this interface to define your own rooms type
// It represents a collection of all rooms, and you can manage the rooms through it
type IRooms interface {

	// Init initializes rooms
	Init() error

	// Room returns room by id
	Room(id string) (IRoom, error)

	// OpenRoom opens a room
	OpenRoom(id string) (IRoom, error)

	// CloseRoom closes a room
	CloseRoom(id string) error

	// Close closes rooms
	Close() error
}
