package rpc

import "context"

type SlackTeamsRepository interface {
	FindTeamByID(ctx context.Context, teamID string) (*SlackTeam, error)
}
