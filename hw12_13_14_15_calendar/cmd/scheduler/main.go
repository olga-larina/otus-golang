package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

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

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatalf("failed building logger %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// storage
	dbDsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		config.Database.DsnPrefix,
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
	)
	sqlStorage := sqlstorage.New(config.Database.Driver, dbDsn)
	if err := sqlStorage.Connect(ctx); err != nil {
		logg.Error(ctx, err, "failed to connect to db")
		return
	}
	defer sqlStorage.Close(ctx)

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
