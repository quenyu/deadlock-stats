package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CrosshairSettings struct {
	Color             string  `json:"color"`
	Thickness         int     `json:"thickness"`
	Length            int     `json:"length"`
	Gap               int     `json:"gap"`
	Dot               bool    `json:"dot"`
	Opacity           float64 `json:"opacity"`
	PipOpacity        float64 `json:"pipOpacity"`
	DotOutlineOpacity float64 `json:"dotOutlineOpacity"`
	HitMarkerDuration float64 `json:"hitMarkerDuration"`
	PipBorder         bool    `json:"pipBorder"`
	PipGapStatic      bool    `json:"pipGapStatic"`
}

type Crosshair struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	AuthorID    uuid.UUID       `json:"author_id" db:"author_id"`
	Author      *User           `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Title       string          `json:"title" db:"title"`
	Description string          `json:"description" db:"description"`
	Settings    json.RawMessage `json:"settings" db:"settings" gorm:"type:jsonb"`
	LikesCount  int             `json:"likes_count" db:"likes_count"`
	IsPublic    bool            `json:"is_public" db:"is_public" gorm:"default:false"`
	ViewCount   int             `json:"view_count" db:"view_count" gorm:"default:0"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type CrosshairLike struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CrosshairID uuid.UUID `json:"crosshair_id" db:"crosshair_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
