package rpc

import (
	"time"
)

type SlackIcon struct {
	Image34      string
	Image44      string
	Image68      string
	Image88      string
	Image102     string
	Image132     string
	Image230     string
	ImageDefault bool `json:"imageDefault"`
}

type SlackTeam struct {
	ID          string     `json:"teamId"`
	Name        string     `json:"name"`
	Domain      string     `json:"domain"`
	EmailDomain string     `json:"emailDomain"`
	Icon        SlackIcon  `json:"icon"`
	IsDeleted   bool       `json:"isDeleted"`
	DeletedAt   *time.Time `json:"deletedAt, omitempty"`
	CreatedAt   time.Time  `json:"createdAt, omitempty"`
	Tags        *[]string  `json:"tags, omitempty"`
}
