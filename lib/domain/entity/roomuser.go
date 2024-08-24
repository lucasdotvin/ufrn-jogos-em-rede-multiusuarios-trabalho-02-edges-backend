package entity

import (
	"time"
)

type RoomUser struct {
	RoomUUID    string
	UserUUID    string
	JoinedAt    time.Time
	WonAt       *time.Time
	LostAt      *time.Time
	AbandonedAt *time.Time
}
