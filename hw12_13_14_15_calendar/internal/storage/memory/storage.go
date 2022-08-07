package memory

import (
	"sync"
	"time"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	bucket map[uuid.UUID]model.Event
}

var _ model.EventsM = (*Storage)(nil)

func NewConnection() (*Storage, error) {
	return &Storage{
		bucket: make(map[uuid.UUID]model.Event),
	}, nil
}

func (s *Storage) Add(event model.Event) (model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event.ID = uuid.New()
	s.bucket[event.ID] = event

	return event, nil
}

func (s *Storage) Update(id uuid.UUID, event model.Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.bucket[event.ID]; !ok {
		return model.ErrEventNotFound
	}

	s.bucket[event.ID] = event

	return nil
}

func (s *Storage) Delete(id uuid.UUID) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.bucket[id]; !ok {
		return model.ErrEventNotFound
	}

	delete(s.bucket, id)

	return nil
}

func (s *Storage) GetListForDay(t time.Time) ([]model.Event, error) {
	result := make([]model.Event, 0)

	y, m, d := t.Date()

	start := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	end := time.Date(y, m, d, 23, 59, 59, 0, t.Location())

	for _, event := range s.bucket {
		if ok := inTimeSegment(event, start, end); ok {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) GetListForWeek(t time.Time) ([]model.Event, error) {
	result := make([]model.Event, 0)

	y, m, d := t.Date()

	start := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 0, 7)

	for _, event := range s.bucket {
		if ok := inTimeSegment(event, start, end); ok {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) GetListForMonth(t time.Time) ([]model.Event, error) {
	result := make([]model.Event, 0)

	y, m, d := t.Date()

	start := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0)

	for _, event := range s.bucket {
		if ok := inTimeSegment(event, start, end); ok {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) GetEvent(id uuid.UUID) (model.Event, error) {
	if _, ok := s.bucket[id]; !ok {
		return model.Event{}, model.ErrEventNotFound
	}

	return s.bucket[id], nil
}

func inTimeSegment(event model.Event, start time.Time, end time.Time) bool {
	isStartPassed := event.StartDate.After(start) || event.StartDate.Equal(start)
	isEndPassed := event.StartDate.Before(end) || event.StartDate.Equal(end)

	return isStartPassed && isEndPassed
}
