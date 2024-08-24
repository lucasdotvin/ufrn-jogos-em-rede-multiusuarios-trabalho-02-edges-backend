package room

import (
	"trabalho-02-edges/lib/domain/service"
)

type HandleUserDisconnectUseCase struct {
	roomService *service.RoomService
}

func NewHandleUserDisconnectUseCase(roomService *service.RoomService) *HandleUserDisconnectUseCase {
	return &HandleUserDisconnectUseCase{
		roomService,
	}
}

func (u *HandleUserDisconnectUseCase) Execute(userUuid string) error {
	return u.roomService.HandleUserDisconnect(userUuid)
}
