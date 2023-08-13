package room

type RoomOptions struct {
	MsgBuffSize   int
	EnterBuffSize int
	LeaveBuffSize int
	OnUserEnter   func(room IRoom, user IUser)
	OnUserLeave   func(room IRoom, user IUser)
}

func DefaultRoomOptions() *RoomOptions {
	return &RoomOptions{
		MsgBuffSize:   1000,
		EnterBuffSize: 1000,
		LeaveBuffSize: 1000,
	}
}

type RoomOption func(opts *RoomOptions)

func SetMsgBuffSize(buffSize int) RoomOption {
	return func(opts *RoomOptions) {
		opts.MsgBuffSize = buffSize
	}
}

func SetEnterBuffSize(buffSize int) RoomOption {
	return func(opts *RoomOptions) {
		opts.EnterBuffSize = buffSize
	}
}

func SetLeaveBuffSize(buffSize int) RoomOption {
	return func(opts *RoomOptions) {
		opts.LeaveBuffSize = buffSize
	}
}

func SetOnUserEnter(onUserEnter func(room IRoom, user IUser)) RoomOption {
	return func(opts *RoomOptions) {
		opts.OnUserEnter = onUserEnter
	}
}

func SetOnUserLeave(onUserLeave func(room IRoom, user IUser)) RoomOption {
	return func(opts *RoomOptions) {
		opts.OnUserLeave = onUserLeave
	}
}
