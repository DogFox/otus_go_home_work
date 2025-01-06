package memorystorage

import (
	"context"
	"sync"

	domain "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	events map[string]domain.Event
	mu     sync.RWMutex //nolint:unused
}

func New() *Storage {
	return &Storage{
		events: map[string]domain.Event{},
		mu:     sync.RWMutex{},
	}
}
func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) CreateEvent(event domain.Event) error {
	s.events[event.ID] = event
	return nil
}
func (s *Storage) UpdateEvent(event domain.Event) error {
	s.events[event.ID] = event
	return nil
}
func (s *Storage) DeleteEvent(event domain.Event) error {
	delete(s.events, event.ID)
	return nil
}
func (s *Storage) EventList() []domain.Event {
	list := make([]domain.Event, 0, len(s.events))
	for _, v := range s.events {
		// fmt.Println("list ", v)
		list = append(list, v)
	}
	return list
}
