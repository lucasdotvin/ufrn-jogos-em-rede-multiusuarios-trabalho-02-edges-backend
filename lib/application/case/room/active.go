package room

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/service"
)

type FindUserActiveRoomUseCase struct {
	roomService *service.RoomService
	userService *service.UserService
}

func NewFindUserActiveRoomUseCase(roomService *service.RoomService, userService *service.UserService) *FindUserActiveRoomUseCase {
	return &FindUserActiveRoomUseCase{
		roomService,
		userService,
	}
}

func (u *FindUserActiveRoomUseCase) Execute(userUuid string) (*entity.Room, *entity.User, []*entity.User, []*entity.RoomUser, error) {
	ro, err := u.roomService.FindUserActiveRoom(userUuid)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	if ro == nil {
		return nil, nil, nil, nil, nil
	}

	roomUsers, err := u.roomService.GetRoomUsers(ro.UUID)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	roomUsersUuids := make([]string, 0, len(roomUsers)+1)
	roomUsersUuids = append(roomUsersUuids, ro.CreatedBy)

	for _, roomUser := range roomUsers {
		roomUsersUuids = append(roomUsersUuids, roomUser.UserUUID)
	}

	users, err := u.userService.GetWhereUuidIn(roomUsersUuids)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	var cr *entity.User

	for _, user := range users {
		if user.UUID == ro.CreatedBy {
			cr = user
			break
		}
	}

	return ro, cr, users, roomUsers, nil
}
