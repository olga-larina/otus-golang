package sender

import (
	"context"
	"encoding/json"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/model"
)

type Sender struct {
	logger Logger
	queue  Queue
	done   chan struct{}
	ctx    context.Context
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Queue interface {
	ReceiveData(ctx context.Context) (<-chan []byte, error)
}

func NewSender(
	ctx context.Context,
	logger Logger,
	queue Queue,
) *Sender {
	return &Sender{
		logger: logger,
		queue:  queue,
		done:   make(chan struct{}),
		ctx:    ctx,
	}
}

func (s *Sender) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting sender")

	err := s.processEvents(ctx)
	if err != nil {
		return err
	}

	s.logger.Info(ctx, "started sender")
	return nil
}

func (s *Sender) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "stopping sender")

	<-ctx.Done()
	<-s.done

	s.logger.Info(ctx, "stopped sender")
	return nil
}

/*
 * Получение событий из очереди.
 */
func (s *Sender) processEvents(ctx context.Context) error {
	data, err := s.queue.ReceiveData(ctx)
	if err != nil {
		defer close(s.done)
		return err
	}

	go func() {
		defer close(s.done)

		for d := range data {
			var notification model.NotificationDto

			err := json.Unmarshal(d, &notification)
			if err != nil {
				s.logger.Error(ctx, err, "failed to read notification")
			} else {
				s.logger.Info(ctx, "received notification", "notification", notification)
			}
		}
	}()

	return nil
}
