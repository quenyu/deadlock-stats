package domain

import (
	"time"

	"github.com/google/uuid"
)

type PlayerStats struct {
	UserID           uuid.UUID `json:"user_id"`
	KDRatio          float64   `json:"kd_ratio"`
	WinRate          float64   `json:"win_rate"`
	AvgMatchesPerDay float64   `json:"avg_matches_per_day"`
	FavoriteHero     string    `json:"favorite_hero"`
	LastUpdatedAt    time.Time `json:"last_updated_at"`
}
