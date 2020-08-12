package mongodb

import "time"

type SlackIcon struct {
	Image34      string `bson:"image34"`
	Image44      string `bson:"image44"`
	Image68      string `bson:"image68"`
	Image88      string `bson:"image88"`
	Image102     string `bson:"image102"`
	Image132     string `bson:"image132"`
	Image230     string `bson:"image230"`
	ImageDefault bool   `bson:"imageDefault"`
}

type slackTeam struct {
	ID          string     `bson:"_id"`
	TeamID      string     `bson:"id"`
	Name        string     `bson:"name"`
	Domain      string     `bson:"domain"`
	EmailDomain string     `bson:"emailDomain"`
	Icon        SlackIcon  `bson:"icon"`
	IsDeleted   bool       `bson:"isDeleted"`
	DeletedAt   *time.Time `bson:"deletedAt"`
	CreatedAt   time.Time  `bson:"createdAt"`
	Tags        *[]string  `bson:"enabled"`
}
