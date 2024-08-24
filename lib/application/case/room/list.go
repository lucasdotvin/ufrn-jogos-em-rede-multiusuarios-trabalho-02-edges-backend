package room

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/service"
)

type ListRoomUseCase struct {
	roomService *service.RoomService
	userService *service.UserService
}

func NewListRoomUseCase(roomService *service.RoomService, userService *service.UserService) *ListRoomUseCase {
	return &ListRoomUseCase{
		roomService,
		userService,
	}
}

func (u *ListRoomUseCase) Execute() ([]*entity.Room, []*entity.User, error) {
	rooms, err := u.roomService.List()

	if err != nil {
		return nil, nil, err
	}

	creatorsUuids := make([]string, 0, len(rooms))

	for _, room := range rooms {
		creatorsUuids = append(creatorsUuids, room.CreatedBy)
	}

	creators, err := u.userService.GetWhereUuidIn(creatorsUuids)

	if err != nil {
		return nil, nil, err
	}

	return rooms, creators, nil
}
