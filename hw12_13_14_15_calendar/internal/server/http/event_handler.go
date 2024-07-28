package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

const (
	userIDHeader = "X-USER-ID"
	eventIDPath  = "eventID"
)

const (
	startDateQueryKey   = "startDate"
	startDateQueryValue = `{startDate:\d{4}-\d{2}-\d{2}}`
)

const (
	periodTypeQueryKey    = "period"
	periodDayQueryValue   = "day"
	periodWeekQueryValue  = "week"
	periodMonthQueryValue = "month"
)

var (
	errNotValidUserID    = errors.New("userID is not valid")
	errNotValidEventID   = errors.New("eventID is not valid")
	errNotValidStartDate = errors.New("startDate is not valid")
)

type EventHandler struct {
	logger Logger
	app    Application
}

func NewEventHandler(logger Logger, app Application) *EventHandler {
	return &EventHandler{
		logger: logger,
		app:    app,
	}
}

func (s *EventHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("failed reading request: %w", err).Error(), http.StatusBadRequest)
		return
	}

	var createEventReq CreateEventRequest
	if err := json.Unmarshal(reqData, &createEventReq); err != nil {
		http.Error(w, fmt.Errorf("failed parsing request: %w", err).Error(), http.StatusBadRequest)
		return
	}

	createEventReq.Event.UserID = userID

	eventID, err := s.app.Create(ctx, *createEventReq.Event)
	if err != nil {
		if errors.Is(err, storage.ErrBusyTime) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateEventResponse{EventID: eventID}
	s.writeResponse(ctx, w, response)
}

func (s *EventHandler) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventID, err := getEventID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := s.app.GetByID(ctx, userID, eventID)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EventResponse{Event: event}
	s.writeResponse(ctx, w, response)
}

func (s *EventHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventID, err := getEventID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("failed reading request: %w", err).Error(), http.StatusBadRequest)
		return
	}

	var updateEventReq UpdateEventRequest
	if err := json.Unmarshal(reqData, &updateEventReq); err != nil {
		http.Error(w, fmt.Errorf("failed parsing request: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if updateEventReq.Event.ID == 0 || updateEventReq.Event.ID != eventID {
		http.Error(w, errNotValidEventID.Error(), http.StatusBadRequest)
		return
	}

	updateEventReq.Event.UserID = userID

	err = s.app.Update(ctx, *updateEventReq.Event)
	if err != nil {
		if errors.Is(err, storage.ErrBusyTime) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *EventHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventID, err := getEventID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Delete(ctx, userID, eventID)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *EventHandler) listForDay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := getStartDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.ListForDay(ctx, userID, startDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EventsResponse{Events: events}
	s.writeResponse(ctx, w, response)
}

func (s *EventHandler) listForWeek(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := getStartDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.ListForWeek(ctx, userID, startDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EventsResponse{Events: events}
	s.writeResponse(ctx, w, response)
}

func (s *EventHandler) listForMonth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := getStartDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.ListForMonth(ctx, userID, startDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EventsResponse{Events: events}
	s.writeResponse(ctx, w, response)
}

func getUserID(r *http.Request) (uint64, error) {
	var userID uint64
	var err error

	userIDStr := r.Header.Get(userIDHeader)
	if len(userIDStr) > 0 {
		userID, err = strconv.ParseUint(userIDStr, 10, 64)
	} else {
		err = errNotValidUserID
	}

	if err == nil && userID == 0 {
		err = errNotValidUserID
	}

	return userID, err
}

func getEventID(r *http.Request) (uint64, error) {
	var eventID uint64
	var err error

	vars := mux.Vars(r)
	eventIDStr := vars[eventIDPath]
	if len(eventIDStr) > 0 {
		eventID, err = strconv.ParseUint(eventIDStr, 10, 64)
	} else {
		err = errNotValidEventID
	}

	if err == nil && eventID == 0 {
		err = errNotValidEventID
	}

	return eventID, err
}

func getStartDate(r *http.Request) (time.Time, error) {
	startDateStr := r.URL.Query().Get(startDateQueryKey)
	if len(startDateStr) == 0 {
		return time.Time{}, errNotValidStartDate
	}

	startDate, err := time.Parse(time.DateOnly, startDateStr)
	if err != nil {
		return time.Time{}, err
	}

	return startDate, nil
}

func (s *EventHandler) writeResponse(ctx context.Context, w http.ResponseWriter, resp any) {
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Errorf("failed encoding response: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		s.logger.Error(ctx, err, "error writing response")
	}
}
