package room

import "trabalho-02-edges/lib/domain/entity"

type RoomsResponse struct {
	Rooms []*RoomResponse `json:"rooms"`
}

func NewRoomsResponse(rooms []*entity.Room, creators []*entity.User) *RoomsResponse {
	creatorsByUuid := make(map[string]*entity.User)

	for _, creator := range creators {
		creatorsByUuid[creator.UUID] = creator
	}

	roomsResponse := &RoomsResponse{
		Rooms: make([]*RoomResponse, 0, len(rooms)),
	}

	for _, room := range rooms {
		roomsResponse.Rooms = append(roomsResponse.Rooms, NewRoomResponse(
			room,
			creatorsByUuid[room.CreatedBy],
			nil,
			nil,
		))
	}

	return roomsResponse
}
