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
	usersEvents map[uint64]map[uint64]struct{}
	mu          sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events:      make(map[uint64]*storage.Event),
		usersEvents: make(map[uint64]map[uint64]struct{}),
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
	eventsByUser, exists := s.usersEvents[event.UserID]
	if !exists {
		eventsByUser = make(map[uint64]struct{})
		s.usersEvents[event.UserID] = eventsByUser
	}
	eventsByUser[event.ID] = struct{}{}

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
	delete(s.usersEvents[userID], eventID)
	if len(s.usersEvents[userID]) == 0 {
		delete(s.usersEvents, userID)
	}

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
	eventsByUser, exists := s.usersEvents[userID]
	if !exists {
		return events, nil
	}

	for eventID := range eventsByUser {
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

func (s *Storage) ListForNotify(_ context.Context, startNotifyDate time.Time, endNotifyDate time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*storage.Event, 0)

	for _, event := range s.events {
		if event.NotifyBefore != 0 {
			notifyDate := event.StartDate.Add(-event.NotifyBefore)
			if notifyDate.Compare(startNotifyDate) >= 0 && notifyDate.Compare(endNotifyDate) <= 0 {
				events = append(events, event)
			}
		}
	}

	sort.Slice(events, func(i, j int) bool {
		notifyDate1 := events[i].StartDate.Add(events[i].NotifyBefore)
		notifyDate2 := events[j].StartDate.Add(events[j].NotifyBefore)
		if notifyDate1.Equal(notifyDate2) {
			if events[i].StartDate.Equal(events[j].StartDate) {
				return events[i].EndDate.Before(events[j].EndDate)
			}
			return events[i].StartDate.Before(events[j].StartDate)
		}
		return notifyDate1.Before(notifyDate2)
	})

	return events, nil
}

func (s *Storage) MarkAsNotified(_ context.Context, eventIDs []uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, eventID := range eventIDs {
		if event := s.events[eventID]; event != nil {
			event.Notified = true
		}
	}
	return nil
}

func (s *Storage) DeleteByEndDate(_ context.Context, maxEndDate time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.events {
		if event.EndDate.Compare(maxEndDate) <= 0 {
			delete(s.events, event.ID)
			delete(s.usersEvents[event.UserID], event.ID)
			if len(s.usersEvents[event.UserID]) == 0 {
				delete(s.usersEvents, event.UserID)
			}
		}
	}

	return nil
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

	for existingEventID := range eventsByUser {
		existingEvent := s.events[existingEventID]
		if event.ID != existingEvent.ID && event.StartDate.Compare(existingEvent.EndDate) <= 0 && event.EndDate.Compare(existingEvent.StartDate) >= 0 {
			return storage.ErrBusyTime
		}
	}

	return nil
}
