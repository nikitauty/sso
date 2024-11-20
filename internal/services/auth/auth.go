package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	"sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
	refreshTTL   time.Duration
}

type UserSaver interface {
	SaveUser(email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	UserByEmail(email string) (models.User, error)
	UserByID(id int64) (models.User, error)
	IsAdmin(userID int64) (bool, error)
}

type AppProvider interface {
	App(appID int32) (models.App, error)
}

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInvalidAppID           = errors.New("invalid app id")
	ErrUserExists             = errors.New("user already exists")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
	refreshTTL time.Duration,
) *Auth {
	return &Auth{
		log,
		userSaver,
		userProvider,
		appProvider,
		tokenTTL,
		refreshTTL,
	}
}

func (a *Auth) Login(
	email string,
	password string,
	appID int32,
) (jwt.TokenPair, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.UserByEmail(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return jwt.TokenPair{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))

		return jwt.TokenPair{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return jwt.TokenPair{}, fmt.Errorf("%s: %w", op, ErrInvalidEmailOrPassword)
	}

	app, err := a.appProvider.App(appID)
	if err != nil {
		return jwt.TokenPair{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged successfully")

	tokens, err := jwt.NewTokenPair(user, app, a.tokenTTL, a.refreshTTL) // Access и Refresh токены
	if err != nil {
		log.Error("failed to generate tokens", sl.Err(err))
		return jwt.TokenPair{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *Auth) RegisterNewUser(
	email string,
	password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(
	userID int64,
) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		log.Error("failed to check is user admin", sl.Err(err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
