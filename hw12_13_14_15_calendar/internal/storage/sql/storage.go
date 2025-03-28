package sqlstorage

import (
	"context"
	"fmt"
	"log"
	"strings"

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
	fmt.Println(s.dsn)
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
	_, err := s.conn.Exec(
		ctx,
		sql,
		event.UserID,
		event.Title,
		event.Date,
		event.Duration,
		event.TimeShift,
		event.Description,
	)
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

	_, err := s.conn.Exec(
		ctx,
		sql,
		event.UserID,
		event.Title,
		event.Date,
		event.Duration,
		event.TimeShift,
		event.Description,
		event.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	sql := `DELETE FROM events WHERE id = $1`

	_, err := s.conn.Exec(ctx, sql, id)
	if err != nil {
		log.Fatal("failed to delete event:", err)
	}
	return nil
}

func (s *Storage) EventList(ctx context.Context, date string, listType string) ([]domain.Event, error) {
	list := make([]domain.Event, 0)

	sqlQuery := "SELECT id, user_id, title, date, duration, timeshift, description FROM events "
	switch strings.ToLower(listType) {
	case "day":
		sqlQuery += " WHERE date = $1"
	case "week":
		sqlQuery += " WHERE date >= $1 AND date < $1 + INTERVAL '7 days'"
	case "month":
		sqlQuery += " WHERE date >= $1 AND date < $1 + INTERVAL '1 month'"
	}
	rows, err := s.conn.Query(ctx, sqlQuery, date)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var event domain.Event
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.Date,
			&event.Duration,
			&event.TimeShift,
			&event.Description,
		)
		if err != nil {
			return list, err
		}
		list = append(list, event)
	}
	return list, nil
}

func (s *Storage) ClearEvents(ctx context.Context, life string) error {
	sql := fmt.Sprintf("DELETE FROM events WHERE date < NOW() - INTERVAL '%s'", life)

	_, err := s.conn.Exec(ctx, sql)
	if err != nil {
		log.Fatal("failed to delete event:", err)
	}
	return nil
}
