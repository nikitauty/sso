package jwt

import (
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	UserID int64 `json:"user_id"`
	AppID  int32 `json:"app_id"`
	jwt.RegisteredClaims
}

func NewTokenPair(user models.User, app models.App, accessTTL, refreshTTL time.Duration) (TokenPair, error) {
	now := time.Now()

	accessClaims := &Claims{
		UserID: user.ID,
		AppID:  app.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(app.Secret)
	if err != nil {
		return TokenPair{}, err
	}

	refreshClaims := &Claims{
		UserID: user.ID,
		AppID:  app.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTTL)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(app.RefreshSecret)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func ValidateToken(app models.App, tokenStr string, isRefresh bool) (*Claims, error) {
	key := app.Secret
	if isRefresh {
		key = app.RefreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}
