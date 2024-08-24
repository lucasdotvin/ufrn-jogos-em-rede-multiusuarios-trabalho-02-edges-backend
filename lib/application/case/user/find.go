package user

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/service"
)

type FindUserUseCase struct {
	userService *service.UserService
}

func NewFindUserUseCase(userService *service.UserService) *FindUserUseCase {
	return &FindUserUseCase{
		userService,
	}
}

func (u *FindUserUseCase) Execute(uuid string) (*entity.User, error) {
	return u.userService.Find(uuid)
}
