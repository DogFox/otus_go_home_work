package memorystorage

import (
	"context"
	"sync"
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

func (s *Storage) CreateEvent() {

}
func (s *Storage) UpdateEvent() {

}
func (s *Storage) DeleteEvent() {

}
func (s *Storage) EventList() {

}
