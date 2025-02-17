package servers

import (
	"context"
	"errors"
	proto "github.com/bruceneco/go-template/internal/adapters/grpc/proto/gen"
	"github.com/bruceneco/go-template/internal/domain/user"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var _ proto.UserServiceServer = (*UserServer)(nil)

type UserServer struct {
	service  *user.Service
	validate *validator.Validate
	proto.UnimplementedUserServiceServer
}

func NewUserServer(s *user.Service, v *validator.Validate) *UserServer {
	return &UserServer{service: s, validate: v}
}
func (u *UserServer) Register(grpcServer *grpc.Server) {
	proto.RegisterUserServiceServer(grpcServer, u)
}

func (u *UserServer) CreateUser(
	ctx context.Context,
	request *proto.CreateUserRequest) (*proto.User, error) {
	err := u.validate.Struct(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	entity := user.Entity{
		Name:     request.GetName(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}
	err = u.service.CreateUser(ctx, &entity)
	switch {
	case errors.Is(err, user.ErrEmailAlreadyExists):
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.User{
		Id:    entity.ID.String(),
		Name:  entity.Name,
		Email: entity.Email,
	}, nil
}

func (u *UserServer) GetUser(
	ctx context.Context,
	request *proto.GetUserRequest) (*proto.User, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	entity, err := u.service.GetUserByID(ctx, id)
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.User{
		Id:    entity.ID.String(),
		Name:  entity.Name,
		Email: entity.Email,
	}, nil
}

func (u *UserServer) UpdateUser(
	ctx context.Context,
	request *proto.User) (*proto.User, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, user.ErrInvalidID.Error())
	}
	entity := user.Entity{
		ID:    &id,
		Name:  request.GetName(),
		Email: request.GetEmail(),
	}
	err = u.service.UpdateUser(ctx, &entity)
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case errors.Is(err, user.ErrEmailAlreadyExists):
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.User{
		Id:    entity.ID.String(),
		Name:  entity.Name,
		Email: entity.Email,
	}, nil
}

func (u *UserServer) DeleteUser(
	ctx context.Context,
	request *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, user.ErrInvalidID.Error())
	}
	err = u.service.DeleteUser(ctx, id)
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
