package app

import (
	"context"

	domain "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/model"
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
	CreateEvent(e domain.Event) error
	UpdateEvent(e domain.Event) error
	DeleteEvent(e domain.Event) error
	EventList() []domain.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(domain.Event{ID: id, Title: title})
}

func (a *App) UpdateEvent(e domain.Event) error {
	a.storage.UpdateEvent(e)
	return nil
}

func (a *App) DeleteEvent(e domain.Event) error {
	a.storage.DeleteEvent(e)
	return nil
}

func (a *App) EventList() []domain.Event {
	return a.storage.EventList()
}
