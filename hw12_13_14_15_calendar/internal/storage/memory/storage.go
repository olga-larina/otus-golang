package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"golang.org/x/exp/rand"
)

type Storage struct {
	events map[int64]*storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[int64]*storage.Event),
	}
}

func (s *Storage) Create(_ context.Context, event *storage.Event) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existingEvent := range s.events {
		if event.StartDate.Compare(existingEvent.EndDate) <= 0 && event.EndDate.Compare(existingEvent.StartDate) >= 0 {
			return 0, storage.ErrBusyTime
		}
	}

	event.ID = s.generateUniqueID()

	s.events[event.ID] = event

	return event.ID, nil
}

func (s *Storage) GetByID(_ context.Context, eventID int64) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[eventID]
	if !exists {
		return nil, storage.ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) Update(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.events[event.ID]
	if !exists {
		return storage.ErrEventNotFound
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) Delete(_ context.Context, eventID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.events[eventID]
	if !exists {
		return storage.ErrEventNotFound
	}

	delete(s.events, eventID)

	return nil
}

func (s *Storage) ListForPeriod(
	_ context.Context,
	startDate time.Time,
	endDateExclusive time.Time,
) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*storage.Event, 0)
	for _, event := range s.events {
		if event.StartDate.Compare(endDateExclusive) < 0 && event.EndDate.Compare(startDate) >= 0 {
			events = append(events, event)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].StartDate.Equal(events[j].StartDate) {
			return events[i].EndDate.Before(events[j].EndDate)
		}
		return events[i].StartDate.Before(events[j].StartDate)
	})

	return events, nil
}

func (s *Storage) generateUniqueID() int64 {
	var eventID int64
	var exists bool

	for exists = true; exists || eventID == 0; {
		eventID = rand.Int63()
		_, exists = s.events[eventID]
	}

	return eventID
}
