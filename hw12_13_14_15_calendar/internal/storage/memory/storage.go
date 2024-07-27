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
	events      map[uint64]*storage.Event
	usersEvents map[uint64][]uint64
	mu          sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events:      make(map[uint64]*storage.Event),
		usersEvents: make(map[uint64][]uint64),
	}
}

func (s *Storage) Create(_ context.Context, event *storage.Event) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.generateUniqueID()

	if err := s.checkBusyTime(event); err != nil {
		return 0, err
	}

	s.events[event.ID] = event
	s.usersEvents[event.UserID] = append(s.usersEvents[event.UserID], event.ID)

	return event.ID, nil
}

func (s *Storage) GetByID(_ context.Context, userID uint64, eventID uint64) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[eventID]
	if !exists || event.UserID != userID {
		return nil, storage.ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) Update(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingEvent, exists := s.events[event.ID]
	if !exists || event.UserID != existingEvent.UserID {
		return storage.ErrEventNotFound
	}

	if err := s.checkBusyTime(event); err != nil {
		return err
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) Delete(_ context.Context, userID uint64, eventID uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingEvent, exists := s.events[eventID]
	if !exists || existingEvent.UserID != userID {
		return storage.ErrEventNotFound
	}

	delete(s.events, eventID)

	idx := 0
	eventsByUser := s.usersEvents[userID]
	for _, existingEventID := range eventsByUser {
		if existingEventID != eventID {
			eventsByUser[idx] = existingEventID
			idx++
		}
	}
	s.usersEvents[userID] = eventsByUser[:idx]

	return nil
}

func (s *Storage) ListForPeriod(
	_ context.Context,
	userID uint64,
	startDate time.Time,
	endDateExclusive time.Time,
) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*storage.Event, 0)
	eventsByUser := s.usersEvents[userID]

	for _, eventID := range eventsByUser {
		event := s.events[eventID]
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

func (s *Storage) generateUniqueID() uint64 {
	var eventID uint64
	var exists bool

	for exists = true; exists || eventID == 0; {
		eventID = rand.Uint64()
		_, exists = s.events[eventID]
	}

	return eventID
}

func (s *Storage) checkBusyTime(event *storage.Event) error {
	eventsByUser := s.usersEvents[event.UserID]

	for _, existingEventID := range eventsByUser {
		existingEvent := s.events[existingEventID]
		if event.ID != existingEvent.ID && event.StartDate.Compare(existingEvent.EndDate) <= 0 && event.EndDate.Compare(existingEvent.StartDate) >= 0 {
			return storage.ErrBusyTime
		}
	}

	return nil
}
