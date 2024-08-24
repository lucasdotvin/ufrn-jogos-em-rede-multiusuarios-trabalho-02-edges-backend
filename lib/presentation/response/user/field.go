package user

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/presentation/response/roomuser"
)

type UserField struct {
	UUID     string                  `json:"uuid"`
	Name     string                  `json:"name"`
	Username string                  `json:"username"`
	RoomUser *roomuser.RoomUserField `json:"room_user"`
}

func NewUserField(user *entity.User, roomUser *entity.RoomUser) *UserField {
	if user == nil {
		return nil
	}

	return &UserField{
		UUID:     user.UUID,
		Name:     user.Name,
		Username: user.Username,
		RoomUser: roomuser.NewRoomUserField(roomUser),
	}
}

func NewUsersField(users []*entity.User, roomUsers []*entity.RoomUser) []*UserField {
	if users == nil {
		return nil
	}

	roomUsersByUserUuid := make(map[string]*entity.RoomUser)

	for _, roomUser := range roomUsers {
		roomUsersByUserUuid[roomUser.UserUUID] = roomUser
	}

	userFields := make([]*UserField, 0)

	for _, user := range users {
		userFields = append(userFields, NewUserField(user, roomUsersByUserUuid[user.UUID]))
	}

	return userFields
}
