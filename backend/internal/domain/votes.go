package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	UserID      uuid.UUID `json:"user_id"`
	ContentType string    `json:"content_type"`
	ContentID   uuid.UUID `json:"content_id"`
	VoteValue   int       `json:"vote_value"`
	CreatedAt   time.Time `json:"created_at"`
}
