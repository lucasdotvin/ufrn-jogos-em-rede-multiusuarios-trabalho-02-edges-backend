package broadcast

import (
	"encoding/json"
	"strings"
	"trabalho-02-edges/lib/domain/entity"
	roomresponse "trabalho-02-edges/lib/presentation/response/room"
	userresponse "trabalho-02-edges/lib/presentation/response/user"
)

const (
	GlobalRoomsEventChannel = "rooms/events"
	RoomEventsChannel       = "rooms/{room}/events"
)

const (
	RoomCreatedEventId       = EventId("room_created")
	RoomDeletedEventId       = EventId("room_deleted")
	RoomStartedEventId       = EventId("room_started")
	RoomUserIngressedEventId = EventId("room_user_ingressed")
	RoomUserEgressedEventId  = EventId("room_user_egressed")
	RoomUserReadyEventId     = EventId("room_user_ready")
	UserWonEventId           = EventId("user_won")
)

type RoomCreatedPayload struct {
	*roomresponse.RoomResponse
}

func NewRoomCreatedEvent(room *entity.Room, createdBy *entity.User) *Event[*RoomCreatedPayload] {
	payload := &RoomCreatedPayload{
		RoomResponse: roomresponse.NewRoomResponse(room, createdBy, nil, nil),
	}

	return NewEvent[*RoomCreatedPayload](RoomCreatedEventId, payload)
}

type RoomDeletedPayload struct {
	UUID string `json:"uuid"`
}

func NewRoomDeletedEvent(room *entity.Room) *Event[*RoomDeletedPayload] {
	payload := &RoomDeletedPayload{
		UUID: room.UUID,
	}

	return NewEvent[*RoomDeletedPayload](RoomDeletedEventId, payload)
}

type NewRoomStartedPayload struct {
	UUID      string `json:"uuid"`
	StartedAt string `json:"started_at"`
}

func NewRoomStartedEvent(room *entity.Room) *Event[*NewRoomStartedPayload] {
	payload := &NewRoomStartedPayload{
		UUID:      room.UUID,
		StartedAt: room.StartedAt.Format("2006-01-02T15:04:05Z"),
	}

	return NewEvent[*NewRoomStartedPayload](RoomStartedEventId, payload)
}

type RoomUserIngressedPayload struct {
	*userresponse.UserField
}

func NewRoomUserIngressedEvent(user *entity.User) *Event[*RoomUserIngressedPayload] {
	payload := &RoomUserIngressedPayload{
		UserField: userresponse.NewUserField(user, nil),
	}

	return NewEvent[*RoomUserIngressedPayload](RoomUserIngressedEventId, payload)
}

type RoomUserEgressedPayload struct {
	*userresponse.UserField
}

func NewRoomUserEgressedEvent(user *entity.User) *Event[*RoomUserEgressedPayload] {
	payload := &RoomUserEgressedPayload{
		UserField: userresponse.NewUserField(user, nil),
	}

	return NewEvent[*RoomUserEgressedPayload](RoomUserEgressedEventId, payload)
}

type NewRoomUserReadyPayload struct {
	*userresponse.UserField
}

func NewRoomUserReadyEvent(user *entity.User) *Event[*NewRoomUserReadyPayload] {
	payload := &NewRoomUserReadyPayload{
		UserField: userresponse.NewUserField(user, nil),
	}

	return NewEvent[*NewRoomUserReadyPayload](RoomUserReadyEventId, payload)
}

type NewUserWonEventPayload struct {
	*userresponse.UserField
}

func NewUserWonEvent(user *entity.User) *Event[*NewUserWonEventPayload] {
	payload := &NewUserWonEventPayload{
		UserField: userresponse.NewUserField(user, nil),
	}

	return NewEvent[*NewUserWonEventPayload](UserWonEventId, payload)
}

type RoomBroadcast struct {
	broadcaster Driver
}

func NewRoomBroadcast(broadcaster Driver) *RoomBroadcast {
	return &RoomBroadcast{
		broadcaster,
	}
}

func (r *RoomBroadcast) NotifyRoomCreated(room *entity.Room, createdBy *entity.User) {
	event := NewRoomCreatedEvent(room, createdBy)
	parsedEvent, _ := json.Marshal(event)

	r.broadcaster.Broadcast(GlobalRoomsEventChannel, parsedEvent)
}

func (r *RoomBroadcast) NotifyRoomDeleted(room *entity.Room) {
	event := NewRoomDeletedEvent(room)
	parsedEvent, _ := json.Marshal(event)

	r.broadcaster.Broadcast(GlobalRoomsEventChannel, parsedEvent)
}

func (r *RoomBroadcast) NotifyRoomStarted(room *entity.Room) {
	roomChannel := strings.Replace(RoomEventsChannel, "{room}", room.UUID, 1)

	event := NewRoomStartedEvent(room)
	parsedEvent, _ := json.Marshal(event)

	r.broadcaster.Broadcast(GlobalRoomsEventChannel, parsedEvent)
	r.broadcaster.Broadcast(roomChannel, parsedEvent)
}

func (r *RoomBroadcast) NotifyUserIngressed(room *entity.Room, createdBy *entity.User) {
	roomChannel := strings.Replace(RoomEventsChannel, "{room}", room.UUID, 1)

	event := NewRoomUserIngressedEvent(createdBy)
	parsedEvent, _ := json.Marshal(event)

	r.broadcaster.Broadcast(roomChannel, parsedEvent)
}

func (r *RoomBroadcast) NotifyUserEgressed(room *entity.Room, user *entity.User) {
	roomChannel := strings.Replace(RoomEventsChannel, "{room}", room.UUID, 1)

	event := NewRoomUserEgressedEvent(user)
	parsedEvent, _ := json.Marshal(event)

	r.broadcaster.Broadcast(roomChannel, parsedEvent)
}

func (r *RoomBroadcast) NotifyUserReady(room *entity.Room, user *entity.User) {
	roomChannel := strings.Replace(RoomEventsChannel, "{room}", room.UUID, 1)

	event := NewRoomUserReadyEvent(user)
	parsedEvent, _ := json.Marshal(event)

	r.broadcaster.Broadcast(roomChannel, parsedEvent)
}

func (r *RoomBroadcast) NotifyUserWon(room *entity.Room, user *entity.User) {
	roomChannel := strings.Replace(RoomEventsChannel, "{room}", room.UUID, 1)

	event := NewUserWonEvent(user)
	parsedEvent, _ := json.Marshal(event)

	r.broadcaster.Broadcast(roomChannel, parsedEvent)
}
