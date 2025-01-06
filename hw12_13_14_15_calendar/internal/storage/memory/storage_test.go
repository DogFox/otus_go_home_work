package memorystorage

import (
	"testing"

	domain "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require" //nolint: depguard
)

func TestStorage(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		storage := New()
		list := storage.EventList()

		require.Equal(t, 0, len(list))
	})

	t.Run("not empty list", func(t *testing.T) {
		storage := New()
		storage.CreateEvent(domain.Event{ID: "test", Title: "Test"})
		list := storage.EventList()

		require.Equal(t, 1, len(list))
	})
}
