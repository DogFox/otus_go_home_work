package app

import (
	"context"

	domain "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/model"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface{}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, e domain.Event) error
	UpdateEvent(ctx context.Context, e domain.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	EventList(ctx context.Context, date string, listType string) ([]domain.Event, error)
	ClearEvents(ctx context.Context, life string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event domain.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, e domain.Event) error {
	a.storage.UpdateEvent(ctx, e)
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	a.storage.DeleteEvent(ctx, id)
	return nil
}

func (a *App) EventList(ctx context.Context, date string, listType string) ([]domain.Event, error) {
	return a.storage.EventList(ctx, date, listType)
}
