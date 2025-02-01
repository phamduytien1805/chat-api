package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/phamduytien1805/package/common"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/server"
	"github.com/phamduytien1805/package/transport"
	userpb "github.com/phamduytien1805/proto/user"
	"github.com/phamduytien1805/user/usecase"
	"google.golang.org/grpc"
)

type Usecases struct {
	AuthBasic  *usecase.AuthBasicUserUsecase
	CreateUser *usecase.CreateUserUsecase
	GetUser    *usecase.GetUserUsecase
	VerifyUser *usecase.VerifyUserEmailUsecase
}

type GrpcServer struct {
	grpcPort string
	host     string
	logger   common.GrpcLog
	s        *grpc.Server

	userpb.UnimplementedUserServiceServer

	uc *Usecases
}

func NewGrpcServer(config *config.UserConfig, uc *Usecases) server.GrpcServer {
	srv := &GrpcServer{
		host:     config.Grpc.Server.Host,
		grpcPort: config.Grpc.Server.Port,
		logger:   common.NewGrpcLog(),
		uc:       uc,
	}
	srv.s = transport.InitializeGrpcServer("user_svc", srv.logger)
	return srv
}

func (srv *GrpcServer) Register() {
	userpb.RegisterUserServiceServer(srv.s, srv)
}

func (srv *GrpcServer) Run() {
	go func() {
		addr := fmt.Sprintf("%s:%s", srv.host, srv.grpcPort)
		srv.logger.Info("grpc server listening", slog.String("addr", addr))
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
		if err := srv.s.Serve(lis); err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (srv *GrpcServer) GracefulStop() error {
	srv.s.GracefulStop()
	return nil
}
