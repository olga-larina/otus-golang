package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/grpc/mocks"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
Генерация моков (кроме файлов в pb):
go install github.com/vektra/mockery/v2@v2.43.2
mockery --all --exclude internal/server/grpc/pb --case underscore

	--keeptree --dir internal/server/grpc --output internal/server/grpc/mocks --with-expecter --log-level warn.
*/
func TestServer(t *testing.T) {
	t.Parallel()

	userID := uint64(12345)
	userIDStr := "12345"
	eventDto := app.EventDto{
		ID:           1,
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}
	eventPb := pb.Event{
		Id:           1,
		Title:        "my event",
		StartDate:    timestamppb.New(getTime(t, "2024-07-06 10:00:00")),
		EndDate:      timestamppb.New(getTime(t, "2024-07-10 00:00:00")),
		Description:  "my event description",
		NotifyBefore: durationpb.New(time.Hour * 24),
	}

	t.Run("create event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		eventDto := eventDto
		eventPb := proto.Clone(&eventPb).(*pb.Event)
		eventID := uint64(1000)

		mockedApplication.EXPECT().Create(ctx, eventDto).Return(eventID, nil)

		actualCreateResponse, err := server.CreateEvent(ctx, &pb.CreateEventRequest{Event: eventPb})
		require.NoError(t, err)
		require.Equal(t, eventID, actualCreateResponse.Id)
	})

	t.Run("create event with busy time", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		eventDto := eventDto
		eventPb := proto.Clone(&eventPb).(*pb.Event)
		eventID := uint64(1000)

		mockedApplication.EXPECT().Create(ctx, eventDto).Return(eventID, storage.ErrBusyTime)

		_, err := server.CreateEvent(ctx, &pb.CreateEventRequest{Event: eventPb})
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("create event without user", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		ctx := context.Background()
		eventPb := proto.Clone(&eventPb).(*pb.Event)

		_, err := server.CreateEvent(ctx, &pb.CreateEventRequest{Event: eventPb})
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("update event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		eventDto := eventDto
		eventPb := proto.Clone(&eventPb).(*pb.Event)

		mockedApplication.EXPECT().Update(ctx, eventDto).Return(nil)

		_, err := server.UpdateEvent(ctx, &pb.UpdateEventRequest{Event: eventPb})
		require.NoError(t, err)
	})

	t.Run("delete event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		eventID := uint64(1000)

		mockedApplication.EXPECT().Delete(ctx, userID, eventID).Return(nil)

		_, err := server.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: eventID})
		require.NoError(t, err)
	})

	t.Run("get event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		eventDto := eventDto
		eventPb := proto.Clone(&eventPb).(*pb.Event)
		eventID := uint64(1000)

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockedApplication.EXPECT().GetByID(ctx, userID, eventID).Return(&eventDto, nil)

		actualEvent, err := server.GetEvent(ctx, &pb.GetEventRequest{Id: eventID})
		require.NoError(t, err)
		require.True(t, proto.Equal(eventPb, actualEvent))
	})

	t.Run("get not existing event", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		eventID := uint64(1000)

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockedApplication.EXPECT().GetByID(ctx, userID, eventID).Return(nil, storage.ErrEventNotFound)

		_, err := server.GetEvent(ctx, &pb.GetEventRequest{Id: eventID})
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})
}

func TestServerListEvents(t *testing.T) {
	t.Parallel()

	userID := uint64(12345)
	userIDStr := "12345"

	eventDto0 := app.EventDto{
		ID:           1,
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}
	eventDto1 := app.EventDto{
		ID:           2,
		Title:        "my event2",
		StartDate:    getTime(t, "2024-07-05 10:00:00"),
		EndDate:      getTime(t, "2024-07-06 09:59:59"),
		Description:  "my event description2",
		UserID:       userID,
		NotifyBefore: time.Hour * 12,
	}

	eventPb0 := pb.Event{
		Id:           1,
		Title:        "my event",
		StartDate:    timestamppb.New(getTime(t, "2024-07-06 10:00:00")),
		EndDate:      timestamppb.New(getTime(t, "2024-07-10 00:00:00")),
		Description:  "my event description",
		NotifyBefore: durationpb.New(time.Hour * 24),
	}
	eventPb1 := pb.Event{
		Id:           2,
		Title:        "my event2",
		StartDate:    timestamppb.New(getTime(t, "2024-07-05 10:00:00")),
		EndDate:      timestamppb.New(getTime(t, "2024-07-06 09:59:59")),
		Description:  "my event description2",
		NotifyBefore: durationpb.New(time.Hour * 12),
	}

	t.Run("list events for day", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		startDate := getTime(t, "2024-07-06 10:00:00")

		eventDto0 := eventDto0
		eventDto1 := eventDto1
		eventPb0 := proto.Clone(&eventPb0).(*pb.Event)
		eventPb1 := proto.Clone(&eventPb1).(*pb.Event)

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockedApplication.EXPECT().ListForDay(ctx, userID, startDate).Return([]*app.EventDto{&eventDto0, &eventDto1}, nil)

		actualEvents, err := server.EventListForDay(ctx, &pb.EventListRequest{StartDate: timestamppb.New(startDate)})
		require.NoError(t, err)
		require.Equal(t, 2, len(actualEvents.Events))
		require.True(t, proto.Equal(eventPb0, actualEvents.Events[0]))
		require.True(t, proto.Equal(eventPb1, actualEvents.Events[1]))
	})

	t.Run("list events for week", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		startDate := getTime(t, "2024-07-06 10:00:00")

		eventDto0 := eventDto0
		eventDto1 := eventDto1
		eventPb0 := proto.Clone(&eventPb0).(*pb.Event)
		eventPb1 := proto.Clone(&eventPb1).(*pb.Event)

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockedApplication.EXPECT().ListForWeek(ctx, userID, startDate).Return([]*app.EventDto{&eventDto0, &eventDto1}, nil)

		actualEvents, err := server.EventListForWeek(ctx, &pb.EventListRequest{StartDate: timestamppb.New(startDate)})
		require.NoError(t, err)
		require.Equal(t, 2, len(actualEvents.Events))
		require.True(t, proto.Equal(eventPb0, actualEvents.Events[0]))
		require.True(t, proto.Equal(eventPb1, actualEvents.Events[1]))
	})

	t.Run("list events for month", func(t *testing.T) {
		mockedLogger := mocks.NewLogger(t)
		mockedApplication := mocks.NewApplication(t)
		server := NewServer(mockedLogger, mockedApplication, "")

		startDate := getTime(t, "2024-07-06 10:00:00")

		eventDto0 := eventDto0
		eventDto1 := eventDto1
		eventPb0 := proto.Clone(&eventPb0).(*pb.Event)
		eventPb1 := proto.Clone(&eventPb1).(*pb.Event)

		md := make(metadata.MD)
		md[userIDHeader] = []string{userIDStr}
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockedApplication.EXPECT().ListForMonth(ctx, userID, startDate).Return([]*app.EventDto{&eventDto0, &eventDto1}, nil)

		actualEvents, err := server.EventListForMonth(ctx, &pb.EventListRequest{StartDate: timestamppb.New(startDate)})
		require.NoError(t, err)
		require.Equal(t, 2, len(actualEvents.Events))
		require.True(t, proto.Equal(eventPb0, actualEvents.Events[0]))
		require.True(t, proto.Equal(eventPb1, actualEvents.Events[1]))
	})
}

func getTime(t *testing.T, value string) time.Time {
	t.Helper()
	time, err := time.Parse(time.DateTime, value)
	require.NoError(t, err)
	return time
}
