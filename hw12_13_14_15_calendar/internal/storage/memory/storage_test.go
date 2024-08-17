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

	userID := uint64(12345)
	event := storage.Event{
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
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
		require.Greater(t, eventID, uint64(0))
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

		actualEvent, err := s.GetByID(ctx, userID, eventID)
		require.NoError(t, err)
		require.Equal(t, expectedEvent, *actualEvent)
	})

	t.Run("get not existing event", func(t *testing.T) {
		s := New()
		_, err := s.GetByID(ctx, userID, uint64(0))
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
			UserID:       userID,
			NotifyBefore: time.Hour * 1,
		}

		err = s.Update(ctx, &updatedEvent)
		require.NoError(t, err)

		actualEvent, err := s.GetByID(ctx, userID, eventID)
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

		err = s.Delete(ctx, userID, eventID)
		require.NoError(t, err)

		_, err = s.GetByID(ctx, userID, uint64(0))
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("delete not existing event", func(t *testing.T) {
		s := New()
		err := s.Delete(ctx, userID, uint64(0))
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}

func TestStorageOtherUser(t *testing.T) {
	t.Parallel()

	userID := uint64(12345)
	otherUserID := uint64(54321)

	event := storage.Event{
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}

	ctx := context.Background()

	t.Run("get created event by another user", func(t *testing.T) {
		s := New()
		event := event

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)

		_, err = s.GetByID(ctx, otherUserID, eventID)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("update created event by another user", func(t *testing.T) {
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
			UserID:       otherUserID,
			NotifyBefore: time.Hour * 1,
		}

		err = s.Update(ctx, &updatedEvent)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})

	t.Run("delete created event by another user", func(t *testing.T) {
		s := New()
		event := event

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)

		err = s.Delete(ctx, otherUserID, eventID)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}

func TestStorageList(t *testing.T) {
	t.Parallel()

	userID := uint64(12345)
	event := storage.Event{
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 25,
	}

	event1 := event
	event1.StartDate = getTime(t, "2024-07-05 10:00:00")
	event1.EndDate = getTime(t, "2024-07-06 09:59:59")
	event1.NotifyBefore = 0

	event2 := event
	event2.StartDate = getTime(t, "2024-07-10 00:00:01")
	event2.EndDate = getTime(t, "2024-07-10 10:00:00")
	event2.NotifyBefore = time.Hour * 24 * 4

	eventOther := event
	eventOther.UserID = uint64(54321)
	eventOther.EndDate = getTime(t, "2024-07-10 00:00:01")

	ctx := context.Background()

	t.Run("list events in period (all)", func(t *testing.T) {
		s := New()
		event := event
		event1 := event1
		event2 := event2
		eventOther := eventOther

		events := []*storage.Event{
			&event1,
			&event,
			&event2,
			&eventOther,
		}

		for i, ev := range events {
			eventID, err := s.Create(ctx, ev)
			require.NoError(t, err)
			events[i].ID = eventID
		}

		expectedEvents := events[:3]

		actualEvents, err := s.ListForPeriod(ctx, userID, getTime(t, "2024-07-06 09:59:59"), getTime(t, "2024-07-10 00:00:02"))
		require.NoError(t, err)
		require.EqualValues(t, expectedEvents, actualEvents)
	})

	t.Run("list events in period (zero)", func(t *testing.T) {
		s := New()
		event := event
		event1 := event1
		event2 := event2
		eventOther := eventOther

		events := []*storage.Event{
			&event1,
			&event,
			&event2,
			&eventOther,
		}

		for i, ev := range events {
			eventID, err := s.Create(ctx, ev)
			require.NoError(t, err)
			events[i].ID = eventID
		}

		actualEvents, err := s.ListForPeriod(ctx, userID, getTime(t, "2024-07-05 00:00:00"), getTime(t, "2024-07-05 09:59:59"))
		require.NoError(t, err)
		require.Equal(t, 0, len(actualEvents))
	})
}

func TestStorageScheduling(t *testing.T) {
	t.Parallel()

	userID := uint64(12345)
	event := storage.Event{
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 25,
	}

	event1 := event
	event1.StartDate = getTime(t, "2024-07-05 10:00:00")
	event1.EndDate = getTime(t, "2024-07-06 09:59:59")
	event1.NotifyBefore = 0

	event2 := event
	event2.StartDate = getTime(t, "2024-07-10 00:00:01")
	event2.EndDate = getTime(t, "2024-07-10 10:00:00")
	event2.NotifyBefore = time.Hour * 24 * 4

	eventOther := event
	eventOther.UserID = uint64(54321)
	eventOther.EndDate = getTime(t, "2024-07-10 00:00:01")

	eventOther1 := event
	eventOther1.UserID = uint64(32154)
	eventOther1.StartDate = getTime(t, "2024-07-05 10:00:00")
	eventOther1.EndDate = getTime(t, "2024-07-10 09:59:59")
	eventOther1.NotifyBefore = time.Hour

	ctx := context.Background()

	t.Run("list events for notify", func(t *testing.T) {
		s := New()
		event := event
		event1 := event1
		event2 := event2
		eventOther1 := eventOther1
		eventOther := eventOther

		events := []*storage.Event{
			&eventOther1,
			&event,
			&eventOther,
			&event2,
			&event1,
		}

		for i, ev := range events {
			eventID, err := s.Create(ctx, ev)
			require.NoError(t, err)
			events[i].ID = eventID
		}

		expectedEvents := events[:4]

		actualEvents, err := s.ListForNotify(ctx, getTime(t, "2024-07-05 09:00:00"), getTime(t, "2024-07-06 09:00:00"))
		require.NoError(t, err)
		require.EqualValues(t, expectedEvents, actualEvents)
	})

	t.Run("set events notify status", func(t *testing.T) {
		s := New()
		event := event
		event1 := event1
		event2 := event2
		eventOther := eventOther

		event1ID, err := s.Create(ctx, &event1)
		require.NoError(t, err)

		event2ID, err := s.Create(ctx, &event2)
		require.NoError(t, err)

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)

		eventOtherID, err := s.Create(ctx, &eventOther)
		require.NoError(t, err)

		err = s.SetNotifyStatus(ctx, []uint64{event1ID, eventID}, storage.NotifyInProgress)
		require.NoError(t, err)

		err = s.SetNotifyStatus(ctx, []uint64{eventOtherID}, storage.Notified)
		require.NoError(t, err)

		actualEvent1, err := s.GetByID(ctx, userID, event1ID)
		require.NoError(t, err)
		require.Equal(t, storage.NotifyInProgress, actualEvent1.NotifyStatus)

		actualEvent2, err := s.GetByID(ctx, userID, event2ID)
		require.NoError(t, err)
		require.Equal(t, storage.NotNotified, actualEvent2.NotifyStatus)

		actualEvent, err := s.GetByID(ctx, userID, eventID)
		require.NoError(t, err)
		require.Equal(t, storage.NotifyInProgress, actualEvent.NotifyStatus)

		actualEventOther, err := s.GetByID(ctx, eventOther.UserID, eventOtherID)
		require.NoError(t, err)
		require.Equal(t, storage.Notified, actualEventOther.NotifyStatus)
	})

	t.Run("delete old events", func(t *testing.T) {
		s := New()

		event := event
		event1 := event1
		event2 := event2

		event1ID, err := s.Create(ctx, &event1)
		require.NoError(t, err)
		event1.ID = event1ID

		event2ID, err := s.Create(ctx, &event2)
		require.NoError(t, err)
		event2.ID = event2ID

		eventID, err := s.Create(ctx, &event)
		require.NoError(t, err)
		event.ID = eventID

		err = s.DeleteByEndDate(ctx, getTime(t, "2024-07-10 00:00:00"))
		require.NoError(t, err)

		actualEvents, err := s.ListForPeriod(ctx, userID, getTime(t, "2024-07-05 10:00:00"), getTime(t, "2024-07-10 10:00:00"))
		require.NoError(t, err)
		require.EqualValues(t, []*storage.Event{&event2}, actualEvents)
	})
}

func getTime(t *testing.T, value string) time.Time {
	t.Helper()
	time, err := time.Parse(time.DateTime, value)
	require.NoError(t, err)
	return time
}
