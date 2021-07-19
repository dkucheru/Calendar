package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dkucheru/Calendar/structs"
	"golang.org/x/crypto/bcrypt"
)

//Generate a salted hash for the input string
func generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

//Compare string to generated hash
func compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

type ArrayRepository struct {
	ArrayRepo []*structs.Event
	ArrayId   int
	// ArrayUsers map[string]*structs.HashedInfo
}

func NewArrayRepository() (*ArrayRepository, error) {
	var events []*structs.Event
	// hashed := make(map[string]*structs.HashedInfo)
	repo := &ArrayRepository{
		ArrayRepo: events,
		ArrayId:   1,
		// ArrayUsers: hashed,
	}
	return repo, nil
}

func (a *ArrayRepository) Get(p structs.EventParams) []*structs.Event {
	var matchedEvents []*structs.Event
	for _, event := range a.ArrayRepo {
		_, weekI := event.Start.ISOWeek()
		if event.Start.Day() == p.Day || p.Day == 0 {
			matchedEvents = append(matchedEvents, event)
		} else if p.Month == 0 || event.Start.Month() == time.Month(p.Month) {
			matchedEvents = append(matchedEvents, event)
		} else if p.Year == 0 || event.Start.Year() == p.Year {
			matchedEvents = append(matchedEvents, event)
		} else if p.Week == 0 || weekI == p.Week {
			matchedEvents = append(matchedEvents, event)
		} else if p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name) {
			matchedEvents = append(matchedEvents, event)
		} else if p.Start == (time.Time{}) || event.Start == p.Start {
			matchedEvents = append(matchedEvents, event)
		} else if p.End == (time.Time{}) || event.End == p.End {
			matchedEvents = append(matchedEvents, event)
		}
	}
	return matchedEvents
}

func (a *ArrayRepository) Update(id int, newEvent structs.Event) (updated structs.Event, err error) {
	foundEvent, err := a.GetByID(id)
	if err != nil {
		return structs.Event{}, err
	}
	foundEvent.Name = newEvent.Name
	foundEvent.Start = newEvent.Start
	foundEvent.End = newEvent.End
	foundEvent.Alert = newEvent.Alert
	foundEvent.Description = newEvent.Description
	return *foundEvent, nil
}

func (a *ArrayRepository) GetByID(id int) (*structs.Event, error) {
	for _, event := range a.ArrayRepo {
		if event.Id == id {
			return event, nil
		}
	}
	message := "event with id [" + fmt.Sprint(id) + "] does not exist"
	return nil, errors.New(message)
}

func (a *ArrayRepository) Add(e structs.Event) structs.Event {
	e.Id = a.ArrayId
	a.ArrayId++
	a.ArrayRepo = append(a.ArrayRepo, &e)
	return e
}

func (a *ArrayRepository) Delete(e structs.Event) {
	for i, event := range a.ArrayRepo {
		if event.Id == e.Id {
			a.ArrayRepo = append(a.ArrayRepo[:i], a.ArrayRepo[i+1:]...)
			return
		}
	}
}

func (a *ArrayRepository) GetLastUsedId() int {
	return a.ArrayId - 1
}

// Implementation of Repository based on map.
type MapRepository struct {
	MapRepo  map[int]structs.Event
	MapId    int
	MapUsers map[string]*structs.HashedInfo
}

func NewMapRepository() (*MapRepository, error) {
	events := make(map[int]structs.Event)
	// hashed := make(map[string]*structs.HashedInfo)
	repo := &MapRepository{
		MapRepo: events,
		MapId:   1,
		// MapUsers: hashed,
	}
	return repo, nil
}

func (m *MapRepository) Get(p structs.EventParams) []*structs.Event {
	var matchedEvents []*structs.Event

	for _, event := range m.MapRepo {
		_, weekI := event.Start.ISOWeek()

		if event.Start.Day() == p.Day || p.Day == 0 {
			matchedEvents = append(matchedEvents, &event)
		} else if p.Month == 0 || event.Start.Month() == time.Month(p.Month) {
			matchedEvents = append(matchedEvents, &event)
		} else if p.Year == 0 || event.Start.Year() == p.Year {
			matchedEvents = append(matchedEvents, &event)
		} else if p.Week == 0 || weekI == p.Week {
			matchedEvents = append(matchedEvents, &event)
		} else if p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name) {
			matchedEvents = append(matchedEvents, &event)
		} else if p.Start == (time.Time{}) || event.Start == p.Start {
			matchedEvents = append(matchedEvents, &event)
		} else if p.End == (time.Time{}) || event.End == p.End {
			matchedEvents = append(matchedEvents, &event)
		}

	}
	return matchedEvents
}

func (m *MapRepository) GetByID(id int) (*structs.Event, error) {
	foundEvent, ok := m.MapRepo[id]
	if !ok {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return nil, errors.New(message)
	}
	return &foundEvent, nil
}

func (m *MapRepository) Update(id int, newEvent structs.Event) (updated structs.Event, err error) {
	foundEvent, ok := m.MapRepo[id]
	if !ok {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return structs.Event{}, errors.New(message)
	}

	foundEvent.Name = newEvent.Name
	foundEvent.Start = newEvent.Start
	foundEvent.End = newEvent.End
	foundEvent.Alert = newEvent.Alert
	foundEvent.Description = newEvent.Description

	m.MapRepo[id] = foundEvent

	return foundEvent, nil
}

func (m *MapRepository) Add(e structs.Event) structs.Event {
	e.Id = m.MapId
	m.MapId++

	m.MapRepo[e.Id] = e
	return e
}

func (m *MapRepository) Delete(e structs.Event) {
	delete(m.MapRepo, e.Id)
}

func (m *MapRepository) GetLastUsedId() int {
	return m.MapId - 1
}

type UsersRepository struct {
	Users map[string]structs.HashedInfo
}

func NewUsersRepository() (*UsersRepository, error) {
	users := make(map[string]structs.HashedInfo)
	repo := &UsersRepository{
		Users: users,
	}
	return repo, nil
}

func (u *UsersRepository) AddUser(e structs.CreateUser) (structs.HashedInfo, error) {
	_, ok := u.Users[e.Username]
	if ok {
		message := "user with username [" + e.Username + "] already exists"
		return structs.HashedInfo{}, errors.New(message)
	}
	generatedHash, err := generate(e.Password)
	if err != nil {
		return structs.HashedInfo{}, err
	}
	loc, err := time.LoadLocation(e.Location)
	if err != nil {
		return structs.HashedInfo{}, err
	}

	hashedData := structs.HashedInfo{
		Username:   e.Username,
		Location:   *loc,
		HashedPass: generatedHash,
	}

	u.Users[e.Username] = hashedData
	return hashedData, nil
}

func (u *UsersRepository) GetUser(username string) (structs.HashedInfo, error) {
	_, ok := u.Users[username]
	if !ok {
		return structs.HashedInfo{}, errors.New("user does not exist")
	}
	return u.Users[username], nil
}

func (u *UsersRepository) UpdateLocation(user string, loc time.Location) (structs.HashedInfo, error) {
	found, ok := u.Users[user]
	if !ok {
		return structs.HashedInfo{}, errors.New("user does not exist")
	}
	found.Location = loc
	u.Users[user] = found
	return u.Users[user], nil
}

func (u *UsersRepository) CheckCredentials(username string, pass string) error {
	userInfo, ok := u.Users[username]
	if !ok {
		return errors.New("username is not valid")
	}

	hashedPassword := userInfo.HashedPass
	if err := compare(hashedPassword, pass); err != nil {
		return errors.New("incorrect password")
	}
	return nil
}
