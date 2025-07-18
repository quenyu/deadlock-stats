package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	SteamID    string    `json:"steam_id"`
	Nickname   string    `json:"nickname"`
	AvatarURL  string    `json:"avatar_url"`
	ProfileURL string    `json:"profile_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
