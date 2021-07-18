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
	ArrayRepo  []*structs.Event
	ArrayId    int
	ArrayUsers map[string]*structs.HashedInfo
}

func NewArrayRepository() (*ArrayRepository, error) {
	var events []*structs.Event
	hashed := make(map[string]*structs.HashedInfo)
	repo := &ArrayRepository{
		ArrayRepo:  events,
		ArrayId:    1,
		ArrayUsers: hashed,
	}
	return repo, nil
}

func (a *ArrayRepository) CheckCredentials(username string, pass string) error {
	userInfo, ok := a.ArrayUsers[username]
	if !ok {
		return errors.New("Username is not valid")
	}

	hashedPassword := userInfo.HashedPass

	return compare(hashedPassword, pass)
}

func (a *ArrayRepository) GetAll() []*structs.Event {
	return a.ArrayRepo
}

func (a *ArrayRepository) GetByID(id int) (structs.Event, error) {
	for _, event := range a.ArrayRepo {
		if event.Id == id {
			return *event, nil
		}
	}
	message := "event with id [" + fmt.Sprint(id) + "] does not exist"
	return structs.Event{}, errors.New(message)
}

func (a *ArrayRepository) AddEvent(e structs.Event) error {
	for _, event := range a.ArrayRepo {
		if event.Id == e.Id {
			message := "event with id [" + fmt.Sprint(e.Id) + "] already exists"
			return errors.New(message)
		}
	}
	a.ArrayRepo = append(a.ArrayRepo, &e)
	return nil
}

func (a *ArrayRepository) AddUser(e structs.CreateUser) error {
	_, ok := a.ArrayUsers[e.Username]
	if ok {
		message := "user with username [" + fmt.Sprint(e.Username) + "] already exists"
		return errors.New(message)
	}
	generatedHash, err := generate(e.Password)
	if err != nil {
		panic("Implement me")
	}

	loc, err := time.LoadLocation(e.Location)
	if err != nil {
		panic("Implement me")
	}

	hashedData := structs.HashedInfo{
		Username:   e.Username,
		Location:   *loc,
		HashedPass: generatedHash,
	}

	a.ArrayUsers[e.Username] = &hashedData
	return nil
}

func (a *ArrayRepository) Delete(e structs.Event) {
	for i, event := range a.ArrayRepo {
		if event.Id == e.Id {
			a.ArrayRepo = append(a.ArrayRepo[:i], a.ArrayRepo[i+1:]...)
			return
		}
	}
}

func (a *ArrayRepository) GetUser(username string) (*structs.HashedInfo, error) {
	_, ok := a.ArrayUsers[username]
	if !ok {
		return nil, errors.New("user does not exist")
	}
	return a.ArrayUsers[username], nil
}

func (a *ArrayRepository) GetNextId() int {
	a.ArrayId++
	return a.ArrayId - 1
}

func (a *ArrayRepository) GetLastUsedId() int {
	return a.ArrayId - 1
}

func (a *ArrayRepository) MatchParams(event structs.Event, p structs.EventParams) bool {
	_, weekI := event.Start.ISOWeek()

	if !(event.Start.Day() == p.Day || p.Day == 0) {
		return false
	}
	if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
		return false
	}
	if !(p.Year == 0 || event.Start.Year() == p.Year) {
		return false
	}
	if !(p.Week == 0 || weekI == p.Week) {
		return false
	}
	if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
		return false
	}
	if !(p.Start == (time.Time{}) || event.Start == p.Start) {
		return false
	}
	if !(p.End == (time.Time{}) || event.End == p.End) {
		return false
	}
	return true
}

// Implementation of Repository based on map.
type MapRepository struct {
	MapRepo  map[int]*structs.Event
	MapId    int
	MapUsers map[string]*structs.HashedInfo
}

func NewMapRepository() (*MapRepository, error) {
	events := make(map[int]*structs.Event)
	hashed := make(map[string]*structs.HashedInfo)
	repo := &MapRepository{
		MapRepo:  events,
		MapId:    1,
		MapUsers: hashed,
	}
	return repo, nil
}

func (m *MapRepository) CheckCredentials(username string, pass string) error {
	userInfo, ok := m.MapUsers[username]
	if !ok {
		return errors.New("Username is not valid")
	}

	hashedPassword := userInfo.HashedPass

	return compare(hashedPassword, pass)
}

func (m *MapRepository) GetAll() []*structs.Event {
	events := make([]*structs.Event, 0, len(m.MapRepo))
	for _, v := range m.MapRepo {
		events = append(events, v)
	}
	return events
}

func (m *MapRepository) GetUser(username string) (*structs.HashedInfo, error) {
	_, ok := m.MapUsers[username]
	if !ok {
		return nil, errors.New("user does not exist")
	}
	return m.MapUsers[username], nil
}

func (m *MapRepository) GetByID(id int) (structs.Event, error) {
	_, ok := m.MapRepo[id]
	if !ok {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return structs.Event{}, errors.New(message)
	}
	return *m.MapRepo[id], nil
}

func (m *MapRepository) AddEvent(e structs.Event) error {
	_, ok := m.MapRepo[e.Id]
	if ok {
		message := "event with id [" + fmt.Sprint(e.Id) + "] already exists"
		return errors.New(message)
	}
	m.MapRepo[e.Id] = &e
	return nil
}

func (m *MapRepository) AddUser(e structs.CreateUser) error {
	_, ok := m.MapUsers[e.Username]
	if ok {
		message := "user with username [" + fmt.Sprint(e.Username) + "] already exists"
		return errors.New(message)
	}
	generatedHash, err := generate(e.Password)
	if err != nil {
		panic("Implement me")
	}
	loc, err := time.LoadLocation(e.Location)
	if err != nil {
		panic("Implement me")
	}

	hashedData := structs.HashedInfo{
		Username:   e.Username,
		Location:   *loc,
		HashedPass: generatedHash,
	}

	m.MapUsers[e.Username] = &hashedData
	return nil
}

func (m *MapRepository) Delete(e structs.Event) {
	delete(m.MapRepo, e.Id)
}

func (m *MapRepository) GetNextId() int {
	m.MapId++
	return m.MapId - 1
}

func (m *MapRepository) GetLastUsedId() int {
	return m.MapId - 1
}

func (m *MapRepository) MatchParams(event structs.Event, p structs.EventParams) bool {
	_, weekI := event.Start.ISOWeek()

	if !(event.Start.Day() == p.Day || p.Day == 0) {
		return false
	}
	if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
		return false
	}
	if !(p.Year == 0 || event.Start.Year() == p.Year) {
		return false
	}
	if !(p.Week == 0 || weekI == p.Week) {
		return false
	}
	if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
		return false
	}
	if !(p.Start == (time.Time{}) || event.Start == p.Start) {
		return false
	}
	if !(p.End == (time.Time{}) || event.End == p.End) {
		return false
	}
	return true
}
