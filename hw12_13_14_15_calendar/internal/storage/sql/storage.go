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

const eventFields = `
event_id, title, start_date, end_date, description, user_id, 
CAST(EXTRACT(EPOCH FROM notify_before) * 1000000000 as BIGINT) AS notify_before, notify_status
`

const checkExistingEventsSQL = `
SELECT count(*) FROM events WHERE start_date <= ? AND end_date >= ? AND user_id = ? AND event_id <> ?
`

const createEventSQL = `
INSERT INTO events (title, start_date, end_date, description, user_id, notify_before)
VALUES (:title, :start_date, :end_date, :description, :user_id, :notify_before)
RETURNING event_id
`

func (s *Storage) Create(ctx context.Context, event *storage.Event) (uint64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var countExistingEvents int
	err = tx.GetContext(ctx, &countExistingEvents, tx.Rebind(checkExistingEventsSQL), event.EndDate, event.StartDate, event.UserID, 0)
	if err != nil {
		return 0, fmt.Errorf("cannot check events: %w", err)
	}
	if countExistingEvents > 0 {
		return 0, storage.ErrBusyTime
	}

	stmt, err := tx.PrepareNamedContext(ctx, createEventSQL)
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

	var eventID uint64
	err = rows.Scan(&eventID)
	if err == nil {
		err = tx.Commit()
	}
	return eventID, err
}

const getEventByIDSQL = `
SELECT ` + eventFields + `
FROM events
WHERE event_id = :event_id AND user_id = :user_id
`

func (s *Storage) GetByID(ctx context.Context, userID uint64, eventID uint64) (*storage.Event, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, getEventByIDSQL)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare context for getting event by id: %w", err)
	}

	rows, err := stmt.QueryxContext(ctx, map[string]interface{}{
		"event_id": eventID,
		"user_id":  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot query context for getting event by id: %w", err)
	}

	if !rows.Next() {
		return nil, storage.ErrEventNotFound
	}

	var event storage.Event
	err = rows.StructScan(&event)
	return &event, err
}

const updateEventSQL = `
UPDATE events
SET title = :title,
	start_date = :start_date,
	end_date = :end_date,
	description = :description,
	user_id = :user_id,
	notify_before = :notify_before
WHERE event_id = :event_id AND user_id = :user_id
`

func (s *Storage) Update(ctx context.Context, event *storage.Event) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var countExistingEvents int
	err = tx.GetContext(ctx, &countExistingEvents, tx.Rebind(checkExistingEventsSQL), event.EndDate, event.StartDate, event.UserID, event.ID)
	if err != nil {
		return fmt.Errorf("cannot check events: %w", err)
	}
	if countExistingEvents > 0 {
		return storage.ErrBusyTime
	}

	stmt, err := tx.PrepareNamedContext(ctx, updateEventSQL)
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
	return tx.Commit()
}

const deleteEventSQL = `
DELETE FROM events
WHERE event_id = :event_id AND user_id = :user_id
`

func (s *Storage) Delete(ctx context.Context, userID uint64, eventID uint64) error {
	stmt, err := s.db.PrepareNamedContext(ctx, deleteEventSQL)
	if err != nil {
		return fmt.Errorf("cannot prepare context for deleting event: %w", err)
	}

	result, err := stmt.ExecContext(ctx, map[string]interface{}{
		"event_id": eventID,
		"user_id":  userID,
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
SELECT ` + eventFields + `
FROM events
WHERE start_date < :end_date AND end_date >= :start_date AND user_id = :user_id
ORDER BY start_date, end_date
`

func (s *Storage) ListForPeriod(
	ctx context.Context,
	userID uint64,
	startDate time.Time,
	endDateExclusive time.Time,
) ([]*storage.Event, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, selectEventsByDatesSQL)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare context for listing events: %w", err)
	}

	rows, err := stmt.QueryxContext(ctx, map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDateExclusive,
		"user_id":    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot query context for listing events: %w", err)
	}

	events := make([]*storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		err = rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("cannot get result for listing events: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

const selectEventsForNotifySQL = `
SELECT ` + eventFields + `
FROM events
WHERE start_date - notify_before BETWEEN :start_notify_date AND :end_notify_date AND notify_status = :notify_status
ORDER BY start_date - notify_before, start_date, end_date
`

func (s *Storage) ListForNotify(ctx context.Context, startNotifyDate time.Time, endNotifyDate time.Time) ([]*storage.Event, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, selectEventsForNotifySQL)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare context for listing events for notify: %w", err)
	}

	rows, err := stmt.QueryxContext(ctx, map[string]interface{}{
		"notify_status":     storage.NotNotified,
		"start_notify_date": startNotifyDate,
		"end_notify_date":   endNotifyDate,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot query context for listing events for notify: %w", err)
	}

	events := make([]*storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		err = rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("cannot get result for listing events for notify: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

const setNotifyStatusSQL = `
UPDATE events
SET notify_status = :notify_status
WHERE event_id = ANY(:event_ids)
`

func (s *Storage) SetNotifyStatus(ctx context.Context, eventIDs []uint64, notifyStatus storage.NotifyStatus) error {
	stmt, err := s.db.PrepareNamedContext(ctx, setNotifyStatusSQL)
	if err != nil {
		return fmt.Errorf("cannot prepare context for setting notify status of events: %w", err)
	}

	_, err = stmt.ExecContext(ctx, map[string]interface{}{
		"notify_status": notifyStatus,
		"event_ids":     eventIDs,
	})
	return err
}

const deleteEventsByEndDateSQL = `
DELETE FROM events
WHERE end_date <= :max_end_date
`

func (s *Storage) DeleteByEndDate(ctx context.Context, maxEndDate time.Time) error {
	stmt, err := s.db.PrepareNamedContext(ctx, deleteEventsByEndDateSQL)
	if err != nil {
		return fmt.Errorf("cannot prepare context for deleting events by end date: %w", err)
	}

	_, err = stmt.ExecContext(ctx, map[string]interface{}{
		"max_end_date": maxEndDate,
	})

	return err
}
