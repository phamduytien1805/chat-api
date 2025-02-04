package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/phamduytien1805/package/common"
	userpb "github.com/phamduytien1805/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *GrpcServer) CreateUserWithCredential(ctx context.Context, req *userpb.CreateUserForm) (*userpb.UserResponse, error) {
	createdUser, err := srv.uc.CreateUser.Exec(ctx, req.Username, req.Email, req.Credential)
	if err != nil {
		if errors.Is(err, common.ErrorUserResourceConflict) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, err
	}
	return &userpb.UserResponse{
		Id:            createdUser.ID.String(),
		Username:      createdUser.Username,
		Email:         createdUser.Email,
		EmailVerified: createdUser.EmailVerified,
	}, nil

}

func (srv *GrpcServer) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.UserResponse, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	user, err := srv.uc.GetUser.ById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		Id:            user.ID.String(),
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (srv *GrpcServer) GetUserByIdentifier(ctx context.Context, req *userpb.GetUserByIdentityRequest) (*userpb.UserResponse, error) {
	user, err := srv.uc.GetUser.ByEmailOrUsername(ctx, req.UsernameOrEmail)
	if err != nil {
		if errors.Is(err, common.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	return &userpb.UserResponse{
		Id:            user.ID.String(),
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (srv *GrpcServer) VerifyUserEmail(ctx context.Context, req *userpb.VerifyUserEmailRequest) (*userpb.Empty, error) {
	_, err := srv.uc.VerifyUser.Exec(ctx, req.UserEmail)
	if err != nil {
		return nil, err
	}

	return &userpb.Empty{}, nil
}
