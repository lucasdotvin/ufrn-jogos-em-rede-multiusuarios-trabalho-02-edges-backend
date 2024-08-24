package controller

import (
	"net/http"
	"strings"
	roomCase "trabalho-02-edges/lib/application/case/room"
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/presentation/auth"
	"trabalho-02-edges/lib/presentation/broadcast"
)

type BroadcastController struct {
	broadcaster                  broadcast.Driver
	validateUserIsPresentUseCase *roomCase.ValidateUserIsPresent
	handleUserDisconnectUseCase  *roomCase.HandleUserDisconnectUseCase
}

func NewBroadcastController(
	broadcaster broadcast.Driver,
	validateUserIsPresentUseCase *roomCase.ValidateUserIsPresent,
	handleUserDisconnectUseCase *roomCase.HandleUserDisconnectUseCase,
) *BroadcastController {
	return &BroadcastController{
		broadcaster,
		validateUserIsPresentUseCase,
		handleUserDisconnectUseCase,
	}
}

func (c *BroadcastController) SubscribeForGlobalRoomEvents(w http.ResponseWriter, r *http.Request) {
	userUuid, ok := auth.GetUserUuid(r.Context())

	if !ok {
		http.Error(w, entity.UserNotFoundError.Error(), http.StatusUnauthorized)
		return
	}

	_ = c.broadcaster.Subscribe(broadcast.GlobalRoomsEventChannel, userUuid, w, r)
}

func (c *BroadcastController) SubscribeForRoomEvents(w http.ResponseWriter, r *http.Request) {
	roomUuid := r.PathValue("room")

	if roomUuid == "" {
		http.Error(w, entity.RoomNotFoundError.Error(), http.StatusNotFound)
		return
	}

	userUuid, ok := auth.GetUserUuid(r.Context())

	if !ok {
		http.Error(w, entity.UserNotFoundError.Error(), http.StatusUnauthorized)
		return
	}

	err := c.validateUserIsPresentUseCase.Execute(roomUuid, userUuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	channel := strings.Replace(broadcast.RoomEventsChannel, "{room}", roomUuid, 1)

	_ = c.broadcaster.Subscribe(channel, userUuid, w, r, &broadcast.Callbacker{
		OnDisconnect: func(_ string, subscriberKey string) {
			_ = c.handleUserDisconnectUseCase.Execute(subscriberKey)
		},
	})
}
