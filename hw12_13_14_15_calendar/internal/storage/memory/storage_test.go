package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Parallel()

	event := storage.Event{
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       12345,
		NotifyBefore: time.Hour * 24,
	}

	event1 := event
	event1.StartDate = getTime(t, "2024-07-05 10:00:00")
	event1.EndDate = getTime(t, "2024-07-06 09:59:59")

	event2 := event
	event2.StartDate = getTime(t, "2024-07-10 00:00:01")
	event2.EndDate = getTime(t, "2024-07-10 10:00:00")

	event3 := event
	event3.StartDate = getTime(t, "2024-07-05 10:00:00")
	event3.EndDate = getTime(t, "2024-07-10 09:59:59")

	ctx := context.Background()

	t.Run("create new event", func(t *testing.T) {
		s := New()
		event := event

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)
		require.Greater(t, eventID, int64(0))
	})

	t.Run("create new event without intersections", func(t *testing.T) {
		s := New()
		event := event
		event1 := event1
		event2 := event2

		_, err := s.Create(ctx, &event1)
		require.NoError(t, err)

		_, err = s.Create(ctx, &event2)
		require.NoError(t, err)

		_, err = s.Create(ctx, &event)
		require.NoError(t, err)
	})

	t.Run("create new event in busy time", func(t *testing.T) {
		s := New()
		event := event
		event3 := event3

		_, err := s.Create(ctx, &event3)
		require.NoError(t, err)

		_, err = s.Create(ctx, &event)
		require.ErrorIs(t, err, storage.ErrBusyTime)
	})

	t.Run("get created event", func(t *testing.T) {
		s := New()
		event := event

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)

		expectedEvent := event
		expectedEvent.ID = eventID

		actualEvent, err := s.GetByID(ctx, eventID)
		require.NoError(t, err)
		require.Equal(t, expectedEvent, *actualEvent)
	})

	t.Run("get not existing event", func(t *testing.T) {
		s := New()
		_, err := s.GetByID(ctx, int64(0))
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("update created event", func(t *testing.T) {
		s := New()
		event := event

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)

		updatedEvent := storage.Event{
			ID:           eventID,
			Title:        "my event 2",
			StartDate:    getTime(t, "2025-07-06 10:00:00"),
			EndDate:      getTime(t, "2025-07-10 00:00:00"),
			Description:  "my event 2 description",
			UserID:       54321,
			NotifyBefore: time.Hour * 1,
		}

		err = s.Update(ctx, &updatedEvent)
		require.NoError(t, err)

		actualEvent, err := s.GetByID(ctx, eventID)
		require.NoError(t, err)
		require.Equal(t, updatedEvent, *actualEvent)
	})

	t.Run("update not existing event", func(t *testing.T) {
		s := New()
		event := event

		err := s.Update(ctx, &event)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("delete created event", func(t *testing.T) {
		s := New()
		event := event

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)

		err = s.Delete(ctx, eventID)
		require.NoError(t, err)

		_, err = s.GetByID(ctx, int64(0))
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("delete not existing event", func(t *testing.T) {
		s := New()
		err := s.Delete(ctx, int64(0))
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}

func TestStorageList(t *testing.T) {
	t.Parallel()

	event := storage.Event{
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       12345,
		NotifyBefore: time.Hour * 24,
	}

	event1 := event
	event1.StartDate = getTime(t, "2024-07-05 10:00:00")
	event1.EndDate = getTime(t, "2024-07-06 09:59:59")

	event2 := event
	event2.StartDate = getTime(t, "2024-07-10 00:00:01")
	event2.EndDate = getTime(t, "2024-07-10 10:00:00")

	event3 := event
	event3.StartDate = getTime(t, "2024-07-05 10:00:00")
	event3.EndDate = getTime(t, "2024-07-10 09:59:59")

	ctx := context.Background()

	t.Run("list events in period (all)", func(t *testing.T) {
		s := New()
		event := event
		event1 := event1
		event2 := event2

		events := []*storage.Event{
			&event1,
			&event,
			&event2,
		}

		for i, ev := range events {
			eventID, err := s.Create(ctx, ev)
			require.NoError(t, err)
			events[i].ID = eventID
		}

		actualEvents, err := s.ListForPeriod(ctx, getTime(t, "2024-07-06 09:59:59"), getTime(t, "2024-07-10 00:00:02"))
		require.NoError(t, err)
		require.EqualValues(t, events, actualEvents)
	})

	t.Run("list events in period (zero)", func(t *testing.T) {
		s := New()
		event := event
		event1 := event1
		event2 := event2

		events := []*storage.Event{
			&event1,
			&event,
			&event2,
		}

		for i, ev := range events {
			eventID, err := s.Create(ctx, ev)
			require.NoError(t, err)
			events[i].ID = eventID
		}

		actualEvents, err := s.ListForPeriod(ctx, getTime(t, "2024-07-05 00:00:00"), getTime(t, "2024-07-05 09:59:59"))
		require.NoError(t, err)
		require.Equal(t, 0, len(actualEvents))
	})
}

func getTime(t *testing.T, value string) time.Time {
	t.Helper()
	time, err := time.Parse(time.DateTime, value)
	require.NoError(t, err)
	return time
}
