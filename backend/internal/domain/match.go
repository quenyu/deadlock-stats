package domain

import "time"

type Match struct {
	ID                   string    `json:"match_id"`
	HeroID               int       `json:"hero_id"`
	PlayerKills          int       `json:"player_kills"`
	PlayerDeaths         int       `json:"player_deaths"`
	PlayerAssists        int       `json:"player_assists"`
	NetWorth             int       `json:"net_worth"`
	MatchDurationS       int       `json:"match_duration_s"`
	MatchResult          int       `json:"match_result"`
	StartTime            int64     `json:"start_time"`
	HeroName             string    `json:"hero_name"`
	HeroAvatar           string    `json:"hero_avatar,omitempty"`
	PlayerRankAfterMatch int       `json:"player_rank_after_match"`
	RankName             string    `json:"rank_name"`
	SubRank              int       `json:"sub_rank"`
	RankImage            string    `json:"rank_image"`
	PlayerRankChange     int       `json:"player_rank_change"`
	Kills                int       `json:"kills,omitempty"`
	Deaths               int       `json:"deaths,omitempty"`
	Assists              int       `json:"assists,omitempty"`
	DurationMinutes      int       `json:"duration_minutes,omitempty"`
	MatchTime            time.Time `json:"match_time,omitempty"`
	Result               string    `json:"result"`
}
