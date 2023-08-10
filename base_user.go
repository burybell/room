package room

// User is a user implementation of IUser
type User struct {
	id   string
	conn IConn
}

func NewUser(id string, conn IConn) *User {
	return &User{id: id, conn: conn}
}

func (r *User) ID() string {
	return r.id
}

func (r *User) Close() error {
	return r.conn.Close()
}

func (r *User) Conn() IConn {
	return r.conn
}
