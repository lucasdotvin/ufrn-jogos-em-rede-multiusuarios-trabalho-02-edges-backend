package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"trabalho-02-edges/config"
)

const driver = "sqlite3"

var cachedDb *sql.DB

func NewDatabase(cfg config.Config) (*sql.DB, error) {
	if cachedDb != nil {
		return cachedDb, nil
	}

	file, err := os.OpenFile(cfg.SQLiteDatabasePath, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	err = file.Close()

	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, cfg.SQLiteDatabasePath)

	if err != nil {
		return nil, err
	}

	err = runMigrations(db)

	if err != nil {
		return nil, err
	}

	cachedDb = db

	return db, nil
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    uuid TEXT PRIMARY KEY NOT NULL,
		    name TEXT NOT NULL,
		    username TEXT NOT NULL,
		    password TEXT NOT NULL,
		    current_score INTEGER DEFAULT 0,
		    created_at DATETIME NOT NULL,
		    updated_at DATETIME DEFAULT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

		CREATE TABLE IF NOT EXISTS rooms (
			uuid TEXT PRIMARY KEY NOT NULL,
			name TEXT NOT NULL,
			max_players INTEGER NOT NULL,
			current_players INTEGER NOT NULL,
			ready_players INTEGER DEFAULT 0,
			created_by TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME DEFAULT NULL,
			started_at DATETIME DEFAULT NULL,
			finished_at DATETIME DEFAULT NULL
		);

		CREATE TABLE IF NOT EXISTS room_user (
			room_uuid TEXT NOT NULL,
			user_uuid TEXT NOT NULL,
		    joined_at DATETIME NOT NULL,
		    won_at DATETIME DEFAULT NULL,
		    lost_at DATETIME DEFAULT NULL,
		    abandoned_at DATETIME DEFAULT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_room_users_room_uuid ON room_user (room_uuid);
		CREATE INDEX IF NOT EXISTS idx_room_users_user_uuid ON room_user (user_uuid);
	`)

	if err != nil {
		return err
	}

	return nil
}
