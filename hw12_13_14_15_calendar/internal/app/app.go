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
	Create(ctx context.Context, event *storage.Event) (int64, error)
	GetByID(ctx context.Context, eventID int64) (*storage.Event, error)
	Update(ctx context.Context, event *storage.Event) error
	Delete(ctx context.Context, eventID int64) error
	ListForPeriod(ctx context.Context, startDate time.Time, endDateExclusive time.Time) ([]*storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, eventDto EventDto) (int64, error) {
	return a.storage.Create(ctx, convertEventToModel(&eventDto))
}

func (a *App) GetByID(ctx context.Context, eventID int64) (*EventDto, error) {
	event, err := a.storage.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return convertEventToDto(event), nil
}

func (a *App) Update(ctx context.Context, eventDto EventDto) error {
	return a.storage.Update(ctx, convertEventToModel(&eventDto))
}

func (a *App) Delete(ctx context.Context, eventID int64) error {
	return a.storage.Delete(ctx, eventID)
}

func (a *App) ListForDay(ctx context.Context, date time.Time) ([]*EventDto, error) {
	endDateExclusive := date.Add(24 * time.Hour)
	return a.listForPeriod(ctx, date, endDateExclusive)
}

func (a *App) ListForWeek(ctx context.Context, startDate time.Time) ([]*EventDto, error) {
	endDateExclusive := startDate.AddDate(0, 0, 7)
	return a.listForPeriod(ctx, startDate, endDateExclusive)
}

func (a *App) ListForMonth(ctx context.Context, startDate time.Time) ([]*EventDto, error) {
	endDateExclusive := startDate.AddDate(0, 1, 0)
	return a.listForPeriod(ctx, startDate, endDateExclusive)
}

func (a *App) listForPeriod(ctx context.Context, startDate time.Time, endDateExclusive time.Time) ([]*EventDto, error) {
	events, err := a.storage.ListForPeriod(ctx, startDate, endDateExclusive)
	if err != nil {
		return nil, err
	}
	eventsDto := make([]*EventDto, 0, len(events))
	for i, event := range events {
		eventsDto[i] = convertEventToDto(event)
	}
	return eventsDto, nil
}
