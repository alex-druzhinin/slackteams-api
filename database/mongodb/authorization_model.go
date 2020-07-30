package mongodb

import "time"

type slackBotAuthorization struct {
	ID          string    `bson:"_id"`
	AccessToken string    `bson:"accessToken"`
	Scope       string    `bson:"scope"`
	UserId      string    `bson:"userId"`
	TeamName    string    `bson:"teamName"`
	TeamId      string    `bson:"teamId"`
	CreatedAt   time.Time `bson:"createdAt"`
	Enabled     bool      `bson:"enabled"`

	Bot struct {
		BotUserId      string `bson:"botUserId"`
		BotAccessToken string `bson:"botAccessToken"`
	} `bson:"bot"`
}
