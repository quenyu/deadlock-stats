package domain

type PersonalRecords struct {
	MaxKills           int     `json:"max_kills"`
	MaxAssists         int     `json:"max_assists"`
	MaxNetWorth        int     `json:"max_net_worth"`
	BestKDA            float64 `json:"best_kda"`
	MaxKillsMatchID    string  `json:"max_kills_match_id"`
	MaxAssistsMatchID  string  `json:"max_assists_match_id"`
	MaxNetWorthMatchID string  `json:"max_net_worth_match_id"`
	BestKDAMatchID     string  `json:"best_kda_match_id"`
}
