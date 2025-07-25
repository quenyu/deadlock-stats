package domain

type DeadlockMMR struct {
	MatchID      int64   `json:"match_id"`
	Rank         int     `json:"rank"`
	StartTime    int64   `json:"start_time"`
	PlayerScore  float64 `json:"player_score"`
	Division     int     `json:"division"`
	DivisionTier int     `json:"division_tier"`
}

type HeroMMRHistory struct {
	HeroID   int           `json:"hero_id"`
	HeroName string        `json:"hero_name"`
	History  []DeadlockMMR `json:"history"`
}
