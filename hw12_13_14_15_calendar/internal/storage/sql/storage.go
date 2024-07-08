package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // for postgres
	"github.com/jmoiron/sqlx"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db         *sqlx.DB
	driverName string
	dsn        string
}

func New(driverName string, dsn string) *Storage {
	return &Storage{
		driverName: driverName,
		dsn:        dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, s.driverName, s.dsn)
	return
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

const checkExistingEventsSQL = `
SELECT count(*) FROM events WHERE start_date <= ? AND end_date >= ?
`

const createEventSQL = `
INSERT INTO events (title, start_date, end_date, description, user_id, notify_before)
VALUES (:title, :start_date, :end_date, :description, :user_id, :notify_before)
RETURNING event_id
`

func (s *Storage) Create(ctx context.Context, event *storage.Event) (int64, error) {
	var countExistingEvents int
	err := s.db.GetContext(ctx, &countExistingEvents, s.db.Rebind(checkExistingEventsSQL), event.EndDate, event.StartDate)
	if err != nil {
		return 0, fmt.Errorf("cannot check events: %w", err)
	}
	if countExistingEvents > 0 {
		return 0, storage.ErrBusyTime
	}

	stmt, err := s.db.PrepareNamedContext(ctx, createEventSQL)
	if err != nil {
		return 0, fmt.Errorf("cannot prepare context for creating event: %w", err)
	}

	rows, err := stmt.QueryxContext(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("cannot query context for creating event: %w", err)
	}

	if !rows.Next() {
		return 0, errors.New("event not created")
	}

	var eventID int64
	err = rows.Scan(&eventID)
	return eventID, err
}

const getEventByIDSQL = `
SELECT event_id, title, start_date, end_date, description, user_id, notify_before
FROM events
WHERE event_id = :event_id
`

func (s *Storage) GetByID(ctx context.Context, eventID int64) (*storage.Event, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, getEventByIDSQL)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare context for getting event by id: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, map[string]interface{}{
		"event_id": eventID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot query context for getting event by id: %w", err)
	}

	if !rows.Next() {
		return nil, storage.ErrEventNotFound
	}

	var event *storage.Event
	err = rows.Scan(event)
	return event, err
}

const updateEventSQL = `
UPDATE events
SET title = :title,
	start_date = :start_date,
	end_date = :end_date,
	description = :description,
	user_id = :user_id,
	notify_before = :notify_before
WHERE event_id = :event_id
`

func (s *Storage) Update(ctx context.Context, event *storage.Event) error {
	stmt, err := s.db.PrepareNamedContext(ctx, updateEventSQL)
	if err != nil {
		return fmt.Errorf("cannot prepare context for updating event: %w", err)
	}

	result, err := stmt.ExecContext(ctx, event)
	if err != nil {
		return fmt.Errorf("cannot query context for updating event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}
	return nil
}

const deleteEventSQL = `
DELETE FROM events
WHERE event_id = :event_id
`

func (s *Storage) Delete(ctx context.Context, eventID int64) error {
	stmt, err := s.db.PrepareNamedContext(ctx, deleteEventSQL)
	if err != nil {
		return fmt.Errorf("cannot prepare context for deleting event: %w", err)
	}

	result, err := stmt.ExecContext(ctx, map[string]interface{}{
		"event_id": eventID,
	})
	if err != nil {
		return fmt.Errorf("cannot query context for deleting event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}
	return nil
}

const selectEventsByDatesSQL = `
SELECT event_id, title, start_date, end_date, description, user_id, notify_before
FROM events
WHERE start_date >= :start_date and end_date < :end_date
ORDER BY start_date, end_date
`

func (s *Storage) ListForPeriod(
	ctx context.Context,
	startDate time.Time,
	endDateExclusive time.Time,
) ([]*storage.Event, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, selectEventsByDatesSQL)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare context for listing events: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDateExclusive,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot query context for listing events: %w", err)
	}

	events := make([]*storage.Event, 0)
	for rows.Next() {
		var event *storage.Event
		err = rows.Scan(event)
		if err != nil {
			return nil, fmt.Errorf("cannot get result for listing events: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}
