package deadlockapi

type RankImage struct {
	Large *string `json:"large"`
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

type DeadlockCardSlot struct {
	Hero struct {
		ID int `json:"id"`
	} `json:"hero"`
}

type DeadlockCard struct {
	AccountID        *int               `json:"account_id"`
	RankedBadgeLevel *int               `json:"ranked_badge_level"`
	RankedRank       *int               `json:"ranked_rank"`
	RankedSubrank    *int               `json:"ranked_subrank"`
	Slots            []DeadlockCardSlot `json:"slots"`
}

type DeadlockMMR struct {
	MatchID      int64   `json:"match_id"`
	Rank         int     `json:"rank"`
	StartTime    int64   `json:"start_time"`
	PlayerScore  float64 `json:"player_score"`
	Division     int     `json:"division"`
	DivisionTier int     `json:"division_tier"`
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
