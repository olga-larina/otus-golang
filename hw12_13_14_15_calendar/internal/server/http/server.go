package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	logger      Logger
	app         Application
	addr        string
	readTimeout time.Duration
	srv         *http.Server
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, eventDto app.EventDto) (int64, error)
	GetByID(ctx context.Context, eventID int64) (*app.EventDto, error)
	Update(ctx context.Context, eventDto app.EventDto) error
	Delete(ctx context.Context, eventID int64) error
	ListForDay(ctx context.Context, date time.Time) ([]*app.EventDto, error)
	ListForWeek(ctx context.Context, startDate time.Time) ([]*app.EventDto, error)
	ListForMonth(ctx context.Context, startDate time.Time) ([]*app.EventDto, error)
}

func NewServer(logger Logger, app Application, addr string, readTimeout time.Duration) *Server {
	return &Server{
		logger:      logger,
		app:         app,
		addr:        addr,
		readTimeout: readTimeout,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting http server")

	mux := mux.NewRouter()
	mux.Handle("/hello", loggingMiddleware(ctx, s.logger, http.HandlerFunc(helloHandler))).Methods("GET")

	s.srv = &http.Server{
		Addr:        s.addr,
		Handler:     mux,
		ReadTimeout: s.readTimeout,
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "stopping http server")
	return s.srv.Shutdown(ctx)
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "hello world")
}
