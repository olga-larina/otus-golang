//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func (s *IntegrationTestSuite) TestHttpProcessEvent() {
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
	eventDto := app.EventDto{
		Title:        event.Title,
		StartDate:    event.StartDate,
		EndDate:      event.EndDate,
		Description:  event.Description,
		NotifyBefore: event.NotifyBefore,
	}
	eventDtoUpdated := app.EventDto{
		Title:        event.Title + " v2",
		StartDate:    event.StartDate.Add(time.Minute * 20),
		EndDate:      event.EndDate.Add(time.Minute * 30),
		Description:  event.Description + " v2",
		NotifyBefore: time.Minute * 10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// добавление нового события
	resp, err := s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetBody(internalhttp.CreateEventRequest{Event: &eventDto}).
		SetResult(internalhttp.CreateEventResponse{}).
		Post("/events")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())
	createEventResponse := resp.Result().(*internalhttp.CreateEventResponse)
	event.ID = createEventResponse.EventID
	eventDto.ID = createEventResponse.EventID
	eventDto.UserID = event.UserID
	eventDtoUpdated.ID = createEventResponse.EventID
	eventDtoUpdated.UserID = event.UserID

	// попытка добавления события без пользователя
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetBody(internalhttp.CreateEventRequest{Event: &eventDto}).
		Post("/events")
	require.NoError(t, err)
	require.True(t, resp.IsError())
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Contains(t, resp.String(), "userID is not valid")

	// попытка добавления пересекающегося события
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetBody(internalhttp.CreateEventRequest{Event: &eventDto}).
		Post("/events")
	require.NoError(t, err)
	require.True(t, resp.IsError())
	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	require.Contains(t, resp.String(), storage.ErrBusyTime.Error())

	// получение события
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetPathParams(map[string]string{
			"eventID": strconv.FormatUint(eventDto.ID, 10),
		}).
		SetResult(internalhttp.EventResponse{}).
		Get("/events/{eventID}")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())
	getEventResponse := resp.Result().(*internalhttp.EventResponse)
	require.Equal(t, eventDto, *getEventResponse.Event)

	// обновление события
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetPathParams(map[string]string{
			"eventID": strconv.FormatUint(eventDto.ID, 10),
		}).
		SetBody(internalhttp.UpdateEventRequest{Event: &eventDtoUpdated}).
		Put("/events/{eventID}")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())

	// получение обновлённого события
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetPathParams(map[string]string{
			"eventID": strconv.FormatUint(eventDto.ID, 10),
		}).
		SetResult(internalhttp.EventResponse{}).
		Get("/events/{eventID}")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())
	getEventUpdatedResponse := resp.Result().(*internalhttp.EventResponse)
	require.Equal(t, eventDtoUpdated, *getEventUpdatedResponse.Event)

	// попытка получения события другим пользователем
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, "12345").
		SetPathParams(map[string]string{
			"eventID": strconv.FormatUint(eventDto.ID, 10),
		}).
		Get("/events/{eventID}")
	require.NoError(t, err)
	require.True(t, resp.IsError())
	require.Equal(t, http.StatusNotFound, resp.StatusCode())
	require.Contains(t, resp.String(), storage.ErrEventNotFound.Error())

	// удаление события
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetPathParams(map[string]string{
			"eventID": strconv.FormatUint(eventDto.ID, 10),
		}).
		Delete("/events/{eventID}")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())

	// попытка запросить удалённое событие
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetPathParams(map[string]string{
			"eventID": strconv.FormatUint(eventDto.ID, 10),
		}).
		Get("/events/{eventID}")
	require.NoError(t, err)
	require.True(t, resp.IsError())
	require.Equal(t, http.StatusNotFound, resp.StatusCode())
	require.Contains(t, resp.String(), storage.ErrEventNotFound.Error())
}

