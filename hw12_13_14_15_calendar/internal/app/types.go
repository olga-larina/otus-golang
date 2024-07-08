package app

import (
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type EventDto struct {
	ID           int64
	Title        string
	StartDate    time.Time
	EndDate      time.Time
	Description  string
	UserID       int64
	NotifyBefore time.Duration
}

func convertEventToModel(dto *EventDto) *storage.Event {
	return &storage.Event{
		ID:           dto.ID,
		Title:        dto.Title,
		StartDate:    dto.StartDate,
		EndDate:      dto.EndDate,
		Description:  dto.Description,
		UserID:       dto.UserID,
		NotifyBefore: dto.NotifyBefore,
	}
}

func convertEventToDto(model *storage.Event) *EventDto {
	return &EventDto{
		ID:           model.ID,
		Title:        model.Title,
		StartDate:    model.StartDate,
		EndDate:      model.EndDate,
		Description:  model.Description,
		UserID:       model.UserID,
		NotifyBefore: model.NotifyBefore,
	}
}
