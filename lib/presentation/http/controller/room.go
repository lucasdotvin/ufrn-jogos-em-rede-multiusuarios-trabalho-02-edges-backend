package controller

import (
	"encoding/json"
	"net/http"
	roomCase "trabalho-02-edges/lib/application/case/room"
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/presentation/auth"
	roomrequest "trabalho-02-edges/lib/presentation/request/room"
	"trabalho-02-edges/lib/presentation/response/room"
)

type RoomController struct {
	listRoomUseCase    *roomCase.ListRoomUseCase
	findUserActiveRoom *roomCase.FindUserActiveRoomUseCase
	createRoomUseCase  *roomCase.CreateRoomUseCase
	ingressUserUseCase *roomCase.IngressUserUseCase
}

func NewRoomController(
	listUseCase *roomCase.ListRoomUseCase,
	findUserActiveRoom *roomCase.FindUserActiveRoomUseCase,
	createUseCase *roomCase.CreateRoomUseCase,
	ingressUseCase *roomCase.IngressUserUseCase,
) *RoomController {
	return &RoomController{
		listUseCase,
		findUserActiveRoom,
		createUseCase,
		ingressUseCase,
	}
}

func (c *RoomController) Index(w http.ResponseWriter, _ *http.Request) {
	rs, us, err := c.listRoomUseCase.Execute()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parsedResponse := room.NewRoomsResponse(rs, us)

	w.Header().Set("Content-Type", JsonContentType)
	_ = json.NewEncoder(w).Encode(parsedResponse)
}

func (c *RoomController) FindMyActiveRoom(w http.ResponseWriter, r *http.Request) {
	userUuid, ok := auth.GetUserUuid(r.Context())

	if !ok {
		http.Error(w, entity.UserNotFoundError.Error(), http.StatusUnauthorized)
		return
	}

	ro, cr, uss, rus, err := c.findUserActiveRoom.Execute(userUuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedResponse := room.NewRoomResponse(ro, cr, uss, rus)

	w.Header().Set("Content-Type", JsonContentType)
	_ = json.NewEncoder(w).Encode(parsedResponse)
}

func (c *RoomController) Store(w http.ResponseWriter, r *http.Request) {
	var input roomrequest.StoreRequest
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userUuid, ok := auth.GetUserUuid(r.Context())

	if !ok {
		http.Error(w, entity.UserNotFoundError.Error(), http.StatusUnauthorized)
		return
	}

	ro, us, err := c.createRoomUseCase.Execute(input.Name, userUuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedResponse := room.NewRoomResponse(ro, us, nil, nil)

	w.Header().Set("Content-Type", JsonContentType)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(parsedResponse)
}

func (c *RoomController) IngressUser(w http.ResponseWriter, r *http.Request) {

	roomUuid := r.PathValue("room")

	if roomUuid == "" {
		http.Error(w, entity.RoomNotFoundError.Error(), http.StatusBadRequest)
		return
	}

	userUuid, ok := auth.GetUserUuid(r.Context())

	if !ok {
		http.Error(w, entity.UserNotFoundError.Error(), http.StatusUnauthorized)
		return
	}

	ro, uss, rus, err := c.ingressUserUseCase.Execute(roomUuid, userUuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedResponse := room.NewRoomResponse(ro, nil, uss, rus)

	w.Header().Set("Content-Type", JsonContentType)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(parsedResponse)
}
