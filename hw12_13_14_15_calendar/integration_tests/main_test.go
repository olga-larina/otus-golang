//go:build integration
// +build integration

package integration

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	sqlstorage "github.com/olga-larina/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const userIDHeader = "x-user-id"

var globalConfig *Config

type IntegrationTestSuite struct {
	suite.Suite
	cfg        *Config
	logg       *logger.Logger
	storage    *sqlstorage.Storage
	grpcClient pb.EventServiceClient
	grpcConn   *grpc.ClientConn
	httpClient *resty.Client
}

func (s *IntegrationTestSuite) SetupSuite() {
	var err error
	s.cfg = globalConfig
	ctx := context.Background()

	// logger
	s.logg, err = logger.New(s.cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed building logger %v", err)
		os.Exit(1)
	}

	// sql storage
	s.storage = sqlstorage.New(s.cfg.Database.Driver, s.cfg.Database.URI)
	if err := s.storage.Connect(ctx); err != nil {
		s.logg.Error(ctx, err, "failed to connect to db")
		os.Exit(1)
	}

	// grpc client
	s.grpcConn, err = grpc.NewClient(s.cfg.Calendar.GrpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logg.Error(ctx, err, "failed to connect to grpc")
		os.Exit(1)
	}
	s.grpcClient = pb.NewEventServiceClient(s.grpcConn)

	// http client
	s.httpClient = resty.New()
	s.httpClient.SetBaseURL(s.cfg.Calendar.HTTPUrl)

	s.logg.Info(ctx, "suite started")
}

func (s *IntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if s.storage != nil {
		defer func() {
			if err := s.storage.Close(ctx); err != nil {
				s.logg.Error(ctx, err, "failed to close sql storage")
			}
		}()
	}
	if s.grpcConn != nil {
		defer func() {
			if err := s.grpcConn.Close(); err != nil {
				s.logg.Error(ctx, err, "failed to close grpc connection")
			}
		}()
	}

	s.logg.Info(ctx, "suite finished")
}

func TestMain(m *testing.M) {
	var err error

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "/etc/integration_tests/config.yaml"
	}

	globalConfig, err = NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed reading config %v", err)
		os.Exit(1)
	}

	location, err := time.LoadLocation(globalConfig.Timezone)
	if err != nil {
		log.Fatalf("failed loading location %v", err)
		os.Exit(1)
	}
	time.Local = location

	code := m.Run()

	os.Exit(code)
}

func TestIntegrationTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IntegrationTestSuite))
}
