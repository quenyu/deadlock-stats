package domain

import (
	"time"

	"github.com/google/uuid"
)

type Crosshair struct {
	ID            uuid.UUID `json:"id"`
	AuthorID      uuid.UUID `json:"author_id"`
	Title         string    `json:"title"`
	CrosshairCode string    `json:"crosshair_code"`
	IsPublic      bool      `json:"is_public"`
	ViewCount     int       `json:"view_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
