package deadlockapi

type RankImage struct {
	Large         *string `json:"large"`
	LargeSubrank1 *string `json:"large_subrank1"`
	LargeSubrank2 *string `json:"large_subrank2"`
	LargeSubrank3 *string `json:"large_subrank3"`
	LargeSubrank4 *string `json:"large_subrank4"`
	LargeSubrank5 *string `json:"large_subrank5"`
	LargeSubrank6 *string `json:"large_subrank6"`
}

type RankV2 struct {
	Tier   int       `json:"tier"`
	Name   string    `json:"name"`
	Images RankImage `json:"images"`
	Color  *string   `json:"color"`
}

type HeroImages struct {
	IconHeroCard *string `json:"icon_hero_card"`
}

type HeroV2 struct {
	ID        int        `json:"id"`
	ClassName string     `json:"class_name"`
	Name      string     `json:"name"`
	Images    HeroImages `json:"images"`
}

type CardHero struct {
	ID    int `json:"id"`
	Kills int `json:"kills,omitempty"`
	Wins  int `json:"wins,omitempty"`
}

type CardStat struct {
	StatID    int `json:"stat_id"`
	StatScore int `json:"stat_score"`
}

type DeadlockCardSlot struct {
	Hero CardHero  `json:"hero"`
	Stat *CardStat `json:"stat,omitempty"`
}

type DeadlockCard struct {
	AccountID        *int               `json:"account_id"`
	RankedBadgeLevel *int               `json:"ranked_badge_level"`
	RankedRank       *int               `json:"ranked_rank"`
	RankedSubrank    *int               `json:"ranked_subrank"`
	Slots            []DeadlockCardSlot `json:"slots"`
}

type DeadlockMatch struct {
	MatchID        int64 `json:"match_id"`
	HeroID         int   `json:"hero_id"`
	PlayerKills    int   `json:"player_kills"`
	PlayerDeaths   int   `json:"player_deaths"`
	PlayerAssists  int   `json:"player_assists"`
	Denies         int   `json:"denies"`
	NetWorth       int   `json:"net_worth"`
	MatchDurationS int   `json:"match_duration_s"`
	MatchResult    int   `json:"match_result"`
	StartTime      int64 `json:"start_time"`
}

type HeroStatAPI struct {
	HeroID        int     `json:"hero_id"`
	Matches       []int64 `json:"matches"`
	MatchesPlayed int     `json:"matches_played"`
	Wins          int     `json:"wins"`
	Kills         int     `json:"kills"`
	Deaths        int     `json:"deaths"`
	Assists       int     `json:"assists"`
}
