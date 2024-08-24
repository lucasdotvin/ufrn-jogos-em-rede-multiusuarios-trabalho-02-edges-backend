package user

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/service"
)

type SignUpUseCase struct {
	userService *service.UserService
}

func NewSignUpUseCase(userService *service.UserService) *SignUpUseCase {
	return &SignUpUseCase{
		userService,
	}
}

func (u *SignUpUseCase) Execute(name string, username string, plainPassword string) (*entity.User, error) {
	return u.userService.SignUp(name, username, plainPassword)
}
