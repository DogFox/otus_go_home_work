package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"

	pb "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/calendar/pb"
	"github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	domain "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/model"
	internalgrpc "github.com/DogFox/otus_go_home_work/hw12_13_14_15_calendar/internal/server/grpc"
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
	UserID:      12345,
	TimeShift:   15, // Уведомление за 15 минут до события
}

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		fmt.Println(configFile)
		fmt.Println(err)
		return
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

	httpServer := internalhttp.NewServer(logg, calendar, storage, config.Server.DSN())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(internalgrpc.UnaryLoggingInterceptor(logg)),
	)
	pb.RegisterEventsServer(grpcServer, internalgrpc.NewServer(logg, calendar, storage))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		startHTTPServer(ctx, httpServer, logg)
	}()

	go func() {
		defer wg.Done()
		startGRPCServer(ctx, grpcServer, logg)
	}()

	wg.Wait()
	logg.Info("calendar is running...")
}

func startHTTPServer(ctx context.Context, server *internalhttp.Server, logger *logger.Logger) error {
	go func() {
		<-ctx.Done()
		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
		logger.Error("server stopped")
	}()

	logger.Println("http server started: ", server.Addr)
	return server.Start(ctx)
}

func startGRPCServer(ctx context.Context, server *grpc.Server, logger *logger.Logger) error {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	go func() {
		<-ctx.Done()
		server.GracefulStop()
		logger.Error("server stopped")
	}()

	logger.Println("grpc server started :50051")
	return server.Serve(listener)
}
