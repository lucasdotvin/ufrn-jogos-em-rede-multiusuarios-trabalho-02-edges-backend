package room

import (
	"trabalho-02-edges/config"
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/service"
)

type CreateRoomUseCase struct {
	cfg         config.Config
	roomService *service.RoomService
}

func NewCreateRoomUseCase(cfg config.Config, roomService *service.RoomService) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		cfg,
		roomService,
	}
}

func (u *CreateRoomUseCase) Execute(name string, creatorUuid string) (*entity.Room, *entity.User, error) {
	ro, us, err := u.roomService.Create(name, u.cfg.DefaultMaxPlayers, creatorUuid)

	if err != nil {
		return nil, nil, err
	}

	return ro, us, nil
}
