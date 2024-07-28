package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/http/mocks"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	userID     = uint64(12345)
	userID2    = uint64(54321)
	userID2Str = "54321"
	eventID    = uint64(1)
	eventIDStr = "1"
)

type eventHandlerTest struct {
	requestBody          interface{}
	headers              map[string]string
	method               string
	url                  string
	route                func(mux *mux.Router, handler *EventHandler)
	appCall              func(app *mocks.Application)
	expectedResponseBody interface{}
	expectedResponseCode int
	testName             string
}

/*
Генерация моков:
go install github.com/vektra/mockery/v2@v2.43.2
mockery --all --case underscore --keeptree --dir internal/server/http --output internal/server/http/mocks --with-expecter --log-level warn.
*/
func TestEventHandler(t *testing.T) {
	t.Parallel()

	tests := make([]eventHandlerTest, 0)
	tests = append(tests, initCreateHandlerTests(t)...)
	tests = append(tests, initUpdateHandlerTests(t)...)
	tests = append(tests, initDeleteHandlerTests(t)...)
	tests = append(tests, initGetByIDHandlerTests(t)...)
	tests = append(tests, initListForDayHandlerTests(t)...)
	tests = append(tests, initListForWeekHandlerTests(t)...)
	tests = append(tests, initListForMonthHandlerTests(t)...)

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt := tt
			t.Parallel()

			mockedLogger := mocks.NewLogger(t)
			mockedApplication := mocks.NewApplication(t)
			handler := NewEventHandler(mockedLogger, mockedApplication)

			tt.appCall(mockedApplication)

			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(context.Background(), tt.method, tt.url, bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			for k, v := range tt.headers {
				req.Header.Add(k, v)
			}

			mux := mux.NewRouter()
			tt.route(mux, handler)

			response := httptest.NewRecorder()
			mux.ServeHTTP(response, req)

			require.Equal(t, tt.expectedResponseCode, response.Code)

			if tt.expectedResponseBody != nil {
				responseBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResponseBody, responseBody)
			}
		})
	}
}

func initCreateHandlerTests(t *testing.T) []eventHandlerTest {
	t.Helper()
	createEventResponse, err := json.Marshal(CreateEventResponse{EventID: eventID})
	require.NoError(t, err)

	return []eventHandlerTest{
		{
			requestBody: CreateEventRequest{Event: eventDto(t, userID)},
			headers:     make(map[string]string),
			method:      "POST",
			url:         "/events",
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc("/events", handler.create).Methods("POST")
			},
			appCall:              func(_ *mocks.Application) {},
			expectedResponseBody: nil,
			expectedResponseCode: http.StatusBadRequest,
			testName:             "create event without userID",
		},
		{
			requestBody: CreateEventRequest{Event: eventDto(t, userID)},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "POST",
			url:    "/events",
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc("/events", handler.create).Methods("POST")
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().Create(mock.Anything, *eventDto(t, userID2)).Return(eventID, storage.ErrBusyTime)
			},
			expectedResponseBody: nil,
			expectedResponseCode: http.StatusBadRequest,
			testName:             "create event with ErrBusyTime",
		},
		{
			requestBody: CreateEventRequest{Event: eventDto(t, userID)},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "POST",
			url:    "/events",
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc("/events", handler.create).Methods("POST")
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().Create(mock.Anything, *eventDto(t, userID2)).Return(eventID, errors.New("error"))
			},
			expectedResponseBody: nil,
			expectedResponseCode: http.StatusInternalServerError,
			testName:             "create event with Internal Error",
		},
		{
			requestBody: CreateEventRequest{Event: eventDto(t, userID)},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "POST",
			url:    "/events",
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc("/events", handler.create).Methods("POST")
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().Create(mock.Anything, *eventDto(t, userID2)).Return(eventID, nil)
			},
			expectedResponseBody: createEventResponse,
			expectedResponseCode: http.StatusOK,
			testName:             "create event",
		},
	}
}

func initUpdateHandlerTests(t *testing.T) []eventHandlerTest {
	t.Helper()
	return []eventHandlerTest{
		{
			requestBody: UpdateEventRequest{Event: eventDto(t, userID)},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "PUT",
			url:    fmt.Sprintf("/events/%s", eventIDStr),
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc(fmt.Sprintf("/events/{%s}", eventIDPath), handler.update).Methods("PUT")
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().Update(mock.Anything, *eventDto(t, userID2)).Return(nil)
			},
			expectedResponseBody: nil,
			expectedResponseCode: http.StatusOK,
			testName:             "update event",
		},
	}
}

func initDeleteHandlerTests(t *testing.T) []eventHandlerTest {
	t.Helper()
	return []eventHandlerTest{
		{
			requestBody: []byte{},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "DELETE",
			url:    fmt.Sprintf("/events/%s", eventIDStr),
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc(fmt.Sprintf("/events/{%s}", eventIDPath), handler.delete).Methods("DELETE")
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().Delete(mock.Anything, userID2, eventID).Return(nil)
			},
			expectedResponseBody: nil,
			expectedResponseCode: http.StatusOK,
			testName:             "delete event",
		},
	}
}

