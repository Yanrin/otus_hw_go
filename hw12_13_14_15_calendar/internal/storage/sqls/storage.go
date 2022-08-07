package sqls

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB
}

var _ model.EventsM = (*Storage)(nil)

func NewConnection(cfg *config.Config) (*Storage, error) {
	var s Storage
	var err error
	s.db, err = sql.Open(
		cfg.Database.Driver,
		fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
		),
	)
	if err != nil {
		return nil, err
	}
	err = s.db.Ping()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Storage) Add(event model.Event) (model.Event, error) {
	query := "INSERT INTO events (id, title, start_date, end_date, description, user_id, notify_at) VALUES ($1, $2, $3, $4, $5, $6, $7);" //nolint:lll
	event.ID = uuid.New()

	_, err := s.db.Exec(
		query,
		event.ID,
		event.Title,
		event.StartDate,
		event.EndDate,
		event.Description,
		event.UserID,
		event.NotifyAt,
	)

	return event, err
}

func (s *Storage) Update(id uuid.UUID, event model.Event) error {
	_, err := s.db.Exec(
		"UPDATE events SET title=$1, start_date=$2, end_date=$3, description=$4, user_id=$5, notify_at=$6 WHERE id = $7;", //nolint:lll
		event.Title,
		event.StartDate,
		event.EndDate,
		event.Description,
		event.UserID,
		event.NotifyAt,
		event.ID,
	)
	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	return err
}

func (s *Storage) GetListForDay(t time.Time) ([]model.Event, error) {
	y, m, d := t.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	end := time.Date(y, m, d, 23, 59, 59, 0, t.Location())

	return s.getListPeriod(start, end)
}

func (s *Storage) GetListForWeek(t time.Time) ([]model.Event, error) {
	y, m, d := t.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 0, 7)

	return s.getListPeriod(start, end)
}

func (s *Storage) GetListForMonth(t time.Time) ([]model.Event, error) {
	y, m, d := t.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0)

	return s.getListPeriod(start, end)
}

func (s *Storage) GetEvent(id uuid.UUID) (model.Event, error) {
	var event model.Event
	row := s.db.QueryRow("SELECT * FROM events WHERE id = $1", id)

	err := row.Scan(&event.ID, &event.Title, &event.StartDate, &event.EndDate, &event.Description, &event.UserID, &event.NotifyAt) //nolint:lll
	if err == sql.ErrNoRows {
		err = model.ErrEventNotFound
	}

	return event, err
}

func (s *Storage) getListPeriod(start, end time.Time) ([]model.Event, error) {
	result := make([]model.Event, 0)
	rows, err := s.db.Query("SELECT * FROM events WHERE start_date >= $1 AND start_date <= $2", start, end)

	if errors.Is(err, sql.ErrNoRows) {
		return result, model.ErrEventNotFound
	}
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var event model.Event
		err := rows.Scan(&event.ID, &event.Title, &event.StartDate, &event.EndDate, &event.Description, &event.UserID, &event.NotifyAt) //nolint:lll
		if err != nil {
			return result, err
		}
		result = append(result, event)
	}

	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}
