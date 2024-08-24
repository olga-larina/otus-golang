package internalgrpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	timeLayout   = "02/Jan/2006:15:04:05 -0700"
	healthMethod = "/grpc.health.v1.Health/Check"
)

func LoggerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)

		if info.FullMethod != healthMethod {
			elapsed := time.Since(startTime)

			var addr string
			if p, exists := peer.FromContext(ctx); exists {
				addr = p.Addr.String()
			}
			respStatus := status.Code(err).String()

			logger.Info(ctx, "grpc request",
				"ip", addr,
				"startTime", startTime.Format(timeLayout),
				"method", info.FullMethod,
				"statusCode", respStatus,
				"latency", elapsed.Milliseconds(),
			)
		}

		return resp, err
	}
}
