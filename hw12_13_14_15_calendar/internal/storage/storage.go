package storage

import "context"

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	CreateEvent() error
	UpdateEvent() error
	DeleteEvent() error
	EventList() []*Event
}
