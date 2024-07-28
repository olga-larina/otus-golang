package internalhttp

import "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"

type CreateEventRequest struct {
	Event *app.EventDto `json:"event"`
}

type CreateEventResponse struct {
	EventID uint64 `json:"eventId"`
}

type UpdateEventRequest struct {
	Event *app.EventDto `json:"event"`
}

type EventResponse struct {
	Event *app.EventDto `json:"event"`
}

type EventsResponse struct {
	Events []*app.EventDto `json:"events"`
}
