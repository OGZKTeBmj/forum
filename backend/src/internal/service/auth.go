package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	models "github.com/OGZKTeBmj/forum/domain"
	"github.com/OGZKTeBmj/forum/internal/storage"
	"github.com/OGZKTeBmj/forum/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	Username string
	UserID   []byte
	jwt.RegisteredClaims
}

type AuthService struct {
	UserProvider UserProvider
	TokenTTL     time.Duration
	Log          *slog.Logger
	Secret       []byte
}

type UserProvider interface {
	User(ctx context.Context, id []byte) (user models.User, err error)
	UserByName(ctx context.Context, name string) (user models.User, err error)
	SaveUser(ctx context.Context, name string, passhash []byte) (id []byte, err error)
	UpdateUser(ctx context.Context, user models.User) error
}

func (a *AuthService) SignUp(ctx context.Context, name string, password string) ([]byte, error) {
	const op string = "authService.Register"

	log := a.Log.With(
		slog.String("op", op),
	)

	_, err := a.UserProvider.UserByName(ctx, name)
	if err == nil {
		return nil, utils.ErrWrap(op, ErrUserIsExist)
	}
	if !errors.Is(err, storage.ErrIsNotExist) {
		log.Error("can't check user")
		return nil, utils.ErrWrap(op, err)
	}

	log.Info("Attempting to register user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed generate password hash", utils.SlogErr(err))
		return nil, utils.ErrWrap(op, err)
	}

	id, err := a.UserProvider.SaveUser(ctx, name, passwordHash)
	if err != nil {
		log.Error("failed to save user", utils.SlogErr(err))
		return nil, utils.ErrWrap(op, err)
	}

	return id, nil
}

func (a *AuthService) SignIn(ctx context.Context, name string, password string) (string, error) {
	const op string = "authService.Login"

	log := a.Log.With(
		slog.String("op", op),
	)

	user, err := a.UserProvider.UserByName(ctx, name)

	if err != nil {
		if errors.Is(err, storage.ErrIsNotExist) {
			return "", utils.ErrWrap(op, ErrInvalidCredentails)
		}
		log.Error("failed to get user", utils.SlogErr(err))
		return "", utils.ErrWrap(op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.Log.Info("invalid", utils.SlogErr(err))

		return "", utils.ErrWrap(op, ErrInvalidCredentails)
	}

	token, err := a.newToken(a.Secret, user)
	if err != nil {
		a.Log.Error("failed to generate token")

		return "", utils.ErrWrap(op, err)
	}

	return token, nil
}

func (a *AuthService) User(ctx context.Context, id []byte) (user models.User, err error) {
	const op = "authService.User"

	user, err = a.UserProvider.User(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrIsNotExist) {
			return models.User{}, utils.ErrWrap(op, ErrUserIsNotExist)
		}
		return models.User{}, utils.ErrWrap(op, err)
	}
	return user, nil
}

func (a *AuthService) UpdateUser(ctx context.Context, user models.User) error {
	const op = "authService.UpdateUser"

	err := a.UserProvider.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrIsNotExist) {
			return utils.ErrWrap(op, ErrUserIsNotExist)
		}
		return utils.ErrWrap(op, err)
	}
	return nil
}

func (a *AuthService) ParseToken(token string) ([]byte, error) {
	const op = "authService.ParseToken"

	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", token.Header["alg"])
		}
		return a.Secret, nil
	})

	if err != nil {
		return nil, utils.ErrWrap(op, err)
	}

	if claims, ok := jwtToken.Claims.(*UserClaims); ok && jwtToken.Valid {
		return claims.UserID, nil
	}

	return nil, utils.ErrWrap(op, ErrInvalidToken)
}

func (a *AuthService) newToken(secret []byte, user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Username: user.Name,
		UserID:   user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(user.Id),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", utils.ErrWrap("can't create user token", err)
	}
	return tokenString, nil

}
