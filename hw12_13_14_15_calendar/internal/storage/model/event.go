package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	Description string    `db:"description"`
	UserID      uuid.UUID `db:"user_id"`
	NotifyAt    time.Time `db:"notify_at"`
}

type EventsM interface {
	Add(event Event) (Event, error)
	Update(id uuid.UUID, event Event) error
	Delete(id uuid.UUID) error
	GetListForDay(date time.Time) ([]Event, error)
	GetListForWeek(date time.Time) ([]Event, error)
	GetListForMonth(date time.Time) ([]Event, error)
	GetEvent(id uuid.UUID) (Event, error)
}

type Storage struct{}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("time is busy")
)
