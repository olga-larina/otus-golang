package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/health"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/sender"
	sqlstorage "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/sender/config.yaml", "Path to configuration file")
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
	consumer := queue.NewConsumer(config.Queue.ConsumerTag)
	if err := queue.Start(ctx); err != nil {
		logg.Error(ctx, err, "queue failed to start")
		return
	}

	// sender
	sender := sender.NewSender(
		ctx,
		logg,
		sqlStorage,
		consumer,
	)

	if err := sender.Start(ctx); err != nil {
		logg.Error(ctx, err, "sender failed to start")
		return
	}

	if err = health.FileHealthcheck(ctx, logg); err != nil {
		logg.Error(ctx, err, "healthcheck failed to start")
		return
	}

	logg.Info(ctx, "sender is running...")

	<-ctx.Done()

	if err := consumer.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop consumer")
	}

	if err := queue.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop queue")
	}

	if err := sender.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop sender")
	}
}
