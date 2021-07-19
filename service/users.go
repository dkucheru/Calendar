package service

import (
	"time"

	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/structs"
)

type usersService struct {
	repository db.UserRepository
}

func newUsersService(repository db.UserRepository) *usersService {
	s := usersService{
		repository: repository,
	}
	return &s
}

func (s *usersService) AddUser(newUser structs.CreateUser) (structs.HashedInfo, error) {
	return s.repository.AddUser(newUser)
}

func (s *usersService) CheckPassword(user string, pass string) error {
	return s.repository.CheckCredentials(user, pass)
}

func (s *usersService) GetUserLocation(username string) (time.Location, error) {
	userInfo, err := s.repository.GetUser(username)
	if err != nil {
		return time.Location{}, err
	}
	return userInfo.Location, nil
}

func (s *usersService) UpdateLocation(user string, newLocation time.Location) (structs.HashedInfo, error) {
	return s.repository.UpdateLocation(user, newLocation)
}
