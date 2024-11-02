package auth

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/services/auth"
	"sso/internal/storage"

	"github.com/go-playground/validator/v10"
	ssov1 "github.com/nikitauty/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	data := LoginReq{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppId:    req.GetAppId(),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(data); err != nil {
		if data.Email == "" {
			return nil, status.Error(codes.InvalidArgument, "email is required")
		}
		if data.Password == "" {
			return nil, status.Error(codes.InvalidArgument, "password is required")
		}
		if data.AppId != 1 {
			return nil, status.Error(codes.InvalidArgument, "wrong app id")
		}

		validationErrors := err.(validator.ValidationErrors)

		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("email is not valid %s", validationErrors))
	}

	token, err := s.auth.Login(ctx, data.Email, data.Password, int(data.AppId))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	data := RegisterReq{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(data); err != nil {
		if data.Email == "" {
			return nil, status.Error(codes.InvalidArgument, "email is required")
		}
		if data.Password == "" {
			return nil, status.Error(codes.InvalidArgument, "password is required")
		}

		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("email is not valid %s", validationErrors))
	}

	userID, err := s.auth.RegisterNewUser(ctx, data.Email, data.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	data := IsAdminReq{
		UserID: req.GetUserId(),
	}

	if data.UserID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, data.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
