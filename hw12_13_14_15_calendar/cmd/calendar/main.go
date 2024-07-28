package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

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
	var storage app.Storage
	if config.StorageType == "SQL" {
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
		storage = sqlStorage
	} else {
		storage = memorystorage.New()
	}

	// app
	calendar := app.New(logg, storage)

	// http server
	httpServerAddr := fmt.Sprintf("%s:%s", config.HTTPServer.Host, config.HTTPServer.Port)
	httpServer := internalhttp.NewServer(logg, calendar, httpServerAddr, config.HTTPServer.ReadTimeout)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logg.Error(ctx, err, "failed to stop http server")
		}
	}()

	// grpc server
	grpcServer := internalgrpc.NewServer(logg, calendar, config.GrpcServer.Port)

	go func() {
		if err = grpcServer.Start(ctx); err != nil {
			logg.Error(ctx, err, "grpc failed to serve")
		}
	}()

	logg.Info(ctx, "calendar is running...")

	if err := httpServer.Start(ctx); err != nil {
		logg.Error(ctx, err, "http server stopped")
	}

	<-ctx.Done()

	if err := grpcServer.Stop(ctx); err != nil {
		logg.Error(ctx, err, "failed to stop grpc server")
	}
}
