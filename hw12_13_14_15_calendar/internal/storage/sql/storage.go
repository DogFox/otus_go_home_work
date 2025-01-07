package sqlstorage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Storage struct {
	conn *pgx.Conn
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	s.conn = conn
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
