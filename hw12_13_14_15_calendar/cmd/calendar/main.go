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
	internalhttp "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

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

	storage := memorystorage.New()

	sqlstorage := sqlstorage.New()
	err = sqlstorage.Connect(ctx, config.Database.DSN())
	if err != nil {
		logg.Error(err)
	}

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar)

	calendar.CreateEvent(ctx, "test", "test")
	calendar.CreateEvent(ctx, "test2", "test2")
	calendar.CreateEvent(ctx, "test3", "test3")
	calendar.CreateEvent(ctx, "test4", "test4")

	fmt.Println(calendar.EventList())

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
