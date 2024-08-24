package health

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/logger"
)

const (
	readyFileDefault       = "/tmp/ready"
	heartbeatFileDefault   = "/tmp/health"
	heartbeatPeriodDefault = 5 * time.Second
)

func FileHealthcheck(ctx context.Context, log *logger.Logger) error {
	var err error

	readyFile, found := os.LookupEnv("READY_FILE")
	if !found {
		readyFile = readyFileDefault
	}

	heartbeatFile, found := os.LookupEnv("HEARTBEAT_FILE")
	if !found {
		heartbeatFile = heartbeatFileDefault
	}

	var heartbeatPeriod time.Duration

	heartbeatPeriodStr, found := os.LookupEnv("HEARTBEAT_PERIOD")
	if found {
		heartbeatPeriod, err = time.ParseDuration(heartbeatPeriodStr)
		if err != nil {
			found = false
		}
	}
	if !found {
		heartbeatPeriod = heartbeatPeriodDefault
	}

	err = createReadyFile(readyFile)
	if err != nil {
		return err
	}

	updateHeartbeatFile(ctx, heartbeatFile, heartbeatPeriod, log)

	return nil
}

// создаёт файл для сигнализации о готовности приложения.
func createReadyFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// обновляет файл с текущим временем.
func updateHeartbeatFile(ctx context.Context, file string, period time.Duration, log *logger.Logger) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// #nosec G306
			err := os.WriteFile(file, []byte(fmt.Sprintf("%d", time.Now().Unix())), 0o644)
			if err != nil {
				log.Error(ctx, err, "failed updating health file")
			}
		}
	}
}
