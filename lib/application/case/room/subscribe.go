package room

import (
	"trabalho-02-edges/lib/domain/service"
)

type ValidateUserIsPresent struct {
	roomService *service.RoomService
}

func NewValidateUserIsPresent(roomService *service.RoomService) *ValidateUserIsPresent {
	return &ValidateUserIsPresent{
		roomService,
	}
}

func (u *ValidateUserIsPresent) Execute(roomUuid string, userUuid string) error {
	return u.roomService.ValidateUserIsPresent(roomUuid, userUuid)
}
