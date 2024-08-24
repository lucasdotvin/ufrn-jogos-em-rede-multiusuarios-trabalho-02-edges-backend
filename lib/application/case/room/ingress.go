package room

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/service"
)

type IngressUserUseCase struct {
	roomService *service.RoomService
	userService *service.UserService
}

func NewIngressUserUseCase(roomService *service.RoomService, userService *service.UserService) *IngressUserUseCase {
	return &IngressUserUseCase{
		roomService,
		userService,
	}
}

func (u *IngressUserUseCase) Execute(roomUuid string, userUuid string) (*entity.Room, []*entity.User, []*entity.RoomUser, error) {
	room, err := u.roomService.IngressUser(roomUuid, userUuid)

	if err != nil {
		return nil, nil, nil, err
	}

	roomUsers, err := u.roomService.GetRoomUsers(roomUuid)

	if err != nil {
		return nil, nil, nil, err
	}

	roomUsersUuids := make([]string, 0, len(roomUsers))

	for _, roomUser := range roomUsers {
		roomUsersUuids = append(roomUsersUuids, roomUser.UserUUID)
	}

	users, err := u.userService.GetWhereUuidIn(roomUsersUuids)

	if err != nil {
		return nil, nil, nil, err
	}

	return room, users, roomUsers, nil
}
