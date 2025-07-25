package domain

type MateStatAPI struct {
	MateID        int `json:"mate_id"`
	Wins          int `json:"wins"`
	MatchesPlayed int `json:"matches_played"`
}
