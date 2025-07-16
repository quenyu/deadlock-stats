package domain

import "github.com/google/uuid"

type ContentTag struct {
	TagID       int       `json:"tag_id"`
	ContentType string    `json:"content_type"`
	ContentID   uuid.UUID `json:"content_id"`
}
