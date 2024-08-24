package user

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/presentation/response/roomuser"
)

type UserResponse struct {
	UUID      string                  `json:"uuid"`
	Name      string                  `json:"name"`
	Username  string                  `json:"username"`
	CreatedAt string                  `json:"created_at"`
	UpdatedAt *string                 `json:"updated_at"`
	RoomUser  *roomuser.RoomUserField `json:"room_user"`
}

func NewUserResponse(user *entity.User, roomUser *entity.RoomUser) *UserResponse {
	response := &UserResponse{
		UUID:      user.UUID,
		Name:      user.Name,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		RoomUser:  roomuser.NewRoomUserField(roomUser),
	}

	if user.UpdatedAt != nil {
		updatedAt := user.UpdatedAt.Format("2006-01-02 15:04:05")
		response.UpdatedAt = &updatedAt
	}

	return response
}
