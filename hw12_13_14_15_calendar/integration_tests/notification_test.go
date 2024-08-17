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
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *IntegrationTestSuite) TestEventNotification() {
	t := s.T()

	var err error

	now := time.Now()
	userID := uint64(now.UnixMilli())
	userIDStr := strconv.FormatUint(userID, 10)
	event := storage.Event{
		Title:        "my event",
		StartDate:    now.Add(s.cfg.Calendar.NotifyPeriod / 2).Truncate(time.Millisecond),
		EndDate:      now.Add(time.Hour).Truncate(time.Millisecond),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: s.cfg.Calendar.NotifyScanPeriod / 2,
	}
	eventPb := pb.Event{
		Title:        event.Title,
		StartDate:    timestamppb.New(event.StartDate),
		EndDate:      timestamppb.New(event.EndDate),
		Description:  event.Description,
		NotifyBefore: durationpb.New(event.NotifyBefore),
	}

	md := make(metadata.MD)
	md[userIDHeader] = []string{userIDStr}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	// добавление нового события
	createEventResponse, err := s.grpcClient.CreateEvent(ctx, &pb.CreateEventRequest{Event: &eventPb})
	require.NoError(t, err)
	eventPb.Id = createEventResponse.Id

	// ожидание отправки нотификации
	require.Eventually(t, func() bool {
		select {
		case <-ctx.Done():
			return false
		default:
			actualEvent, err := s.storage.GetByID(ctx, userID, eventPb.Id)
			if err != nil {
				s.logg.Error(ctx, err, "failed to get event")
				return false
			}
			return actualEvent.NotifyStatus == storage.Notified
		}
	}, 3*s.cfg.Calendar.NotifyCronPeriod, time.Second, "event was notified?")
}
