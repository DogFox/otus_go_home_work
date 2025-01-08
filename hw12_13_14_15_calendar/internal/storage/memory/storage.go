package memorystorage

import (
	"context"
	"strconv"
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

func (s *Storage) CreateEvent(ctx context.Context, event domain.Event) error {
	s.mu.Lock()
	event.ID = int64(len(s.events) + 1)
	s.events[strconv.Itoa(int(event.ID))] = event
	s.mu.Unlock()
	return nil
}
func (s *Storage) UpdateEvent(ctx context.Context, event domain.Event) error {
	s.mu.Lock()
	s.events[strconv.Itoa(int(event.ID))] = event
	s.mu.Unlock()
	return nil
}
func (s *Storage) DeleteEvent(ctx context.Context, event domain.Event) error {
	delete(s.events, strconv.Itoa(int(event.ID)))
	return nil
}
func (s *Storage) EventList(ctx context.Context) ([]domain.Event, error) {
	list := make([]domain.Event, 0, len(s.events))
	s.mu.Lock()
	for _, v := range s.events {
		// fmt.Println("list ", v)
		list = append(list, v)
	}
	s.mu.Unlock()
	return list, nil
}
