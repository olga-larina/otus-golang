package app

import (
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type EventDto struct {
	ID           uint64        `json:"id"`
	Title        string        `json:"title"`
	StartDate    time.Time     `json:"startDate"`
	EndDate      time.Time     `json:"endDate"`
	Description  string        `json:"description"`
	UserID       uint64        `json:"userId"`
	NotifyBefore time.Duration `json:"notifyBefore"`
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
