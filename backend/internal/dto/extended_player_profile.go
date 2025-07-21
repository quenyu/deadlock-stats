package dto

import (
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"github.com/quenyu/deadlock-stats/internal/domain"
)

type ExtendedPlayerProfile struct {
	Card         *deadlockapi.DeadlockCard `json:"card"`
	MatchHistory []domain.Match            `json:"match_history"`
	HeroStats    []domain.HeroStat         `json:"hero_stats"`
	MMRHistory   []deadlockapi.DeadlockMMR `json:"mmr_history"`

	TotalMatches        int                        `json:"total_matches"`
	WinRate             float64                    `json:"win_rate"`
	KDRatio             float64                    `json:"kd_ratio"`
	PerformanceDynamics domain.PerformanceDynamics `json:"performance_dynamics"`

	Nickname       string  `json:"nickname"`
	AvatarURL      string  `json:"avatar_url"`
	RankImage      string  `json:"rank_image"`
	AvgSoulsPerMin float64 `json:"avg_souls_per_min"`
}
