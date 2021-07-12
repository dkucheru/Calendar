package db

import "github.com/dkucheru/Calendar/structs"

// FIXME: Rename 
type DatabaseRepository struct {
	DBrepository []*structs.Event
}

func NewDBRepository() (*DatabaseRepository, error) {
	var events []*structs.Event
	repo := &DatabaseRepository{
		DBrepository: events,
	}
	return repo, nil
}

func (d *DatabaseRepository) Array() []*structs.Event {
	return d.DBrepository
}

func (d *DatabaseRepository) AddToArray(s *structs.Event) {
	d.DBrepository = append(d.DBrepository, s)
}

func (d *DatabaseRepository) RemoveFromArray(i int) {
	d.DBrepository = append(d.DBrepository[:i], d.DBrepository[i+1:]...)
}
