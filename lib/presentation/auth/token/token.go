package token

import (
	"errors"
	"time"
	"trabalho-02-edges/lib/domain/entity"
)

const (
	AccessTokenKey    = "access_token"
	RefreshTokenKey   = "refresh_token"
	BroadcastTokenKey = "broadcast_token"
)

var (
	InvalidTokenError     = errors.New("invalid token")
	InvalidTokenGoalError = errors.New("invalid token goal")
)

type Type string

const (
	Bearer Type = "Bearer"
)

type Goal string

const (
	Access    Goal = "acc"
	Refresh   Goal = "ref"
	Broadcast Goal = "bro"
)

func ParseGoal(s string) (Goal, error) {
	switch s {
	case "acc":
		return Access, nil
	case "ref":
		return Refresh, nil
	case "bro":
		return Broadcast, nil
	default:
		return "", InvalidTokenGoalError
	}
}

type Token struct {
	Content   string
	Uid       string
	Type      Type
	Goal      Goal
	ExpiresAt time.Time
	IssuedAt  time.Time
}

func (t *Token) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

func (t *Token) IsValid() bool {
	return !t.IsExpired()
}

func (t *Token) IsAccessToken() bool {
	return t.Goal == Access
}

func (t *Token) IsRefreshToken() bool {
	return t.Goal == Refresh
}

func (t *Token) IsBroadcastToken() bool {
	return t.Goal == Broadcast
}

type Service interface {
	GenerateAccessToken(user *entity.User) (*Token, error)

	GenerateRefreshToken(user *entity.User) (*Token, error)

	GenerateBroadcastToken(user *entity.User) (*Token, error)

	ParseTokenFromContent(content string) (*Token, error)
}
