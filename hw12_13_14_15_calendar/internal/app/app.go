package app

import (
	"context"

	"github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	CreateEvent(e storage.Event) error
	UpdateEvent() error
	DeleteEvent() error
	EventList() []*storage.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) UpdateEvent() error {
	return nil
}

func (a *App) DeleteEvent() error {
	return nil
}

// func (a *App) EventList() []*storage.Event {
// 	return nil
// }
