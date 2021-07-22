package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/dkucheru/Calendar/structs"
	"golang.org/x/crypto/bcrypt"
)

func Initialize(username, password, database string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)
	fmt.Println(dsn)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Database connection established")
	return conn, nil
}

const (
	HOST = "database"
	PORT = 5432
)

type UsersDBRepository struct {
	Conn *sql.DB
}

func NewUsersDBRepository(conn *sql.DB) (*UsersDBRepository, error) {
	return &UsersDBRepository{Conn: conn}, nil
}

func (db *UsersDBRepository) AddUser(e structs.CreateUser) (structs.HashedInfo, error) {
	generatedHash, err := generate(e.Password)
	if err != nil {
		return structs.HashedInfo{}, err
	}
	loc, err := time.LoadLocation(e.Location)
	if err != nil {
		return structs.HashedInfo{}, err
	}
	query := `INSERT INTO users (username, hashedpass, userlocation) VALUES ($1, $2, $3);`
	_, err = db.Conn.Query(query, e.Username, generatedHash, e.Location)
	if err != nil {
		fmt.Println("error query insert user")
		fmt.Println(err.Error())
		return structs.HashedInfo{}, err
	}

	return structs.HashedInfo{
		Username:   e.Username,
		Location:   *loc,
		HashedPass: generatedHash,
	}, nil
}

func (db *UsersDBRepository) GetUser(user string) (structs.HashedInfo, error) {
	var item structs.CreateUser
	justAdded :=
		`SELECT username,hashedpass,userlocation
	FROM users
	WHERE username = $1;`
	err := db.Conn.QueryRow(justAdded, user).Scan(&item.Username, &item.Password, &item.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.HashedInfo{}, ErrNoMatch
		}
		return structs.HashedInfo{}, err
	}
	if err != nil {
		fmt.Println("issue with scan")
		fmt.Println(err.Error())
		return structs.HashedInfo{}, err
	}
	loc, err := time.LoadLocation(item.Location)
	if err != nil {
		return structs.HashedInfo{}, err
	}
	return structs.HashedInfo{
		Username:   item.Username,
		Location:   *loc,
		HashedPass: item.Password,
	}, nil
}

func (db *UsersDBRepository) UpdateLocation(user string, loc time.Location) (structs.HashedInfo, error) {
	var event structs.CreateUser
	query := `UPDATE events 
	SET userlocation = $1
	 WHERE username=$2 RETURNING username,hashedpass,userlocation;`
	err := db.Conn.QueryRow(query, loc.String(), user).
		Scan(&event.Username, &event.Password, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.HashedInfo{}, ErrNoMatch
		}
		return structs.HashedInfo{}, err
	}
	newloc, err := time.LoadLocation(event.Location)
	if err != nil {
		return structs.HashedInfo{}, err
	}

	return structs.HashedInfo{
		Username:   event.Username,
		Location:   *newloc,
		HashedPass: event.Password,
	}, nil
}

type EventsDBRepository struct {
	Conn *sql.DB
}

var ErrNoMatch = fmt.Errorf("no matching record")

func NewDatabaseRepository(conn *sql.DB) (*EventsDBRepository, error) {
	return &EventsDBRepository{Conn: conn}, nil
}

func (db *EventsDBRepository) Add(e structs.Event) (structs.Event, error) {
	var id int
	query := `INSERT INTO events (event_name,event_start,event_end,event_description, event_alert) 
	VALUES ($1, $2,$3,$4,$5) RETURNING eventid`
	err := db.Conn.QueryRow(query, e.Name, e.Start, e.End, e.Description, e.Alert).Scan(&id)
	if err != nil {
		fmt.Println("error with query add event")
		return structs.Event{}, err
	}
	justAdded, err := db.GetByID(id)
	if err != nil {
		fmt.Println("error finding event")
		return structs.Event{}, err
	}
	return justAdded, nil
}

func (db *EventsDBRepository) Get(p structs.EventParams) ([]structs.Event, error) {
	var list []structs.Event
	query :=
		`SELECT eventid,event_name,event_start,event_end,event_description, event_alert
	FROM events
	WHERE event_name = $1 OR date_part('day', timestamp event_start) = $2 OR 
	date_part('week', timestamp event_start) = $3 OR date_part('month', timestamp event_start) = $4 
	OR date_part('year', timestamp event_start) = $5 OR event_start = $6 OR event_end = $7
	IF $8 THEN
		ORDER BY ID DESC
	END IF;`
	rows, err := db.Conn.Query(query, p.Name, p.Day, p.Week, p.Month, p.Year, p.Start, p.End, p.Sorting)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var item structs.Event
		err := rows.Scan(&item.Id, &item.Name, &item.Start, &item.End, &item.Description, &item.Alert)
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (db *EventsDBRepository) GetByID(id int) (structs.Event, error) {
	var item structs.Event
	justAdded :=
		`SELECT eventid,event_name,event_start,event_end,event_description, event_alert
	FROM events
	WHERE eventid = $1;`
	err := db.Conn.QueryRow(justAdded, id).Scan(&item.Id, &item.Name, &item.Start, &item.End, &item.Description, &item.Alert)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.Event{}, ErrNoMatch
		}
		return structs.Event{}, err
	}
	return item, nil
}

