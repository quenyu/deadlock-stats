package domain

type FeaturedHero struct {
	HeroID    int    `json:"hero_id"`
	HeroName  string `json:"hero_name"`
	HeroImage string `json:"hero_image"`
	Kills     int    `json:"kills,omitempty"`
	Wins      int    `json:"wins,omitempty"`
	StatID    int    `json:"stat_id,omitempty"`
	StatScore int    `json:"stat_score,omitempty"`
}
