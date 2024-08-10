package model

import (
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type NotificationDto struct {
	EventID        uint64    `json:"eventId"`
	EventTitle     string    `json:"eventTitle"`
	EventStartDate time.Time `json:"eventStartDate"`
	UserID         uint64    `json:"userId"`
}

func ConvertEventToNotification(event *storage.Event) *NotificationDto {
	return &NotificationDto{
		EventID:        event.ID,
		EventTitle:     event.Title,
		EventStartDate: event.StartDate,
		UserID:         event.UserID,
	}
}
