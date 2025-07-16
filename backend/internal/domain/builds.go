package domain

import (
	"time"

	"github.com/google/uuid"
)

type Build struct {
	ID          uuid.UUID `json:"id"`
	AuthorID    uuid.UUID `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GameVersion string    `json:"game_version"`
	IsPublic    bool      `json:"is_public"`
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
