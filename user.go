package room

// IUser is a user interface
// You can implement this interface to define your own user type
// It represents a user in a room
type IUser interface {
	ID() string
	Close() error
	Conn() IConn
}
