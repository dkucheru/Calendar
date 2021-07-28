package service

import (
	"os"
	"testing"
	"time"

	"github.com/dkucheru/Calendar/db"

	"github.com/dkucheru/Calendar/structs"
)

func AddUser(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewUsersDBRepository(repo)
	var testService = newUsersService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}

	testService.AddUser(structs.CreateUser{Username: "testUserExists", Password: "o!", Location: "Local"})
	testCases := map[string]struct {
		user         structs.CreateUser
		errorMessage string
	}{
		"Username already exists": {
			structs.CreateUser{
				Username: "testUserExists",
				Password: "1123",
				Location: "Local",
			},
			"user with username [testUserExists] already exists",
		},
		"Ok new user": {
			structs.CreateUser{
				Username: "newUser",
				Password: "1123",
				Location: "Local",
			},
			"",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			returnedInfo, err := testService.AddUser(test.user)

			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("wrong error : got %q, wanted %q", err, test.errorMessage)
			}

			foundUser, err := testService.repository.GetUser(test.user.Username)
			if err != nil {
				t.Errorf("error finding new user")
			}

			if foundUser.Location.String() != returnedInfo.Location.String() {
				t.Errorf("location adding malfunction")
			}
			if foundUser.Username != returnedInfo.Username {
				t.Errorf("location adding malfunction")
			}
			if foundUser.HashedPass != returnedInfo.HashedPass {
				t.Errorf("location adding malfunction")
			}
		})
	}
}

func UpdateUserTimezoneInDB(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewUsersDBRepository(repo)
	var testService = newUsersService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testService.AddUser(structs.CreateUser{Username: "testUserExists", Password: "o!", Location: "Local"})
	testCases := map[string]struct {
		user         string
		loc          time.Location
		errorMessage string
	}{
		"Username already exists": {
			user:         "Nonexistant",
			loc:          *time.Now().UTC().Location(),
			errorMessage: "user does not exist",
		},
		"Ok change of location": {
			user:         "testUserExists",
			loc:          *time.Now().UTC().Location(),
			errorMessage: "",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			returnedInfo, err := testService.UpdateLocation(test.user, test.loc)

			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("wrong error : got %q, wanted %q", err, test.errorMessage)
			}

			foundUser, err := testService.repository.GetUser(test.user)
			if err != nil {
				t.Errorf("error finding changed user")
			}

			if foundUser.Location.String() != returnedInfo.Location.String() {
				t.Errorf("location change malfunction")
			}
		})

	}

}

func CheckPasswordFromDB(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewUsersDBRepository(repo)
	var testService = newUsersService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testService.AddUser(structs.CreateUser{Username: "testUserExists", Password: "o!", Location: "Local"})
	testCases := map[string]struct {
		user         string
		pass         string
		errorMessage string
	}{
		"Incorrect password": {
			user:         "testUserExists",
			pass:         "0",
			errorMessage: "incorrect password",
		},
		"Correct password": {
			user:         "testUserExists",
			pass:         "o!",
			errorMessage: "",
		},
		"Incorrect login": {
			user:         "badLogin",
			pass:         "o!",
			errorMessage: "username is not valid",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			err := testService.CheckPassword(test.user, test.pass)

			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("wrong error : got %q, wanted %q", err, test.errorMessage)
			}
		})
	}
}

func GetUserLocationFromDB(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewUsersDBRepository(repo)
	var testService = newUsersService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testService.AddUser(structs.CreateUser{Username: "testUserExists", Password: "o!", Location: "Local"})
	testCases := map[string]struct {
		user         string
		errorMessage string
	}{
		"Ok user": {
			user:         "testUserExists",
			errorMessage: "",
		},
		"User does not exist": {
			user:         "badUser",
			errorMessage: "user does not exist",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := testService.GetUserLocation(test.user)

			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("wrong error : got %q, wanted %q", err, test.errorMessage)
			}
		})
	}
}
