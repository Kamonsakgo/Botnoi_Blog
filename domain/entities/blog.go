package entities

import (
	"time"
)

type Image struct {
	Name string `json:"name" bson:"name,omitempty"`
	URL  string `json:"url" bson:"url,omitempty"`
}
type BlogModel struct {
	BlogID       string    `json:"blog_id" bson:"blog_id,omitempty"`
	UserID       string    `json:"user_id" bson:"user_id,omitempty"`
	Title        string    `json:"title" bson:"title,omitempty"`
	Content      string    `json:"content" bson:"content,omitempty"`
	ImageURL     []Image   `json:"image_url" bson:"image_url,omitempty"`
	Category     string    `json:"category" bson:"category,omitempty"`
	Tag          []string  `json:"tag" bson:"tag,omitempty"`
	Type         []string  `json:"type" bson:"type,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at,omitempty"`
	LastUpdateAt time.Time `json:"last_update_at" bson:"last_update_at,omitempty"`
	Location     string    `json:"location" bson:"location,omitempty"`
	HL_id        string    `json:"hl_id" bson:"hl_id,omitempty"`
}

type BlogResponseModel struct {
	Blogs      []BlogModel `json:"blogs"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
	TotalCount int64       `json:"total_items"`
}
