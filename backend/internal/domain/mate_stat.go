package domain

type MateStat struct {
	SteamID   string  `json:"steam_id"`
	Nickname  string  `json:"nickname"`
	AvatarURL string  `json:"avatar_url"`
	Games     int     `json:"games"`
	Wins      int     `json:"wins"`
	WinRate   float64 `json:"win_rate"`
}
