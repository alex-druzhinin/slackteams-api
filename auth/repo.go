package auth

import "context"

type Repository interface {
	FindUserByToken(ctx context.Context, token string) (string, error)
}
