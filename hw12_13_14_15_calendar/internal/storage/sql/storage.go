package sqlstorage

import (
	"context"
	"fmt"
	"log"

	domain "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	conn *pgx.Conn
	dsn  string
}

func New(dsn string) *Storage {
	return &Storage{
		dsn:  dsn,
		conn: nil,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.conn.Close(ctx)
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event domain.Event) error {
	sql := `INSERT INTO events (user_id, title, date, duration, timeshift, description) 
			VALUES ($1, $2, $3, $4, $5, $6)`
	fmt.Println(event)
	_, err := s.conn.Exec(ctx, sql, event.User_ID, event.Title, event.Date, event.Duration, event.TimeShift, event.Description)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (s *Storage) UpdateEvent(ctx context.Context, event domain.Event) error {
	sql := `UPDATE events 
	SET user_id = $1, title = $2, date = $3, duration = $4, timeshift = $5, description = $6
	WHERE id = $7`

	_, err := s.conn.Exec(ctx, sql, event.User_ID, event.Title, event.Date, event.Duration, event.TimeShift, event.Description, event.ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) DeleteEvent(ctx context.Context, event domain.Event) error {
	sql := `DELETE FROM events WHERE id = $1`

	_, err := s.conn.Exec(ctx, sql, event.ID)
	if err != nil {
		log.Fatal("failed to delete event:", err)
	}
	return nil
}

func (s *Storage) EventList(ctx context.Context) ([]domain.Event, error) {
	list := make([]domain.Event, 0)
	rows, err := s.conn.Query(ctx, "SELECT id, user_id, title, date, duration, timeshift, description FROM events")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var event domain.Event
		err := rows.Scan(&event.ID, &event.User_ID, &event.Title, &event.Date, &event.Duration, &event.TimeShift, &event.Description)
		if err != nil {
			return list, err
		}
		list = append(list, event)
	}
	return list, nil
}
