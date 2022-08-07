package main

import (
	"testing"
	"time"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/model"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/sqls"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const CHECKSQL = false

var (
	cfgPath  = "testdata/config/config.yml"
	userID1  = uuid.New()
	userID2  = uuid.New()
	testdata = []model.Event{
		{
			ID:          uuid.New(),
			Title:       "Event 1",
			StartDate:   time.Now().Round(time.Microsecond),
			EndDate:     time.Now().Round(time.Microsecond).Add(time.Hour),
			Description: "First event",
			UserID:      userID1,
			NotifyAt:    time.Now().Round(time.Microsecond).Add(-30 * time.Minute),
		},
		{
			ID:          uuid.New(),
			Title:       "Event 2",
			StartDate:   time.Now().AddDate(0, 0, 4),
			EndDate:     time.Now().AddDate(0, 0, 5),
			Description: "Previous",
			UserID:      userID1,
			NotifyAt:    time.Now().AddDate(0, 0, 3),
		},
		{
			ID:          uuid.New(),
			Title:       "Test Event",
			StartDate:   time.Now().AddDate(0, 0, 15),
			EndDate:     time.Now().AddDate(0, 0, 16),
			Description: "Far far far after",
			UserID:      userID2,
			NotifyAt:    time.Now().AddDate(0, 0, 12),
		},
	}
)

type ConnType struct {
	mode    string
	storage model.EventsM
}

func TestConnection(t *testing.T) {
	t.Run("memory", func(t *testing.T) {
		_, err := memory.NewConnection()
		require.NoError(t, err)
	})
	if CHECKSQL {
		t.Run("sqls", func(t *testing.T) {
			cfg, err := config.New(cfgPath)
			require.NoError(t, err)

			_, err = sqls.NewConnection(cfg)
			require.NoError(t, err)
		})
	}
}

func TestStorage(t *testing.T) {
	connections := []ConnType{}
	var connect model.EventsM

	connect, _ = memory.NewConnection()
	connections = append(connections, ConnType{mode: "memory", storage: connect})

	if CHECKSQL {
		cfg, _ := config.New(cfgPath)
		connect, _ = sqls.NewConnection(cfg)
		connections = append(connections, ConnType{mode: "sqls", storage: connect})
	}

	for _, tc := range connections {
		t.Run(tc.mode+"Single add/get/delete", func(t *testing.T) {
			event, err := tc.storage.Add(testdata[0])

			require.NoError(t, err)
			require.True(t, compareAddedEvents(event, testdata[0]))

			eventGot, err := tc.storage.GetEvent(event.ID)
			require.NoError(t, err)
			require.True(t, compareAddedEvents(event, eventGot))

			err = tc.storage.Delete(event.ID)
			require.NoError(t, err)

			_, err = tc.storage.GetEvent(event.ID)
			require.ErrorIs(t, err, model.ErrEventNotFound)
		})

		t.Run(tc.mode+"GetEvent empty result", func(t *testing.T) {
			eventAdded, err := tc.storage.Add(testdata[0])
			require.NoError(t, err)
			require.True(t, compareAddedEvents(eventAdded, testdata[0]))

			defer func() {
				err = tc.storage.Delete(eventAdded.ID)
				require.NoError(t, err)
			}()

			var diffID uuid.UUID
			for {
				diffID = uuid.New()
				if diffID != eventAdded.ID {
					break
				}
			}

			_, err = tc.storage.GetEvent(diffID)
			require.ErrorIs(t, err, model.ErrEventNotFound)
		})

		t.Run(tc.mode+"GetList", func(t *testing.T) {
			var addedList []model.Event
			for _, td := range testdata {
				addedEvent, err := tc.storage.Add(td)
				addedList = append(addedList, addedEvent)
				require.NoError(t, err)
			}
			defer func() {
				for _, td := range addedList {
					err := tc.storage.Delete(td.ID)
					require.NoError(t, err)
				}
			}()

			eventList, err := tc.storage.GetListForDay(time.Now())
			require.NoError(t, err)
			require.Equal(t, 1, len(eventList))
			require.True(t, compareAddedEvents(eventList[0], testdata[0]))

			eventList, err = tc.storage.GetListForWeek(time.Now())
			require.NoError(t, err)
			require.Equal(t, 2, len(eventList))

			eventList, err = tc.storage.GetListForMonth(time.Now())
			require.NoError(t, err)
			require.Equal(t, 3, len(eventList))
		})
	}
}

// compareAddedEvents compares all fields of Events except ID.
func compareAddedEvents(e1, e2 model.Event) bool {
	if e1.Title != e2.Title {
		return false
	}
	if e1.StartDate.String() != e2.StartDate.String() { // because of time.Local and time.Location() magic in time.Date()
		return false
	}
	if e1.EndDate.String() != e2.EndDate.String() {
		return false
	}
	if e1.Title != e2.Title {
		return false
	}
	if e1.Description != e2.Description {
		return false
	}
	if e1.NotifyAt.String() != e2.NotifyAt.String() {
		return false
	}
	return true
}
