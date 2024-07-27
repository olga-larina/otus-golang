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
	addr        string
	readTimeout time.Duration
	handler     *EventHandler
	srv         *http.Server
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

type Application interface {
	Create(ctx context.Context, eventDto app.EventDto) (uint64, error)
	GetByID(ctx context.Context, userID uint64, eventID uint64) (*app.EventDto, error)
	Update(ctx context.Context, eventDto app.EventDto) error
	Delete(ctx context.Context, userID uint64, eventID uint64) error
	ListForDay(ctx context.Context, userID uint64, date time.Time) ([]*app.EventDto, error)
	ListForWeek(ctx context.Context, userID uint64, startDate time.Time) ([]*app.EventDto, error)
	ListForMonth(ctx context.Context, userID uint64, startDate time.Time) ([]*app.EventDto, error)
}

func NewServer(logger Logger, app Application, addr string, readTimeout time.Duration) *Server {
	return &Server{
		logger:      logger,
		addr:        addr,
		readTimeout: readTimeout,
		handler:     NewEventHandler(logger, app),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting http server")

	mux := mux.NewRouter()

	mux.Handle("/hello", loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.helloHandler))).Methods("GET")

	mux.Handle("/events", loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.handler.create))).Methods("POST")
	mux.Handle(fmt.Sprintf("/events/{%s}", eventIDPath), loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.handler.update))).Methods("PUT")
	mux.Handle(fmt.Sprintf("/events/{%s}", eventIDPath), loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.handler.getByID))).Methods("GET")
	mux.Handle(fmt.Sprintf("/events/{%s}", eventIDPath), loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.handler.delete))).Methods("DELETE")
	mux.Handle("/events", loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.handler.listForDay))).Methods("GET").Queries(
		startDateQueryKey, startDateQueryValue,
		periodTypeQueryKey, periodDayQueryValue,
	)
	mux.Handle("/events", loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.handler.listForWeek))).Methods("GET").Queries(
		startDateQueryKey, startDateQueryValue,
		periodTypeQueryKey, periodWeekQueryValue,
	)
	mux.Handle("/events", loggingMiddleware(ctx, s.logger, http.HandlerFunc(s.handler.listForMonth))).Methods("GET").Queries(
		startDateQueryKey, startDateQueryValue,
		periodTypeQueryKey, periodMonthQueryValue,
	)

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

func (s *Server) helloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "hello world")
}
