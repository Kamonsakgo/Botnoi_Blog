package entities

import (
	"time"
)

type HighlightModel struct {
	HighlightID   string    `json:"highlight_id" bson:"highlight_id,omitempty"`
	UserID        string    `json:"user_id" bson:"user_id,omitempty"`
	Date          string    `json:"date" bson:"date,omitempty"`
	Speaker       string    `json:"speaker" bson:"speaker,omitempty"`
	Title         string    `json:"title" bson:"title,omitempty"`
	Content       string    `json:"content" bson:"content,omitempty"`
	ImageURL      string    `json:"image_url" bson:"image_url,omitempty"`
	Category      []string  `json:"category" bson:"category,omitempty"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at,omitempty"`
	LastUpdateAt  time.Time `json:"last_update_at" bson:"last_update_at,omitempty"`
	Location      string    `json:"location" bson:"location,omitempty"`
	LocationEvent string    `json:"location_event" bson:"location_event,omitempty"`
}
type HighlightResponseModel struct {
	Highlights []HighlightModel `json:"Highlight"`
	Page       int              `json:"page"`
	TotalPages int              `json:"total_pages"`
	TotalCount int64            `json:"total_items"`
}
