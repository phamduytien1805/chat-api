package userclient

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/phamduytien1805/auth/domain"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/transport"
	userpb "github.com/phamduytien1805/proto/user"
	"google.golang.org/grpc"
)

var UserClientConn *grpc.ClientConn

type UserClientService struct {
	conn *grpc.ClientConn

	createUserWithCredential endpoint.Endpoint
	getUserByIdentity        endpoint.Endpoint
	getUserById              endpoint.Endpoint
	verifyUserEmail          endpoint.Endpoint
}

func NewUserClientService(config *config.UserConfig) (domain.UserService, error) {
	fmt.Println("Connecting to user service", config)
	hostsvc := fmt.Sprintf("%s:%s", config.Grpc.Server.Host, config.Grpc.Server.Port)
	conn, err := transport.InitializeGrpcClient(hostsvc)
	if err != nil {
		return nil, err
	}
	UserClientConn = conn

	return &UserClientService{
		conn: conn,
		createUserWithCredential: transport.NewGrpcEndpoint(
			conn,
			"auth",
			"auth.UserService",
			"CreateUserWithCredential",
			&userpb.CreateUserForm{},
		),
		getUserByIdentity: transport.NewGrpcEndpoint(
			conn,
			"auth",
			"auth.UserService",
			"GetUserByIdentifier",
			&userpb.GetUserByIdentityRequest{},
		),
		getUserById: transport.NewGrpcEndpoint(
			conn,
			"auth",
			"auth.UserService",
			"GetUserById",
			&userpb.GetUserByIdRequest{},
		),
		verifyUserEmail: transport.NewGrpcEndpoint(
			conn,
			"auth",
			"auth.UserService",
			"VerifyUserEmail",
			&userpb.VerifyUserEmailRequest{},
		),
	}, nil
}
func (s *UserClientService) CreateUserWithCredential(ctx context.Context, username string, email string, hashed_password string) (*domain.User, error) {
	req := &userpb.CreateUserForm{
		Username:   username,
		Email:      email,
		Credential: hashed_password,
	}
	resp, err := s.createUserWithCredential(ctx, req)
	if err != nil {
		return nil, err
	}
	user := resp.(*userpb.UserResponse)
	return mapUserResponseToDomainUser(user)
}

func (s *UserClientService) VerifyUserByIdentity(ctx context.Context, identity string, hashed_password string) (*domain.User, error) {
	req := &userpb.GetUserByIdentityRequest{
		UsernameOrEmail: identity,
		Credential:      hashed_password,
	}
	resp, err := s.getUserByIdentity(ctx, req)
	if err != nil {
		return nil, err
	}
	user := resp.(*userpb.UserResponse)
	return mapUserResponseToDomainUser(user)
}

func (s *UserClientService) GetUserById(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	req := &userpb.GetUserByIdRequest{
		Id: userID.String(),
	}
	resp, err := s.getUserById(ctx, req)
	if err != nil {
		return nil, err
	}
	user := resp.(*userpb.UserResponse)
	return mapUserResponseToDomainUser(user)
}

func (s *UserClientService) VerifyUserEmail(ctx context.Context, userEmail string) error {
	req := &userpb.VerifyUserEmailRequest{
		UserEmail: userEmail,
	}
	_, err := s.verifyUserEmail(ctx, req)
	return err
}

func mapUserResponseToDomainUser(user *userpb.UserResponse) (*domain.User, error) {
	userId, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:       userId,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
