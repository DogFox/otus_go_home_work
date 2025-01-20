package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	domain "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/model"
	internalhttp "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

var testEvent = domain.Event{
	ID:          1,
	Title:       "Morning Jog",
	Date:        time.Date(2025, time.January, 8, 6, 0, 0, 0, time.UTC), // 8 Jan 2025, 06:00 UTC
	Duration:    time.Hour * 1,                                          // 1 час
	Description: "A refreshing morning jog through the park.",
	User_ID:     12345,
	TimeShift:   15, // Уведомление за 15 минут до события
}

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig()
	if err != nil {
		fmt.Println(err)
	}
	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage
	switch config.Storage.Type {
	case "postgres":
		storage = sqlstorage.New(config.Database.DSN())
		err = storage.Connect(ctx)
		if err != nil {
			logg.Error(err)
		}
	case "memory":
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, config.Server.DSN())

	calendar.CreateEvent(ctx, testEvent)
	calendar.CreateEvent(ctx, testEvent)
	calendar.CreateEvent(ctx, testEvent)
	calendar.CreateEvent(ctx, testEvent)
	calendar.CreateEvent(ctx, testEvent)

	fmt.Println(calendar.EventList(ctx))

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	// logg.Info("Это INFO сообщение")
	// logg.Warn("Это WARN сообщение")
	// logg.Error("Это ERROR сообщение")
	// logg.Fatal("Это FATAL сообщение")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
