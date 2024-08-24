package sqlite

import (
	"database/sql"
	"github.com/google/uuid"
	"strings"
	"time"
	"trabalho-02-edges/lib/domain/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewSqliteUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (s *UserRepository) CheckIfUsernameExists(username string) (bool, error) {
	row := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? LIMIT 1", username)

	var count int
	err := row.Scan(&count)

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (s *UserRepository) FindByUsername(username string) (*entity.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE username = ? LIMIT 1", username)

	user := &entity.User{}
	err := row.Scan(s.getEntityFieldPointers(user)...)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserRepository) FindByUuid(uuid string) (*entity.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE uuid = ? LIMIT 1", uuid)

	user := &entity.User{}
	err := row.Scan(s.getEntityFieldPointers(user)...)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserRepository) Store(user *entity.User) error {
	u := uuid.NewString()
	createdAt := time.Now()

	_, err := s.db.Exec(`
		INSERT INTO users (uuid, name, username, password, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, u, user.Name, user.Username, user.Password, createdAt)

	if err != nil {
		return err
	}

	user.UUID = u
	user.CreatedAt = createdAt

	return nil
}

func (s *UserRepository) GetWhereUuidIn(uuids []string) ([]*entity.User, error) {
	questionMarks := strings.Repeat("?, ", len(uuids))
	questionMarks = strings.TrimSuffix(questionMarks, ", ")

	queryParams := make([]any, 0, len(uuids))

	for _, u := range uuids {
		queryParams = append(queryParams, u)
	}

	query := "SELECT * FROM users WHERE uuid IN (" + questionMarks + ")"

	rows, err := s.db.Query(query, queryParams...)

	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0)

	for rows.Next() {
		user := &entity.User{}
		err = rows.Scan(s.getEntityFieldPointers(user)...)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserRepository) getEntityFieldPointers(user *entity.User) []interface{} {
	return []interface{}{
		&user.UUID,
		&user.Name,
		&user.Username,
		&user.Password,
		&user.CurrentScore,
		&user.CreatedAt,
		&user.UpdatedAt,
	}
}
