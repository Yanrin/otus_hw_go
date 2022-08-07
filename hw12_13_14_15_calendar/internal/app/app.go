package app

import (
	"time"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/google/uuid"
)

type Calendar struct {
	Stg model.EventsM
}

func NewCalendar(em model.EventsM) Calendar {
	return Calendar{
		Stg: em,
	}
}

func (c *Calendar) Add(event model.Event) (model.Event, error) {
	return c.Stg.Add(event)
}

func (c *Calendar) Update(id uuid.UUID, event model.Event) error {
	return c.Stg.Update(id, event)
}

func (c *Calendar) Delete(id uuid.UUID) error {
	return c.Stg.Delete(id)
}

func (c *Calendar) GetEvent(id uuid.UUID) (model.Event, error) {
	return c.Stg.GetEvent(id)
}

func (c *Calendar) GetListForDay(date time.Time) ([]model.Event, error) {
	return c.Stg.GetListForDay(date)
}

func (c *Calendar) GetListForWeek(date time.Time) ([]model.Event, error) {
	return c.Stg.GetListForDay(date)
}

func (c *Calendar) GetListForMonth(date time.Time) ([]model.Event, error) {
	return c.Stg.GetListForDay(date)
}