func (s *IntegrationTestSuite) TestHttpListEvents() {
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
	eventDto1 := app.EventDto{
		Title:        event.Title,
		StartDate:    event.StartDate,
		EndDate:      event.EndDate,
		Description:  event.Description,
		NotifyBefore: event.NotifyBefore,
	}
	eventDto2 := app.EventDto{
		Title:        event.Title + " v2",
		StartDate:    event.StartDate.Add(time.Hour * 24 * 5),
		EndDate:      event.EndDate.Add(time.Hour * 25 * 5),
		Description:  event.Description + " v2",
		NotifyBefore: time.Minute * 10,
	}
	eventDto3 := app.EventDto{
		Title:        event.Title + " v3",
		StartDate:    event.StartDate.Add(time.Hour * 24 * 15),
		EndDate:      event.EndDate.Add(time.Hour * 25 * 15),
		Description:  event.Description + " v3",
		NotifyBefore: time.Minute * 30,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// добавление событий
	eventsDto := []*app.EventDto{&eventDto1, &eventDto2, &eventDto3}
	for _, eventDto := range eventsDto {
		resp, err := s.httpClient.R().
			// SetDebug(true).
			SetContext(ctx).
			SetHeader(userIDHeader, userIDStr).
			SetBody(internalhttp.CreateEventRequest{Event: eventDto}).
			SetResult(internalhttp.CreateEventResponse{}).
			Post("/events")
		require.NoError(t, err)
		require.True(t, resp.IsSuccess())
		createEventResponse := resp.Result().(*internalhttp.CreateEventResponse)
		eventDto.ID = createEventResponse.EventID
		eventDto.UserID = event.UserID
	}

	// получение событий
	startDate := time.Date(event.StartDate.Year(), event.StartDate.Month(), event.StartDate.Day(), 0, 0, 0, 0, event.StartDate.Location()).Format("2006-01-02")

	// за день
	resp, err := s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetQueryParam("startDate", startDate).
		SetQueryParam("period", "day").
		SetResult(internalhttp.EventsResponse{}).
		Get("/events")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())
	eventListForDayResponse := resp.Result().(*internalhttp.EventsResponse)
	require.Equal(t, 1, len(eventListForDayResponse.Events))
	require.Equal(t, eventDto1, *eventListForDayResponse.Events[0])

	// за неделю
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetQueryParam("startDate", startDate).
		SetQueryParam("period", "week").
		SetResult(internalhttp.EventsResponse{}).
		Get("/events")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())
	eventListForWeekResponse := resp.Result().(*internalhttp.EventsResponse)
	require.Equal(t, 2, len(eventListForWeekResponse.Events))
	require.Equal(t, eventDto1, *eventListForWeekResponse.Events[0])
	require.Equal(t, eventDto2, *eventListForWeekResponse.Events[1])

	// за месяц
	resp, err = s.httpClient.R().
		// SetDebug(true).
		SetContext(ctx).
		SetHeader(userIDHeader, userIDStr).
		SetQueryParam("startDate", startDate).
		SetQueryParam("period", "month").
		SetResult(internalhttp.EventsResponse{}).
		Get("/events")
	require.NoError(t, err)
	require.True(t, resp.IsSuccess())
	eventListForMonthResponse := resp.Result().(*internalhttp.EventsResponse)
	require.Equal(t, 3, len(eventListForMonthResponse.Events))
	require.Equal(t, eventDto1, *eventListForMonthResponse.Events[0])
	require.Equal(t, eventDto2, *eventListForMonthResponse.Events[1])
	require.Equal(t, eventDto3, *eventListForMonthResponse.Events[2])

	// удаление событий
	for _, eventDto := range eventsDto {
		resp, err = s.httpClient.R().
			// SetDebug(true).
			SetContext(ctx).
			SetHeader(userIDHeader, userIDStr).
			SetPathParams(map[string]string{
				"eventID": strconv.FormatUint(eventDto.ID, 10),
			}).
			Delete("/events/{eventID}")
		require.NoError(t, err)
		require.True(t, resp.IsSuccess())
	}
}
