package domain

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID          uuid.UUID `json:"id"`
	AuthorID    uuid.UUID `json:"author_id"`
	ParentID    uuid.UUID `json:"parent_id"`
	ContentType string    `json:"content_type"`
	ContentID   uuid.UUID `json:"content_id"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
}
