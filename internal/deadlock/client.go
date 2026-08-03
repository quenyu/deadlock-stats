package deadlock

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxResponseBytes = 16 << 20 // 16 MiB

type Match struct {
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
	PlayerTeam     int   `json:"player_team"`
}

// Won compares the player's team with the winning team returned by the API.
// It must not be reduced to match_result == 1: players on team 0 can also win.
func (m Match) Won() bool {
	return m.PlayerTeam == m.MatchResult
}

type HistoryClient interface {
	MatchHistory(ctx context.Context, accountID string) ([]Match, error)
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) MatchHistory(ctx context.Context, accountID string) ([]Match, error) {
	endpoint := fmt.Sprintf("%s/v1/players/%s/match-history", c.baseURL, url.PathEscape(accountID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create match-history request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "deadlock-ghost-match/0.1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch match history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("deadlock API returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var matches []Match
	decoder := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes))
	if err := decoder.Decode(&matches); err != nil {
		return nil, fmt.Errorf("decode match history: %w", err)
	}

	return matches, nil
}
