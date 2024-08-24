package sqlite

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
	"trabalho-02-edges/lib/domain/entity"
)

type RoomRepository struct {
	db *sql.DB
}

func NewSqliteRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{
		db,
	}
}

func (s *RoomRepository) GetAllOpen() ([]*entity.Room, error) {
	rows, err := s.db.Query("SELECT * FROM rooms WHERE started_at IS NULL AND finished_at IS NULL")

	if err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, 0)

	for rows.Next() {
		room := &entity.Room{}
		err := rows.Scan(s.getEntityFieldPointers(room)...)

		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	err = rows.Close()

	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *RoomRepository) Store(room *entity.Room) error {
	u := uuid.NewString()
	createdAt := time.Now()

	_, err := s.db.Exec(`
		INSERT INTO rooms (uuid, name, max_players, current_players, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, u, room.Name, room.MaxPlayers, room.CurrentPlayers, room.CreatedBy, createdAt)

	if err != nil {
		return err
	}

	room.UUID = u
	room.CreatedAt = createdAt

	return nil
}

func (s *RoomRepository) FindByUuid(uuid string) (*entity.Room, error) {
	row := s.db.QueryRow("SELECT * FROM rooms WHERE uuid = ? LIMIT 1", uuid)

	room := &entity.Room{}
	err := row.Scan(s.getEntityFieldPointers(room)...)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomRepository) Update(room *entity.Room) error {
	updatedAt := time.Now()

	_, err := s.db.Exec(
		`
		UPDATE rooms
		SET
		    name = ?,
		    max_players = ?,
		    current_players = ?,
		    ready_players = ?,
		    created_by = ?,
		    updated_at = ?,
		    started_at = ?,
		    finished_at = ?
		WHERE uuid = ?
		`,
		room.Name,
		room.MaxPlayers,
		room.CurrentPlayers,
		room.ReadyPlayers,
		room.CreatedBy,
		updatedAt,
		room.StartedAt,
		room.FinishedAt,
		room.UUID,
	)

	if err != nil {
		return err
	}

	room.UpdatedAt = &updatedAt

	return nil
}

func (s *RoomRepository) Delete(room *entity.Room) error {
	_, err := s.db.Exec("DELETE FROM rooms WHERE uuid = ?", room.UUID)

	return err
}

func (s *RoomRepository) getEntityFieldPointers(room *entity.Room) []interface{} {
	return []interface{}{
		&room.UUID,
		&room.Name,
		&room.MaxPlayers,
		&room.CurrentPlayers,
		&room.ReadyPlayers,
		&room.CreatedBy,
		&room.CreatedAt,
		&room.UpdatedAt,
		&room.StartedAt,
		&room.FinishedAt,
	}
}