func initGetByIDHandlerTests(t *testing.T) []eventHandlerTest {
	t.Helper()
	getEventResponse, err := json.Marshal(EventResponse{Event: eventDto(t, userID2)})
	require.NoError(t, err)

	return []eventHandlerTest{
		{
			requestBody: []byte{},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "GET",
			url:    fmt.Sprintf("/events/%s", eventIDStr),
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc(fmt.Sprintf("/events/{%s}", eventIDPath), handler.getByID).Methods("GET")
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().GetByID(mock.Anything, userID2, eventID).Return(eventDto(t, userID2), nil)
			},
			expectedResponseBody: getEventResponse,
			expectedResponseCode: http.StatusOK,
			testName:             "get event by id",
		},
		{
			requestBody: []byte{},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "GET",
			url:    fmt.Sprintf("/events/%s", eventIDStr),
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc(fmt.Sprintf("/events/{%s}", eventIDPath), handler.getByID).Methods("GET")
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().GetByID(mock.Anything, userID2, eventID).Return(nil, storage.ErrEventNotFound)
			},
			expectedResponseBody: nil,
			expectedResponseCode: http.StatusNotFound,
			testName:             "get not existsing event by id",
		},
	}
}

func initListForDayHandlerTests(t *testing.T) []eventHandlerTest {
	t.Helper()
	events := []*app.EventDto{eventDto(t, userID2), eventDto2(t, userID2)}
	getEventsResponse, err := json.Marshal(EventsResponse{
		Events: events,
	})
	require.NoError(t, err)

	return []eventHandlerTest{
		{
			requestBody: []byte{},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "GET",
			url:    fmt.Sprintf("/events?%s=%s&%s=%s", startDateQueryKey, "2024-08-01", periodTypeQueryKey, periodDayQueryValue),
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc("/events", handler.listForDay).Methods("GET").Queries(
					startDateQueryKey, startDateQueryValue,
					periodTypeQueryKey, periodDayQueryValue,
				)
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().ListForDay(mock.Anything, userID2, getTime(t, "2024-08-01 00:00:00")).Return(events, nil)
			},
			expectedResponseBody: getEventsResponse,
			expectedResponseCode: http.StatusOK,
			testName:             "list events for day",
		},
	}
}

func initListForWeekHandlerTests(t *testing.T) []eventHandlerTest {
	t.Helper()
	events := []*app.EventDto{eventDto(t, userID2), eventDto2(t, userID2)}
	getEventsResponse, err := json.Marshal(EventsResponse{
		Events: events,
	})
	require.NoError(t, err)

	return []eventHandlerTest{
		{
			requestBody: []byte{},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "GET",
			url:    fmt.Sprintf("/events?%s=%s&%s=%s", startDateQueryKey, "2024-08-01", periodTypeQueryKey, periodWeekQueryValue),
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc("/events", handler.listForWeek).Methods("GET").Queries(
					startDateQueryKey, startDateQueryValue,
					periodTypeQueryKey, periodWeekQueryValue,
				)
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().ListForWeek(mock.Anything, userID2, getTime(t, "2024-08-01 00:00:00")).Return(events, nil)
			},
			expectedResponseBody: getEventsResponse,
			expectedResponseCode: http.StatusOK,
			testName:             "list events for week",
		},
	}
}

func initListForMonthHandlerTests(t *testing.T) []eventHandlerTest {
	t.Helper()
	events := []*app.EventDto{eventDto(t, userID2), eventDto2(t, userID2)}
	getEventsResponse, err := json.Marshal(EventsResponse{
		Events: events,
	})
	require.NoError(t, err)

	return []eventHandlerTest{
		{
			requestBody: []byte{},
			headers: map[string]string{
				userIDHeader: userID2Str,
			},
			method: "GET",
			url:    fmt.Sprintf("/events?%s=%s&%s=%s", startDateQueryKey, "2024-08-01", periodTypeQueryKey, periodMonthQueryValue),
			route: func(mux *mux.Router, handler *EventHandler) {
				mux.HandleFunc("/events", handler.listForMonth).Methods("GET").Queries(
					startDateQueryKey, startDateQueryValue,
					periodTypeQueryKey, periodMonthQueryValue,
				)
			},
			appCall: func(app *mocks.Application) {
				app.EXPECT().ListForMonth(mock.Anything, userID2, getTime(t, "2024-08-01 00:00:00")).Return(events, nil)
			},
			expectedResponseBody: getEventsResponse,
			expectedResponseCode: http.StatusOK,
			testName:             "list events for month",
		},
	}
}

func eventDto(t *testing.T, userID uint64) *app.EventDto {
	t.Helper()
	return &app.EventDto{
		ID:           1,
		Title:        "my event",
		StartDate:    getTime(t, "2024-07-06 10:00:00"),
		EndDate:      getTime(t, "2024-07-10 00:00:00"),
		Description:  "my event description",
		UserID:       userID,
		NotifyBefore: time.Hour * 24,
	}
}

func eventDto2(t *testing.T, userID uint64) *app.EventDto {
	t.Helper()
	return &app.EventDto{
		ID:           2,
		Title:        "my event2",
		StartDate:    getTime(t, "2024-08-06 10:00:00"),
		EndDate:      getTime(t, "2024-08-10 00:00:00"),
		Description:  "my event description2",
		UserID:       userID,
		NotifyBefore: time.Hour * 12,
	}
}

func getTime(t *testing.T, value string) time.Time {
	t.Helper()
	time, err := time.Parse(time.DateTime, value)
	require.NoError(t, err)
	return time
}
