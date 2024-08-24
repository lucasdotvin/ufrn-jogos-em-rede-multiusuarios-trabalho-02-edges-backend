package service

import (
	"trabalho-02-edges/lib/domain/entity"
	"trabalho-02-edges/lib/domain/repository"
	"trabalho-02-edges/lib/domain/repository/database"
)

type UserService struct {
	hashRepository         repository.HashRepository
	userDatabaseRepository database.UserRepository
}

func NewUserService(hashRepository repository.HashRepository, userDatabaseRepository database.UserRepository) *UserService {
	return &UserService{
		hashRepository,
		userDatabaseRepository,
	}
}

func (s *UserService) SignUp(name string, username string, plainPassword string) (*entity.User, error) {
	isUsernameTaken, err := s.userDatabaseRepository.CheckIfUsernameExists(username)

	if err != nil {
		return nil, err
	}

	if isUsernameTaken {
		return nil, entity.UsernameIsTakenError
	}

	hashedPassword, err := s.hashRepository.Hash(plainPassword)

	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:     name,
		Username: username,
		Password: hashedPassword,
	}

	err = s.userDatabaseRepository.Store(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) SignIn(username string, plainPassword string) (*entity.User, error) {
	user, err := s.userDatabaseRepository.FindByUsername(username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, entity.UserNotFoundError
	}

	isPasswordCorrect := s.hashRepository.Compare(user.Password, plainPassword)

	if !isPasswordCorrect {
		return nil, entity.WrongPasswordError
	}

	return user, nil
}

func (s *UserService) Find(uuid string) (*entity.User, error) {
	u, err := s.userDatabaseRepository.FindByUuid(uuid)

	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, entity.UserNotFoundError
	}

	return u, nil
}

func (s *UserService) GetWhereUuidIn(uuids []string) ([]*entity.User, error) {
	return s.userDatabaseRepository.GetWhereUuidIn(uuids)
}
