package database

import "trabalho-02-edges/lib/domain/entity"

type RoomUserRepository interface {
	Store(roomUser *entity.RoomUser) error

	FindActiveRoomForUser(userUUID string) (*entity.RoomUser, error)

	Delete(roomUser *entity.RoomUser) error

	GetByRoomUuid(roomUUID string) ([]*entity.RoomUser, error)

	Update(roomUser *entity.RoomUser) error

	FindRandomActivePlayer(roomUUID string) (*entity.RoomUser, error)
}
