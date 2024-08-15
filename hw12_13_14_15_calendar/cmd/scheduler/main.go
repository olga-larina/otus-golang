package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed reading config %v", err)
		return
	}

	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Fatalf("failed loading location %v", err)
		return
	}
	time.Local = location

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatalf("failed building logger %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// storage
	sqlStorage := sqlstorage.New(config.Database.Driver, config.Database.URI)
	if err := sqlStorage.Connect(ctx); err != nil {
		logg.Error(ctx, err, "failed to connect to db")
		return
	}
	defer func() {
		if err := sqlStorage.Close(ctx); err != nil {
			logg.Error(ctx, err, "failed to close sql storage")
		}
	}()

	// queue
	queue := rabbit.NewQueue(
		logg,
		config.Queue.URI,
		config.Queue.ExchangeName,
		config.Queue.ExchangeType,
		config.Queue.QueueName,
		config.Queue.RoutingKey,
	)
	publisher := queue.NewPublisher()
	if err := queue.Start(ctx); err != nil {
		logg.Error(ctx, err, "queue failed to start")
		return
	}

	// scheduler
	scheduler := scheduler.NewScheduler(
		ctx,
		logg,
		sqlStorage,
		publisher,
		config.Schedule.NotifyCron,
		config.Schedule.ClearCron,
		config.Schedule.NotifyPeriod,
		config.Schedule.NotifyScanPeriod,
		config.Schedule.ClearPeriod,
	)

	if err := scheduler.Start(ctx); err != nil {
		logg.Error(ctx, err, "scheduler failed to start")
		return
	}

	logg.Info(ctx, "scheduler is running...")

	<-ctx.Done()

	if err := queue.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop queue")
	}

	if err := scheduler.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop scheduler")
	}
}
