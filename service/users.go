package service

import (
	"errors"
	"time"

	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/structs"
	"golang.org/x/crypto/bcrypt"
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
	userInfo, err := s.repository.GetUser(user)
	if err != nil {
		return errors.New("username is not valid")
	}

	hashedPassword := userInfo.HashedPass
	if err := compare(hashedPassword, pass); err != nil {
		return errors.New("incorrect password")
	}
	return nil
}

//Compare string to generated hash
func compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
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
