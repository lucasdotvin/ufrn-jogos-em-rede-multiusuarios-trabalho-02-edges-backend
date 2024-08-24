package user

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/service"
)

type SignInUseCase struct {
	userService *service.UserService
}

func NewSignInUseCase(userService *service.UserService) *SignInUseCase {
	return &SignInUseCase{
		userService,
	}
}

func (u *SignInUseCase) Execute(username string, password string) (*entity.User, error) {
	return u.userService.SignIn(username, password)
}
