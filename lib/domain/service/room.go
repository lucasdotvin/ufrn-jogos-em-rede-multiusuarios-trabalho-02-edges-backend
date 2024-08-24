package service

import (
	"time"
	"trabalho-02-edges/lib/domain/broadcast"
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/repository/database"
)

type RoomService struct {
	roomBroadcast broadcast.RoomBroadcast

	roomDatabaseRepository     database.RoomRepository
	roomUserDatabaseRepository database.RoomUserRepository
	userDatabaseRepository     database.UserRepository
}

func NewRoomService(
	roomBroadcast broadcast.RoomBroadcast,
	roomDatabaseRepository database.RoomRepository,
	roomUserDataRepository database.RoomUserRepository,
	userDatabaseRepository database.UserRepository,
) *RoomService {
	return &RoomService{
		roomBroadcast,
		roomDatabaseRepository,
		roomUserDataRepository,
		userDatabaseRepository,
	}
}

func (s *RoomService) List() ([]*entity.Room, error) {
	return s.roomDatabaseRepository.GetAllOpen()
}

func (s *RoomService) Create(name string, maxPlayers int, createdBy string) (*entity.Room, *entity.User, error) {
	u, err := s.userDatabaseRepository.FindByUuid(createdBy)

	if err != nil {
		return nil, nil, err
	}

	if u == nil {
		return nil, nil, entity.UserNotFoundError
	}

	alreadyActiveRoom, err := s.roomUserDatabaseRepository.FindActiveRoomForUser(createdBy)

	if err != nil {
		return nil, nil, err
	}

	if alreadyActiveRoom != nil {
		return nil, nil, entity.UserAlreadyInRoomError
	}

	r := &entity.Room{
		Name:           name,
		MaxPlayers:     maxPlayers,
		CurrentPlayers: 1,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}

	err = s.roomDatabaseRepository.Store(r)

	if err != nil {
		return nil, nil, err
	}

	ru := &entity.RoomUser{
		RoomUUID: r.UUID,
		UserUUID: r.CreatedBy,
	}

	err = s.roomUserDatabaseRepository.Store(ru)

	if err != nil {
		return nil, nil, err
	}

	s.roomBroadcast.NotifyRoomCreated(r, u)

	return r, u, nil
}

func (s *RoomService) IngressUser(roomUuid string, userUuid string) (*entity.Room, error) {
	ro, err := s.roomDatabaseRepository.FindByUuid(roomUuid)

	if err != nil {
		return nil, err
	}

	if ro == nil {
		return nil, entity.RoomNotFoundError
	}

	if !ro.CanIngress() {
		return nil, entity.RoomIsFullError
	}

	us, err := s.userDatabaseRepository.FindByUuid(userUuid)

	if err != nil {
		return nil, err
	}

	if us == nil {
		return nil, entity.UserNotFoundError
	}

	alreadyActiveRoom, err := s.roomUserDatabaseRepository.FindActiveRoomForUser(userUuid)

	if err != nil {
		return nil, err
	}

	if alreadyActiveRoom != nil {
		return nil, entity.UserAlreadyInRoomError
	}

	ru := &entity.RoomUser{
		RoomUUID: ro.UUID,
		UserUUID: us.UUID,
	}

	err = s.roomUserDatabaseRepository.Store(ru)

	if err != nil {
		return nil, err
	}

	ro.CurrentPlayers++

	err = s.roomDatabaseRepository.Update(ro)

	if err != nil {
		return nil, err
	}

	s.roomBroadcast.NotifyUserIngressed(ro, us)

	return ro, nil
}

func (s *RoomService) ValidateUserIsPresent(roomUuid string, userUuid string) error {
	ro, err := s.roomDatabaseRepository.FindByUuid(roomUuid)

	if err != nil {
		return err
	}

	if ro == nil {
		return entity.RoomNotFoundError
	}

	us, err := s.userDatabaseRepository.FindByUuid(userUuid)

	if err != nil {
		return err
	}

	if us == nil {
		return entity.UserNotFoundError
	}

	ru, err := s.roomUserDatabaseRepository.FindActiveRoomForUser(userUuid)

	if err != nil {
		return err
	}

	if ru == nil {
		return entity.UserNotInRoomError
	}

	if ru.RoomUUID != ro.UUID {
		return entity.UserNotInRoomError
	}

	return nil
}

func (s *RoomService) HandleUserDisconnect(userUuid string) error {
	ru, err := s.roomUserDatabaseRepository.FindActiveRoomForUser(userUuid)

	if err != nil {
		return err
	}

	if ru == nil {
		return nil
	}

	ro, err := s.roomDatabaseRepository.FindByUuid(ru.RoomUUID)

	if err != nil {
		return err
	}

	if ro == nil {
		return nil
	}

	if !ro.IsStarted() {
		return s.handleUserDisconnectBeforeGameStart(ru, ro)
	}

	return nil
}

func (s *RoomService) handleUserDisconnectBeforeGameStart(ru *entity.RoomUser, ro *entity.Room) error {
	err := s.roomUserDatabaseRepository.Delete(ru)

	if err != nil {
		return err
	}

	ro.CurrentPlayers--

	if ro.IsEmpty() {
		return s.deleteRoom(ro)
	}

	err = s.roomDatabaseRepository.Update(ro)

	if err != nil {
		return err
	}

	us, err := s.userDatabaseRepository.FindByUuid(ru.UserUUID)

	if err != nil {
		return err
	}

	if us == nil {
		return entity.UserNotFoundError
	}

	s.roomBroadcast.NotifyUserEgressed(ro, us)

	return nil
}

func (s *RoomService) FindUserActiveRoom(userUuid string) (*entity.Room, error) {
	ru, err := s.roomUserDatabaseRepository.FindActiveRoomForUser(userUuid)

	if err != nil {
		return nil, err
	}

	if ru == nil {
		return nil, nil
	}

	ro, err := s.roomDatabaseRepository.FindByUuid(ru.RoomUUID)

	if err != nil {
		return nil, err
	}

	if ro == nil {
		return nil, entity.RoomNotFoundError
	}

	if !ro.IsActive() {
		return nil, entity.RoomNotFoundError
	}

	return ro, nil
}

func (s *RoomService) deleteRoom(room *entity.Room) error {
	err := s.roomDatabaseRepository.Delete(room)

	if err != nil {
		return err
	}

	s.roomBroadcast.NotifyRoomDeleted(room)

	return nil
}

func (s *RoomService) GetRoomUsers(roomUuid string) ([]*entity.RoomUser, error) {
	rus, err := s.roomUserDatabaseRepository.GetByRoomUuid(roomUuid)

	if err != nil {
		return nil, err
	}

	return rus, nil
}
