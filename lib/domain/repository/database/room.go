package database

import "trabalho-02-edges/lib/domain/entity"

type RoomRepository interface {
	GetAllOpen() ([]*entity.Room, error)

	Store(room *entity.Room) error

	Update(room *entity.Room) error

	FindByUuid(uuid string) (*entity.Room, error)

	Delete(room *entity.Room) error
}
