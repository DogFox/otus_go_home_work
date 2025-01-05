package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu sync.RWMutex //nolint:unused
}

func New() *Storage {
	return &Storage{}
}
func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) CreateEvent(event storage.Event) error {
	fmt.Println("event ", event)
	return nil
}
func (s *Storage) UpdateEvent() error {
	return nil
}
func (s *Storage) DeleteEvent() error {
	return nil
}
func (s *Storage) EventList() []*storage.Event {
	return make([]*storage.Event, 0)
}
