//go:build integration
// +build integration

package integration

import (
	"context"
	"strconv"
	"time"

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

func (s *IntegrationTestSuite) TestGrpcProcessEvent() {
	t := s.T()

	var err error

	now := time.Now()
	userID := uint64(now.UnixMilli())
	userIDStr := strconv.FormatUint(userID, 10)
	event := storage.Event{
		Title:        "my event",
		StartDate:    now.Add(time.Minute * 30).Truncate(time.Millisecond),
		EndDate:      now.Add(time.Hour).Truncate(time.Millisecond),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Minute * 5,
	}
	eventPb := pb.Event{
		Title:        event.Title,
		StartDate:    timestamppb.New(event.StartDate),
		EndDate:      timestamppb.New(event.EndDate),
		Description:  event.Description,
		NotifyBefore: durationpb.New(event.NotifyBefore),
	}
	eventPbUpdated := pb.Event{
		Title:        event.Title + " v2",
		StartDate:    timestamppb.New(event.StartDate.Add(time.Minute * 20)),
		EndDate:      timestamppb.New(event.EndDate.Add(time.Minute * 30)),
		Description:  event.Description + " v2",
		NotifyBefore: durationpb.New(time.Minute * 10),
	}

	md := make(metadata.MD)
	md[userIDHeader] = []string{userIDStr}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	// добавление нового события
	createEventResponse, err := s.grpcClient.CreateEvent(ctx, &pb.CreateEventRequest{Event: &eventPb})
	require.NoError(t, err)
	event.ID = createEventResponse.Id
	eventPb.Id = createEventResponse.Id
	eventPbUpdated.Id = createEventResponse.Id

	// попытка добавления события без пользователя
	_, err = s.grpcClient.CreateEvent(ctxTimeout, &pb.CreateEventRequest{Event: &eventPb})
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Contains(t, st.Message(), "userID is not valid")

	// попытка добавления пересекающегося события
	_, err = s.grpcClient.CreateEvent(ctx, &pb.CreateEventRequest{Event: &eventPb})
	st, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Contains(t, st.Message(), storage.ErrBusyTime.Error())

	// получение события
	getEventResponse, err := s.grpcClient.GetEvent(ctx, &pb.GetEventRequest{Id: eventPb.Id})
	require.NoError(t, err)
	require.True(t, proto.Equal(&eventPb, getEventResponse))

	// обновление события
	_, err = s.grpcClient.UpdateEvent(ctx, &pb.UpdateEventRequest{Event: &eventPbUpdated})
	require.NoError(t, err)

	// получение обновлённого события
	getEventUpdatedResponse, err := s.grpcClient.GetEvent(ctx, &pb.GetEventRequest{Id: eventPb.Id})
	require.NoError(t, err)
	require.True(t, proto.Equal(&eventPbUpdated, getEventUpdatedResponse))

	// попытка получения события другим пользователем
	mdOther := make(metadata.MD)
	mdOther[userIDHeader] = []string{"12345"}
	ctxOther := metadata.NewOutgoingContext(ctxTimeout, mdOther)
	_, err = s.grpcClient.GetEvent(ctxOther, &pb.GetEventRequest{Id: eventPb.Id})
	st, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
	require.Contains(t, st.Message(), storage.ErrEventNotFound.Error())

	// удаление события
	_, err = s.grpcClient.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: eventPb.Id})
	require.NoError(t, err)

	// попытка запросить удалённое событие
	_, err = s.grpcClient.GetEvent(ctx, &pb.GetEventRequest{Id: eventPb.Id})
	st, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
	require.Contains(t, st.Message(), storage.ErrEventNotFound.Error())
}

func (s *IntegrationTestSuite) TestGrpcListEvents() {
	t := s.T()

	var err error

	now := time.Now()
	userID := uint64(now.UnixMilli())
	userIDStr := strconv.FormatUint(userID, 10)
	event := storage.Event{
		Title:        "my event1",
		StartDate:    now.Add(time.Minute * 30).Truncate(time.Millisecond),
		EndDate:      now.Add(time.Hour).Truncate(time.Millisecond),
		Description:  "my event1 description",
		UserID:       userID,
		NotifyBefore: time.Minute * 5,
	}
	eventPb1 := pb.Event{
		Title:        event.Title,
		StartDate:    timestamppb.New(event.StartDate),
		EndDate:      timestamppb.New(event.EndDate),
		Description:  event.Description,
		NotifyBefore: durationpb.New(event.NotifyBefore),
	}
	eventPb2 := pb.Event{
		Title:        event.Title + " v2",
		StartDate:    timestamppb.New(event.StartDate.Add(time.Hour * 24 * 5)),
		EndDate:      timestamppb.New(event.EndDate.Add(time.Hour * 25 * 5)),
		Description:  event.Description + " v2",
		NotifyBefore: durationpb.New(time.Minute * 10),
	}
	eventPb3 := pb.Event{
		Title:        event.Title + " v3",
		StartDate:    timestamppb.New(event.StartDate.Add(time.Hour * 24 * 15)),
		EndDate:      timestamppb.New(event.EndDate.Add(time.Hour * 25 * 15)),
		Description:  event.Description + " v3",
		NotifyBefore: durationpb.New(time.Minute * 30),
	}

	md := make(metadata.MD)
	md[userIDHeader] = []string{userIDStr}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	// добавление событий
	eventsPb := []*pb.Event{&eventPb1, &eventPb2, &eventPb3}
	for _, eventPb := range eventsPb {
		createEventResponse, err := s.grpcClient.CreateEvent(ctx, &pb.CreateEventRequest{Event: eventPb})
		require.NoError(t, err)
		eventPb.Id = createEventResponse.Id
	}

	// получение событий
	startDate := time.Date(event.StartDate.Year(), event.StartDate.Month(), event.StartDate.Day(), 0, 0, 0, 0, event.StartDate.Location())

	// за день
	eventListForDayResponse, err := s.grpcClient.EventListForDay(ctx, &pb.EventListRequest{StartDate: timestamppb.New(startDate)})
	require.NoError(t, err)
	require.Equal(t, 1, len(eventListForDayResponse.Events))
	require.True(t, proto.Equal(&eventPb1, eventListForDayResponse.Events[0]))

	// за неделю
	eventListForWeekResponse, err := s.grpcClient.EventListForWeek(ctx, &pb.EventListRequest{StartDate: timestamppb.New(startDate)})
	require.NoError(t, err)
	require.Equal(t, 2, len(eventListForWeekResponse.Events))
	require.True(t, proto.Equal(&eventPb1, eventListForWeekResponse.Events[0]))
	require.True(t, proto.Equal(&eventPb2, eventListForWeekResponse.Events[1]))

	// за месяц
	eventListForMonthResponse, err := s.grpcClient.EventListForMonth(ctx, &pb.EventListRequest{StartDate: timestamppb.New(startDate)})
	require.NoError(t, err)
	require.Equal(t, 3, len(eventListForMonthResponse.Events))
	require.True(t, proto.Equal(&eventPb1, eventListForMonthResponse.Events[0]))
	require.True(t, proto.Equal(&eventPb2, eventListForMonthResponse.Events[1]))
	require.True(t, proto.Equal(&eventPb3, eventListForMonthResponse.Events[2]))

	// удаление событий
	for _, eventPb := range eventsPb {
		_, err := s.grpcClient.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: eventPb.Id})
		require.NoError(t, err)
	}
}
