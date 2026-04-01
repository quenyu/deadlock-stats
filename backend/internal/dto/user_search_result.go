package dto

import "time"

type UserSearchResult struct {
	ID         string     `json:"id"`
	SteamID    string     `json:"steam_id"`
	Nickname   string     `json:"nickname"`
	AvatarURL  string     `json:"avatar_url"`
	ProfileURL string     `json:"profile_url"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`

	AccountID   int    `json:"account_id,omitempty"`
	CountryCode string `json:"countrycode,omitempty"`
	LastUpdated int64  `json:"last_updated,omitempty"`
	Realname    string `json:"realname,omitempty"`

	IsDeadlockPlayer    bool `json:"is_deadlock_player"`
	DeadlockStatusKnown bool `json:"deadlock_status_known"`
}
