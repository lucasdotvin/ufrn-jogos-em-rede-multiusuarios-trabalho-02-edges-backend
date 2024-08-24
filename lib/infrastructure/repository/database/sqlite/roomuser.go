package sqlite

import (
	"database/sql"
	"errors"
	"time"
	"trabalho-02-edges/lib/domain/entity"
)

type RoomUserRepository struct {
	db *sql.DB
}

func NewSqliteRoomUserRepository(db *sql.DB) *RoomUserRepository {
	return &RoomUserRepository{
		db,
	}
}

func (s *RoomUserRepository) Store(roomUser *entity.RoomUser) error {
	joinedAt := time.Now()

	_, err := s.db.Exec(`
		INSERT INTO room_user (room_uuid, user_uuid, joined_at)
		VALUES (?, ?, ?)
	`, roomUser.RoomUUID, roomUser.UserUUID, joinedAt)

	if err != nil {
		return err
	}

	roomUser.JoinedAt = joinedAt

	return nil
}

func (s *RoomUserRepository) FindActiveRoomForUser(userUUID string) (*entity.RoomUser, error) {
	var roomUser entity.RoomUser

	err := s.db.QueryRow(`
		SELECT *
		FROM room_user
		WHERE
			user_uuid = ?
			AND joined_at IS NOT NULL
			AND won_at IS NULL
			AND lost_at IS NULL
			AND abandoned_at IS NULL
		LIMIT 1
	`, userUUID).Scan(s.getEntityFieldPointers(&roomUser)...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &roomUser, nil
}

func (s *RoomUserRepository) Delete(roomUser *entity.RoomUser) error {
	_, err := s.db.Exec(`
		DELETE FROM room_user
		WHERE room_uuid = ? AND user_uuid = ?
	`, roomUser.RoomUUID, roomUser.UserUUID)

	if err != nil {
		return err
	}

	return nil
}

func (s *RoomUserRepository) GetByRoomUuid(roomUUID string) ([]*entity.RoomUser, error) {
	rows, err := s.db.Query(`
		SELECT *
		FROM room_user
		WHERE room_uuid = ?
	`, roomUUID)

	if err != nil {
		return nil, err
	}

	var roomUsers []*entity.RoomUser

	for rows.Next() {
		var roomUser entity.RoomUser

		err := rows.Scan(s.getEntityFieldPointers(&roomUser)...)

		if err != nil {
			return nil, err
		}

		roomUsers = append(roomUsers, &roomUser)
	}

	err = rows.Close()

	if err != nil {
		return nil, err
	}

	return roomUsers, nil
}

func (s *RoomUserRepository) Update(roomUser *entity.RoomUser) error {
	_, err := s.db.Exec(`
		UPDATE room_user
		SET
			joined_at = ?,
			won_at = ?,
			lost_at = ?,
			abandoned_at = ?
		WHERE room_uuid = ? AND user_uuid = ?
	`, roomUser.JoinedAt, roomUser.WonAt, roomUser.LostAt, roomUser.AbandonedAt, roomUser.RoomUUID, roomUser.UserUUID)

	if err != nil {
		return err
	}

	return nil
}

func (s *RoomUserRepository) FindRandomActivePlayer(roomUUID string) (*entity.RoomUser, error) {
	var roomUser entity.RoomUser

	err := s.db.QueryRow(`
		SELECT *
		FROM room_user
		WHERE
			room_uuid = ?
			AND joined_at IS NOT NULL
			AND won_at IS NULL
			AND lost_at IS NULL
			AND abandoned_at IS NULL
		ORDER BY RANDOM()
		LIMIT 1
	`, roomUUID).Scan(s.getEntityFieldPointers(&roomUser)...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &roomUser, nil
}

func (s *RoomUserRepository) getEntityFieldPointers(roomUser *entity.RoomUser) []any {
	return []any{
		&roomUser.RoomUUID,
		&roomUser.UserUUID,
		&roomUser.JoinedAt,
		&roomUser.WonAt,
		&roomUser.LostAt,
		&roomUser.AbandonedAt,
	}
}
