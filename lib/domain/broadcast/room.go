package broadcast

import "trabalho-02-edges/lib/domain/entity"

type RoomBroadcast interface {
	NotifyRoomCreated(room *entity.Room, createdBy *entity.User)

	NotifyUserIngressed(room *entity.Room, createdBy *entity.User)

	NotifyUserEgressed(room *entity.Room, user *entity.User)

	NotifyUserReady(room *entity.Room, user *entity.User)

	NotifyRoomDeleted(room *entity.Room)

	NotifyRoomStarted(room *entity.Room)

	NotifyUserWon(room *entity.Room, user *entity.User)
}
