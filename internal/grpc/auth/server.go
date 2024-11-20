package auth

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/lib/jwt"
	"sso/internal/services/auth"
	"sso/internal/storage"

	"github.com/go-playground/validator/v10"
	ssov1 "github.com/nikitauty/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(email string, password string, appID int32) (pair jwt.TokenPair, err error)
	RegisterNewUser(email string, password string) (userID int64, err error)
	IsAdmin(userID int64) (bool, error)
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
		if data.AppId == 0 {
			return nil, status.Error(codes.InvalidArgument, "wrong app id")
		}

		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("email is not valid %s", validationErrors))
	}

	pair, err := s.auth.Login(data.Email, data.Password, data.AppId)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidEmailOrPassword) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		if errors.Is(err, auth.ErrInvalidAppID) {
			return nil, status.Error(codes.InvalidArgument, "invalid app id")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: pair.AccessToken,
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

	userID, err := s.auth.RegisterNewUser(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
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

	isAdmin, err := s.auth.IsAdmin(data.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
