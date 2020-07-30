package auth

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

var (
	ErrUserFound          = errors.New("user found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserConfirmed      = errors.New("user confirmed")
	ErrTokenNotFound      = errors.New("token not found")
	ErrNoPassword         = errors.New("no password")
	ErrWrongPassword      = errors.New("wrong password")
	ErrGoogleTokenExpired = errors.New("google token expired")
)

type CtxKey string

const (
	CtxKeyAuthUser CtxKey = "authUser"
)

type Service interface {
	FindUserByToken(ctx context.Context, token string) (string, error)
}

type service struct {
	repo  Repository
	cache *cache.Cache
}

type Config struct {
	UserRepository Repository
}

func NewAuthService(conf Config) (Service, error) {
	as := &service{
		repo:  conf.UserRepository,
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}

	return as, nil
}

func (a *service) FindUserByToken(ctx context.Context, token string) (string, error) {
	return a.repo.FindUserByToken(ctx, token)
}
