package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/sender"
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

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatalf("failed building logger %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

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
	}

	// sender
	sender := sender.NewSender(
		ctx,
		logg,
		consumer,
	)

	if err := sender.Start(ctx); err != nil {
		logg.Error(ctx, err, "sender failed to start")
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
