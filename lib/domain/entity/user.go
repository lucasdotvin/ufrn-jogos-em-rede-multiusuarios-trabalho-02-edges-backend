package entity

import (
	"errors"
	"time"
)

type User struct {
	UUID         string
	Name         string
	Username     string
	Password     string
	CurrentScore int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

var (
	UsernameIsTakenError = errors.New("username is taken")
	UserNotFoundError    = errors.New("user not found")
	WrongPasswordError   = errors.New("wrong password")
)
