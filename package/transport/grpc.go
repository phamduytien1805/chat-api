package transport

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/phamduytien1805/package/common"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

var (
	ServiceIdHeader string = "Service-Id"
)

func interceptorLogger(l common.GrpcLog) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, fields...)
		case logging.LevelInfo:
			l.Info(msg, fields...)
		case logging.LevelWarn:
			l.Warn(msg, fields...)
		case logging.LevelError:
			l.Error(msg, fields...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func InitializeGrpcServer(name string, logger common.GrpcLog) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 8), // increase to 8 MB (default: 4 MB)
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // terminate the connection if a client pings more than once every 5 seconds
			PermitWithoutStream: true,            // allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,  // if a client is idle for 15 seconds, send a GOAWAY
			MaxConnectionAge:      600 * time.Second, // if any connection is alive for more than maxConnectionAge, send a GOAWAY
			MaxConnectionAgeGrace: 5 * time.Second,   // allow 5 seconds for pending RPCs to complete before forcibly closing connections
			Time:                  5 * time.Second,   // ping the client if it is idle for 5 seconds to ensure the connection is still active
			Timeout:               1 * time.Second,   // wait 1 second for the ping ack before assuming the connection is dead
		}),
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		logger.Error("recovered from panic, stack: " + string(debug.Stack()))
		return status.Errorf(codes.Internal, "%s", p)
	}
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithDurationField(logging.DurationToTimeMillisFields),
	}
	logTraceID := func(ctx context.Context) logging.Fields {
		return nil
	}

	opts = append(opts,
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptorLogger(logger), logOpts...),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger), logging.WithFieldsFromContext(logTraceID)),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)
	grpcSrv := grpc.NewServer(opts...)
	return grpcSrv
}

func InitializeGrpcClient(svcHost string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	scheme := "dns"

	retryOpts := []retry.CallOption{
		// generate waits between 900ms to 1100ms
		retry.WithBackoff(retry.BackoffLinearWithJitter(1*time.Second, 0.1)),
		retry.WithMax(3),
		retry.WithCodes(codes.Unavailable, codes.Aborted),
		retry.WithPerRetryTimeout(3 * time.Second),
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	dialOpts = append(dialOpts,
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(`{
			"loadBalancingPolicy": "round_robin"
		}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		}),
		grpc.WithChainStreamInterceptor(
			retry.StreamClientInterceptor(retryOpts...),
		),
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(retryOpts...),
		),
		//grpc.WithBlock(),
	)

	slog.Info("connecting to grpc host: " + svcHost)
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:///%s", scheme, svcHost),
		dialOpts...,
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewGrpcEndpoint(conn *grpc.ClientConn, serviceID, serviceName, method string, grpcReply interface{}) endpoint.Endpoint {
	var options []grpctransport.ClientOption
	var (
		ep         endpoint.Endpoint
		endpointer sd.FixedEndpointer
	)

	ep = grpctransport.NewClient(
		conn,
		serviceName,
		method,
		encodeGRPCRequest,
		decodeGRPCResponse,
		grpcReply,
		append(options, grpctransport.ClientBefore(grpctransport.SetRequestHeader(ServiceIdHeader, serviceID)))...,
	).Endpoint()
	ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    serviceName + "." + method,
		Timeout: 60 * time.Second,
	}))(ep)
	endpointer = append(endpointer, ep)
	// timeout for the whole invocation
	ep = lb.Retry(1, 15*time.Second, lb.NewRoundRobin(endpointer))

	return ep
}

func encodeGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeGRPCResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	return grpcReply, nil
}
