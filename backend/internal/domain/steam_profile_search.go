package domain

type SteamProfileSearch struct {
	AccountID   int    `json:"account_id"`
	Avatar      string `json:"avatar"`
	CountryCode string `json:"countrycode"`
	LastUpdated int64  `json:"last_updated"`
	Personaname string `json:"personaname"`
	Profileurl  string `json:"profileurl"`
	Realname    string `json:"realname"`
}