func (db *EventsDBRepository) Update(id int, e structs.Event) (updated structs.Event, err error) {
	var event structs.Event
	query := `UPDATE events 
	SET event_name = $1, event_start = $2, event_end = $3, event_description = &4, event_alert = $5
	 WHERE eventid=$6 RETURNING eventid,event_name,event_start,event_end,event_description, event_alert;`
	err = db.Conn.QueryRow(query, e.Name, e.Start, e.End, e.Description, e.Alert, id).
		Scan(&event.Id, &event.Name, &event.Start, &event.End, &event.Description, &event.Alert)
	if err != nil {
		if err == sql.ErrNoRows {
			return event, ErrNoMatch
		}
		return event, err
	}
	return event, nil
}

func (db *EventsDBRepository) Delete(e structs.Event) error {
	query := `DELETE FROM events WHERE eventid = $1;`
	_, err := db.Conn.Exec(query, e.Id)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db *EventsDBRepository) GetLastUsedId() int {
	var lastUsedId int
	query :=
		`SELECT eventid
	FROM events
	WHERE eventid >= ALL(
		SELECT eventid
		FROM events
	);`
	err := db.Conn.QueryRow(query).Scan(&lastUsedId)
	if err != nil {
		return 0
	}
	return lastUsedId
}

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

type ArrayRepository struct {
	ArrayRepo []*structs.Event
	ArrayId   int
}

func NewArrayRepository() (*ArrayRepository, error) {
	var events []*structs.Event
	repo := &ArrayRepository{
		ArrayRepo: events,
		ArrayId:   1,
	}
	return repo, nil
}

func (a *ArrayRepository) Get(p structs.EventParams) ([]structs.Event, error) {
	var matchedEvents []structs.Event
	for _, event := range a.ArrayRepo {
		_, weekI := event.Start.ISOWeek()
		if event.Start.Day() == p.Day || p.Day == 0 {
			matchedEvents = append(matchedEvents, *event)
		} else if p.Month == 0 || event.Start.Month() == time.Month(p.Month) {
			matchedEvents = append(matchedEvents, *event)
		} else if p.Year == 0 || event.Start.Year() == p.Year {
			matchedEvents = append(matchedEvents, *event)
		} else if p.Week == 0 || weekI == p.Week {
			matchedEvents = append(matchedEvents, *event)
		} else if p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name) {
			matchedEvents = append(matchedEvents, *event)
		} else if p.Start == (time.Time{}) || event.Start == p.Start {
			matchedEvents = append(matchedEvents, *event)
		} else if p.End == (time.Time{}) || event.End == p.End {
			matchedEvents = append(matchedEvents, *event)
		}
	}
	return matchedEvents, nil
}

func (a *ArrayRepository) Update(id int, newEvent structs.Event) (updated structs.Event, err error) {
	var foundEvent *structs.Event
	for _, event := range a.ArrayRepo {
		if event.Id == id {
			foundEvent = event
		}
	}
	if foundEvent == nil {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return structs.Event{}, errors.New(message)
	}
	// foundEvent, err := a.GetByID(id)
	// if err != nil {
	// 	return structs.Event{}, err
	// }
	foundEvent.Name = newEvent.Name
	foundEvent.Start = newEvent.Start
	foundEvent.End = newEvent.End
	foundEvent.Alert = newEvent.Alert
	foundEvent.Description = newEvent.Description
	return *foundEvent, nil
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

func (a *ArrayRepository) Add(e structs.Event) (structs.Event, error) {
	e.Id = a.ArrayId
	a.ArrayId++
	a.ArrayRepo = append(a.ArrayRepo, &e)
	return e, nil
}

func (a *ArrayRepository) Delete(e structs.Event) error {
	for i, event := range a.ArrayRepo {
		if event.Id == e.Id {
			a.ArrayRepo = append(a.ArrayRepo[:i], a.ArrayRepo[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
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

func (m *MapRepository) Get(p structs.EventParams) ([]structs.Event, error) {
	var matchedEvents []structs.Event

	for _, event := range m.MapRepo {
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
	return matchedEvents, nil
}

func (m *MapRepository) GetByID(id int) (structs.Event, error) {
	foundEvent, ok := m.MapRepo[id]
	if !ok {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return structs.Event{}, errors.New(message)
	}
	return foundEvent, nil
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

func (m *MapRepository) Add(e structs.Event) (structs.Event, error) {
	e.Id = m.MapId
	m.MapId++

	m.MapRepo[e.Id] = e
	return e, nil
}

func (m *MapRepository) Delete(e structs.Event) error {
	delete(m.MapRepo, e.Id)
	return nil
}

func (m *MapRepository) GetLastUsedId() int {
	return m.MapId - 1
}

type UsersRepository struct {
	Users map[string]structs.HashedInfo
}

func NewUsersInMemoryRepository() (*UsersRepository, error) {
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
