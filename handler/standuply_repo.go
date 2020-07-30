package handler

import "context"

type AuthorizationsRepository interface {
	GetAllAuthorizations(ctx context.Context) ([]*SlackBotAuthorization, error)
}
