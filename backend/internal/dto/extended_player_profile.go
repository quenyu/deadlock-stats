package dto

import (
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"github.com/quenyu/deadlock-stats/internal/domain"
)

type ExtendedPlayerProfile struct {
	Card                *deadlockapi.DeadlockCard  `json:"card"`
	MatchHistory        []domain.Match             `json:"match_history"`
	HeroStats           []domain.HeroStat          `json:"hero_stats"`
	MMRHistory          []domain.DeadlockMMR       `json:"mmr_history"`
	TotalMatches        int                        `json:"total_matches"`
	WinRate             float64                    `json:"win_rate"`
	KDRatio             float64                    `json:"kd_ratio"`
	PerformanceDynamics domain.PerformanceDynamics `json:"performance_dynamics"`

	PlayerRank     int     `json:"player_rank"`
	Nickname       string  `json:"nickname"`
	AvatarURL      string  `json:"avatar_url"`
	RankImage      string  `json:"rank_image"`
	RankName       string  `json:"rank_name"`
	SubRank        int     `json:"sub_rank"`
	AvgSoulsPerMin float64 `json:"avg_souls_per_min"`

	FeaturedHeroes     []domain.FeaturedHero   `json:"featured_heroes"`
	PeakRank           int                     `json:"peak_rank"`
	PeakRankName       string                  `json:"peak_rank_name"`
	PeakRankImage      string                  `json:"peak_rank_image"`
	PersonalRecords    domain.PersonalRecords  `json:"personal_records"`
	AvgKillsPerMatch   float64                 `json:"avg_kills_per_match"`
	AvgDeathsPerMatch  float64                 `json:"avg_deaths_per_match"`
	AvgAssistsPerMatch float64                 `json:"avg_assists_per_match"`
	AvgMatchDuration   float64                 `json:"avg_match_duration"`
	MateStats          []domain.MateStat       `json:"mate_stats"`
	HeroMMRHistory     []domain.HeroMMRHistory `json:"hero_mmr_history"`
}
