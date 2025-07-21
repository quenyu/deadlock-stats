package domain

type HeroStat struct {
	HeroID     int     `json:"hero_id"`
	HeroName   string  `json:"hero_name"`
	Matches    int     `json:"matches_played"`
	WinRate    float64 `json:"win_rate"`
	KDA        float64 `json:"kda"`
	HeroAvatar string  `json:"hero_avatar,omitempty"`
}
