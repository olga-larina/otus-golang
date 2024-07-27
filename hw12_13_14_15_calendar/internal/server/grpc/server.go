package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const userIDHeader = "X-USER-ID"

var errNotValidUserID = errors.New("userID is not valid")

//go:generate protoc -I ../../../api EventService.proto --go_out=. --go-grpc_out=.
type Server struct {
	logger   Logger
	app      Application
	grpcPort string
	srv      *grpc.Server
	pb.UnimplementedEventServiceServer
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

func NewServer(logger Logger, app Application, grpcPort string) *Server {
	return &Server{
		logger:   logger,
		app:      app,
		grpcPort: grpcPort,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(ctx, "starting grpc server", "port", s.grpcPort)

	lsn, err := net.Listen("tcp", fmt.Sprintf(":%s", s.grpcPort))
	if err != nil {
		s.logger.Error(ctx, err, "failed to create grpc server")
		return err
	}

	s.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggerInterceptor(s.logger),
		),
	)
	reflection.Register(s.srv)
	pb.RegisterEventServiceServer(s.srv, s)

	return s.srv.Serve(lsn)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "stopping grpc server")
	s.srv.GracefulStop()
	return nil
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	if req == nil || req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event := repackEventToDto(req.Event, userID)
	eventID, err := s.app.Create(ctx, *event)
	if err != nil {
		if errors.Is(err, storage.ErrBusyTime) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateEventResponse{Id: eventID}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.Event, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event, err := s.app.GetByID(ctx, userID, req.Id)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return repackEventToProto(event), nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*emptypb.Empty, error) {
	if req == nil || req.Event == nil || req.Event.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event := repackEventToDto(req.Event, userID)
	err = s.app.Update(ctx, *event)
	if err != nil {
		if errors.Is(err, storage.ErrBusyTime) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.app.Delete(ctx, userID, req.Id)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) EventListForDay(ctx context.Context, req *pb.EventListRequest) (*pb.EventList, error) {
	if req == nil || req.StartDate == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := s.app.ListForDay(ctx, userID, req.StartDate.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EventList{Events: repackEventsToProto(events)}, nil
}

func (s *Server) EventListForWeek(ctx context.Context, req *pb.EventListRequest) (*pb.EventList, error) {
	if req == nil || req.StartDate == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := s.app.ListForWeek(ctx, userID, req.StartDate.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EventList{Events: repackEventsToProto(events)}, nil
}

func (s *Server) EventListForMonth(ctx context.Context, req *pb.EventListRequest) (*pb.EventList, error) {
	if req == nil || req.StartDate == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := s.app.ListForMonth(ctx, userID, req.StartDate.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EventList{Events: repackEventsToProto(events)}, nil
}

func getUserID(ctx context.Context) (uint64, error) {
	var userID uint64
	var err error

	if md, exists := metadata.FromIncomingContext(ctx); exists {
		values := md.Get(userIDHeader)
		if len(values) > 0 {
			userID, err = strconv.ParseUint(values[0], 10, 64)
		} else {
			err = errNotValidUserID
		}
	} else {
		err = errNotValidUserID
	}

	if err == nil && userID == 0 {
		err = errNotValidUserID
	}

	return userID, err
}

func repackEventToDto(in *pb.Event, userID uint64) *app.EventDto {
	return &app.EventDto{
		ID:           in.Id,
		Title:        in.Title,
		StartDate:    in.StartDate.AsTime(),
		EndDate:      in.EndDate.AsTime(),
		Description:  in.Description,
		UserID:       userID,
		NotifyBefore: in.NotifyBefore.AsDuration(),
	}
}

func repackEventToProto(in *app.EventDto) *pb.Event {
	return &pb.Event{
		Id:           in.ID,
		Title:        in.Title,
		StartDate:    timestamppb.New(in.StartDate),
		EndDate:      timestamppb.New(in.EndDate),
		Description:  in.Description,
		NotifyBefore: durationpb.New(in.NotifyBefore),
	}
}

func repackEventsToProto(in []*app.EventDto) []*pb.Event {
	events := make([]*pb.Event, len(in))
	for i, event := range in {
		events[i] = repackEventToProto(event)
	}
	return events
}
