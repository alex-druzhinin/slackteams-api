package database

import (
	"context"

	"bitbucket.org/iwlab-standuply/slackteams-api/auth"
	"bitbucket.org/iwlab-standuply/slackteams-api/config"
	log "github.com/sirupsen/logrus"
)

type userRepository struct {
	users []config.User
}

func NewLocalAuthRepository(users []config.User) auth.Repository {
	return &userRepository{
		users,
	}
}

func (r *userRepository) FindUserByToken(ctx context.Context, token string) (string, error) {
	for i := range r.users {
		if r.users[i].Token == token {
			log.Debugf("Found user %s", r.users[i].Name)
			return r.users[i].Name, nil
		}
	}

	return "", auth.ErrUserFound
}
