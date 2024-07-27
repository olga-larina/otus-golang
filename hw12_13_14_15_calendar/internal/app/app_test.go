package app

import (
	"context"
	"testing"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app/mocks"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

/*
Генерация моков:
go install github.com/vektra/mockery/v2@v2.43.2
mockery --all --case underscore --keeptree --dir internal/app --output internal/app/mocks --with-expecter --log-level warn.
*/
func TestApp(t *testing.T) {
	t.Parallel()

	userID := uint64(12345)
	event := storage.Event{
		ID:           1,
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}
	eventDto := EventDto{
		ID:           1,
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}

	ctx := context.Background()

	t.Run("create event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		eventDto := eventDto
		event := event
		eventID := uint64(1000)

		mockedStorage.EXPECT().Create(ctx, &event).Return(eventID, nil)

		actualEventID, err := app.Create(ctx, eventDto)
		require.NoError(t, err)
		require.Equal(t, eventID, actualEventID)
	})

	t.Run("update event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		eventDto := eventDto
		event := event

		mockedStorage.EXPECT().Update(ctx, &event).Return(nil)

		err := app.Update(ctx, eventDto)
		require.NoError(t, err)
	})

	t.Run("delete event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		eventID := uint64(1000)

		mockedStorage.EXPECT().Delete(ctx, userID, eventID).Return(nil)

		err := app.Delete(ctx, userID, eventID)
		require.NoError(t, err)
	})

	t.Run("get event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		eventID := uint64(1000)

		eventDto := eventDto
		event := event

		mockedStorage.EXPECT().GetByID(ctx, userID, eventID).Return(&event, nil)

		actualEventDto, err := app.GetByID(ctx, userID, eventID)
		require.NoError(t, err)
		require.Equal(t, eventDto, *actualEventDto)
	})

	t.Run("get event with error", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		eventID := uint64(1000)

		mockedStorage.EXPECT().GetByID(ctx, userID, eventID).Return(nil, storage.ErrEventNotFound)

		actualEventDto, err := app.GetByID(ctx, userID, eventID)
		require.Nil(t, actualEventDto)
		require.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}

func TestAppListEvents(t *testing.T) {
	t.Parallel()

	userID := uint64(12345)
	event0 := storage.Event{
		ID:           1,
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}
	event1 := event0
	event1.ID = 2
	event1.StartDate = getTime(t, "2024-07-05 10:00:00")
	event1.EndDate = getTime(t, "2024-07-06 09:59:59")

	eventDto0 := EventDto{
		ID:           1,
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}
	eventDto1 := eventDto0
	eventDto1.ID = 2
	eventDto1.StartDate = getTime(t, "2024-07-05 10:00:00")
	eventDto1.EndDate = getTime(t, "2024-07-06 09:59:59")

	ctx := context.Background()

	t.Run("list events for day", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		startDate := getTime(t, "2024-07-06 10:00:00")
		endDate := getTime(t, "2024-07-07 10:00:00")

		eventDto0 := eventDto0
		eventDto1 := eventDto1
		event0 := event0
		event1 := event1

		mockedStorage.EXPECT().ListForPeriod(ctx, userID, startDate, endDate).Return([]*storage.Event{&event0, &event1}, nil)

		actualEvents, err := app.ListForDay(ctx, userID, startDate)
		require.NoError(t, err)
		require.Equal(t, 2, len(actualEvents))
		require.Equal(t, eventDto0, *actualEvents[0])
		require.Equal(t, eventDto1, *actualEvents[1])
	})

	t.Run("list events for week", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		startDate := getTime(t, "2024-07-06 10:00:00")
		endDate := getTime(t, "2024-07-13 10:00:00")

		eventDto0 := eventDto0
		eventDto1 := eventDto1
		event0 := event0
		event1 := event1

		mockedStorage.EXPECT().ListForPeriod(ctx, userID, startDate, endDate).Return([]*storage.Event{&event0, &event1}, nil)

		actualEvents, err := app.ListForWeek(ctx, userID, startDate)
		require.NoError(t, err)
		require.Equal(t, 2, len(actualEvents))
		require.Equal(t, eventDto0, *actualEvents[0])
		require.Equal(t, eventDto1, *actualEvents[1])
	})

	t.Run("list events for month", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedStorage := mocks.NewStorage(t)
		app := New(mockedLogger, mockedStorage)

		startDate := getTime(t, "2024-07-06 10:00:00")
		endDate := getTime(t, "2024-08-06 10:00:00")

		eventDto0 := eventDto0
		eventDto1 := eventDto1
		event0 := event0
		event1 := event1

		mockedStorage.EXPECT().ListForPeriod(ctx, userID, startDate, endDate).Return([]*storage.Event{&event0, &event1}, nil)

		actualEvents, err := app.ListForMonth(ctx, userID, startDate)
		require.NoError(t, err)
		require.Equal(t, 2, len(actualEvents))
		require.Equal(t, eventDto0, *actualEvents[0])
		require.Equal(t, eventDto1, *actualEvents[1])
	})
}

func getTime(t *testing.T, value string) time.Time {
	t.Helper()
	time, err := time.Parse(time.DateTime, value)
	require.NoError(t, err)
	return time
}
