package domain

import "time"

type Trend struct {
	Trend     string    `json:"trend"` // "up", "down", "stable"
	Value     string    `json:"value"`
	Sparkline []float64 `json:"sparkline"`
}

type PerformanceDynamics struct {
	WinLoss Trend `json:"win_loss"`
	KDA     Trend `json:"kda"`
	Rank    Trend `json:"rank"`
}

type PlayerProfile struct {
	SteamID               string              `json:"steam_id" gorm:"primaryKey"`
	Nickname              string              `json:"nickname"`
	AvatarURL             string              `json:"avatar_url"`
	ProfileURL            string              `json:"profile_url"`
	CreatedAt             time.Time           `json:"created_at"`
	LastMatchTime         time.Time           `json:"last_match_time"`
	PlayerRank            int                 `json:"player_rank"`
	RankName              string              `json:"rank_name"`
	SubRank               int                 `json:"sub_rank"`
	RankImage             string              `json:"rank_image"`
	WinRate               float64             `json:"win_rate"`
	KDRatio               float64             `json:"kd_ratio"`
	AvgMatchesPerDay      float64             `json:"avg_matches_per_day"`
	FavoriteHero          string              `json:"favorite_hero"`
	LastUpdatedAt         time.Time           `json:"last_updated_at"`
	TotalMatches          int                 `json:"total_matches"`
	TotalKills            int                 `json:"total_kills"`
	TotalDeaths           int                 `json:"total_deaths"`
	TotalAssists          int                 `json:"total_assists"`
	MaxKillsInMatch       int                 `json:"max_kills_in_match"`
	AvgDamagePerMatch     float64             `json:"avg_damage_per_match"`
	AvgObjectivesPerMatch float64             `json:"avg_objectives_per_match"`
	AvgSoulsPerMin        float64             `json:"avg_souls_per_min"`
	RecentMatches         []Match             `json:"recent_matches" gorm:"-"`
	HeroStats             []HeroStat          `json:"hero_stats" gorm:"-"`
	PerformanceDynamics   PerformanceDynamics `json:"performance_dynamics" gorm:"-"`
}
