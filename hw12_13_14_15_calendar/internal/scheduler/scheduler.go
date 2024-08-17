package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/model"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	logger           Logger
	storage          Storage
	queue            Queue
	cron             *cron.Cron
	notifyCron       string
	clearCron        string
	notifyPeriod     time.Duration
	notifyScanPeriod time.Duration
	clearPeriod      time.Duration
	ctx              context.Context
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Storage interface {
	ListForNotify(ctx context.Context, startNotifyDate time.Time, endNotifyDate time.Time) ([]*storage.Event, error)
	SetNotifyStatus(ctx context.Context, eventIDs []uint64, notifyStatus storage.NotifyStatus) error
	DeleteByEndDate(ctx context.Context, maxEndDate time.Time) error
}

type Queue interface {
	SendData(ctx context.Context, data []byte) error
}

func NewScheduler(
	ctx context.Context,
	logger Logger,
	storage Storage,
	queue Queue,
	notifyCron string,
	clearCron string,
	notifyPeriod time.Duration,
	notifyScanPeriod time.Duration,
	clearPeriod time.Duration,
) *Scheduler {
	return &Scheduler{
		logger:           logger,
		storage:          storage,
		queue:            queue,
		cron:             cron.New(cron.WithSeconds()),
		notifyCron:       notifyCron,
		clearCron:        clearCron,
		notifyPeriod:     notifyPeriod,
		notifyScanPeriod: notifyScanPeriod,
		clearPeriod:      clearPeriod,
		ctx:              ctx,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting scheduler")

	s.cron.AddFunc(s.notifyCron, func() {
		s.notifyEvents()
	})
	s.cron.AddFunc(s.clearCron, func() {
		s.clearEvents()
	})
	s.cron.Start()

	s.logger.Info(ctx, "started scheduler")
	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "stopping scheduler")

	ctx = s.cron.Stop()
	<-ctx.Done()

	s.logger.Info(ctx, "stopped scheduler")
	return nil
}

/*
 * Уведомление о предстоящих событиях.
 */
func (s *Scheduler) notifyEvents() {
	now := time.Now()
	startNotifyDate := now.Add(-s.notifyScanPeriod)
	endNotifyDate := now.Add(s.notifyPeriod)
	s.logger.Debug(s.ctx, "start notifying events", "startNotifyDate", startNotifyDate, "endNotifyDate", endNotifyDate)

	events, err := s.storage.ListForNotify(s.ctx, startNotifyDate, endNotifyDate)
	if err != nil {
		s.logger.Error(
			s.ctx, err, "failed notifying events",
			"stage", "storage",
			"startNotifyDate", startNotifyDate,
			"endNotifyDate", endNotifyDate,
		)
		return
	}

	eventIDs := make([]uint64, 0)
	for _, event := range events {
		notification := model.ConvertEventToNotification(event)
		notificationStr, err := json.Marshal(notification)
		if err != nil {
			s.logger.Error(
				s.ctx, err, "failed notifying events",
				"stage", "marshal",
				"eventId", event.ID,
				"startNotifyDate", startNotifyDate,
				"endNotifyDate", endNotifyDate,
			)
			continue
		}

		err = s.queue.SendData(s.ctx, notificationStr)
		if err != nil {
			s.logger.Error(
				s.ctx, err, "failed notifying events",
				"stage", "send",
				"eventId", event.ID,
				"startNotifyDate", startNotifyDate,
				"endNotifyDate", endNotifyDate,
			)
			continue
		}
		s.logger.Debug(s.ctx, "notification sent", "eventID", event.ID)

		eventIDs = append(eventIDs, event.ID)
	}

	err = s.storage.SetNotifyStatus(s.ctx, eventIDs, storage.NotifyInProgress)
	if err != nil {
		s.logger.Error(
			s.ctx, err, "failed notifying events",
			"stage", "marking",
			"startNotifyDate", startNotifyDate,
			"endNotifyDate", endNotifyDate,
		)
	} else {
		s.logger.Debug(s.ctx, "succeeded notifying events")
	}
}

/*
 * Очистка старых событий.
 */
func (s *Scheduler) clearEvents() {
	maxEndDateToDelete := time.Now().Add(-s.clearPeriod)
	s.logger.Debug(s.ctx, "start clearing events", "maxEndDate", maxEndDateToDelete)

	err := s.storage.DeleteByEndDate(s.ctx, maxEndDateToDelete)
	if err != nil {
		s.logger.Error(s.ctx, err, "failed clearing events", "maxEndDate", maxEndDateToDelete)
	} else {
		s.logger.Debug(s.ctx, "succeeded clearing events", "maxEndDate", maxEndDateToDelete)
	}
}
