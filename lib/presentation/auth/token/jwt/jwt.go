package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"trabalho-02-edges/config"
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/presentation/auth/token"
)

type Service struct {
	key      []byte
	duration time.Duration
	renewDue time.Duration
}

func NewJwtService(cfg config.Config) *Service {
	return &Service{
		key:      []byte(cfg.JwtSecret),
		duration: time.Duration(cfg.JwtDurationInMinutes) * time.Minute,
		renewDue: time.Duration(cfg.JwtRenewDueInMinutes) * time.Minute,
	}
}

func (s *Service) GenerateAccessToken(user *entity.User) (*token.Token, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(s.duration)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.UUID,
		"exp": expiresAt.Unix(),
		"iat": issuedAt.Unix(),
		"goa": token.Access,
	})

	signed, err := t.SignedString(s.key)

	if err != nil {
		return nil, err
	}

	return &token.Token{
		Content:   signed,
		Uid:       user.UUID,
		Type:      token.Bearer,
		Goal:      token.Access,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}, nil
}

func (s *Service) GenerateRefreshToken(user *entity.User) (*token.Token, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(s.renewDue)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.UUID,
		"exp": expiresAt.Unix(),
		"iat": issuedAt.Unix(),
		"goa": token.Refresh,
	})

	signed, err := t.SignedString(s.key)

	if err != nil {
		return nil, err
	}

	return &token.Token{
		Content:   signed,
		Uid:       user.UUID,
		Type:      token.Bearer,
		Goal:      token.Refresh,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}, nil
}

func (s *Service) GenerateBroadcastToken(user *entity.User) (*token.Token, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(s.duration)

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.UUID,
		"exp": expiresAt.Unix(),
		"iat": issuedAt.Unix(),
		"goa": token.Broadcast,
	})

	signed, err := t.SignedString(s.key)

	if err != nil {
		return nil, err
	}

	return &token.Token{
		Content:   signed,
		Uid:       user.UUID,
		Type:      token.Bearer,
		Goal:      token.Broadcast,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}, nil
}

func (s *Service) ParseTokenFromContent(content string) (*token.Token, error) {
	rawToken, err := jwt.Parse(content, func(t *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := rawToken.Claims.(jwt.MapClaims)

	if !ok || !rawToken.Valid {
		return nil, token.InvalidTokenError
	}

	uid, ok := claims["uid"].(string)

	if !ok {
		return nil, token.InvalidTokenError
	}

	exp, ok := claims["exp"].(float64)

	if !ok {
		return nil, token.InvalidTokenError
	}

	iat, ok := claims["iat"].(float64)

	if !ok {
		return nil, token.InvalidTokenError
	}

	goa, ok := claims["goa"].(string)

	if !ok {
		return nil, token.InvalidTokenError
	}

	goal, err := token.ParseGoal(goa)

	if err != nil {
		return nil, err
	}

	return &token.Token{
		Content:   content,
		Uid:       uid,
		Type:      token.Bearer,
		Goal:      goal,
		ExpiresAt: time.Unix(int64(exp), 0),
		IssuedAt:  time.Unix(int64(iat), 0),
	}, nil
}
