package app

import (
	"context"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Storage interface {
	Create(ctx context.Context, event *storage.Event) (uint64, error)
	GetByID(ctx context.Context, userID uint64, eventID uint64) (*storage.Event, error)
	Update(ctx context.Context, event *storage.Event) error
	Delete(ctx context.Context, userID uint64, eventID uint64) error
	ListForPeriod(ctx context.Context, userID uint64, startDate time.Time, endDateExclusive time.Time) ([]*storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Create(ctx context.Context, eventDto EventDto) (uint64, error) {
	return a.storage.Create(ctx, convertEventToModel(&eventDto))
}

func (a *App) GetByID(ctx context.Context, userID uint64, eventID uint64) (*EventDto, error) {
	event, err := a.storage.GetByID(ctx, userID, eventID)
	if err != nil {
		return nil, err
	}
	return convertEventToDto(event), nil
}

func (a *App) Update(ctx context.Context, eventDto EventDto) error {
	return a.storage.Update(ctx, convertEventToModel(&eventDto))
}

func (a *App) Delete(ctx context.Context, userID uint64, eventID uint64) error {
	return a.storage.Delete(ctx, userID, eventID)
}

func (a *App) ListForDay(ctx context.Context, userID uint64, date time.Time) ([]*EventDto, error) {
	endDateExclusive := date.Add(24 * time.Hour)
	return a.listForPeriod(ctx, userID, date, endDateExclusive)
}

func (a *App) ListForWeek(ctx context.Context, userID uint64, startDate time.Time) ([]*EventDto, error) {
	endDateExclusive := startDate.AddDate(0, 0, 7)
	return a.listForPeriod(ctx, userID, startDate, endDateExclusive)
}

func (a *App) ListForMonth(ctx context.Context, userID uint64, startDate time.Time) ([]*EventDto, error) {
	endDateExclusive := startDate.AddDate(0, 1, 0)
	return a.listForPeriod(ctx, userID, startDate, endDateExclusive)
}

func (a *App) listForPeriod(ctx context.Context, userID uint64, startDate time.Time, endDateExclusive time.Time) ([]*EventDto, error) {
	events, err := a.storage.ListForPeriod(ctx, userID, startDate, endDateExclusive)
	if err != nil {
		return nil, err
	}
	eventsDto := make([]*EventDto, len(events))
	for i, event := range events {
		eventsDto[i] = convertEventToDto(event)
	}
	return eventsDto, nil
}
