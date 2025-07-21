package deadlockapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/quenyu/deadlock-stats/internal/domain"
)

const baseURL = "https://api.deadlock-api.com/v1"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) FetchPlayerCard(steamID string) (*DeadlockCard, error) {
	url := fmt.Sprintf("%s/players/%s/card", baseURL, steamID)
	var card DeadlockCard
	err := c.doRequest(url, &card)
	return &card, err
}

func (c *Client) FetchMatchHistory(steamID string) ([]DeadlockMatch, error) {
	url := fmt.Sprintf("%s/players/%s/match-history", baseURL, steamID)

	var apiMatches []DeadlockMatch
	if err := c.doRequest(url, &apiMatches); err != nil {
		return nil, err
	}

	return apiMatches, nil
}

func (c *Client) FetchHeroStats(steamID string) ([]domain.HeroStat, error) {
	url := fmt.Sprintf("%s/players/%s/hero-stats", baseURL, steamID)

	var apiHeroStats []HeroStatAPI
	if err := c.doRequest(url, &apiHeroStats); err != nil {
		return nil, err
	}

	domainHeroStats := make([]domain.HeroStat, len(apiHeroStats))
	for i, apiStat := range apiHeroStats {
		matchesCount := apiStat.MatchesPlayed
		if matchesCount == 0 {
			matchesCount = len(apiStat.Matches)
		}

		var winRate float64
		if matchesCount > 0 {
			winRate = (float64(apiStat.Wins) / float64(matchesCount)) * 100
		}

		var kda float64
		if apiStat.Deaths > 0 {
			kda = float64(apiStat.Kills+apiStat.Assists) / float64(apiStat.Deaths)
		} else {
			kda = float64(apiStat.Kills + apiStat.Assists)
		}

		domainHeroStats[i] = domain.HeroStat{
			HeroID:     apiStat.HeroID,
			HeroName:   fmt.Sprintf("HeroID_%d", apiStat.HeroID),
			Matches:    matchesCount,
			WinRate:    winRate,
			KDA:        kda,
			HeroAvatar: "",
		}
	}

	return domainHeroStats, nil
}

func (c *Client) FetchMMRHistory(steamID string) ([]DeadlockMMR, error) {
	url := fmt.Sprintf("%s/players/%s/mmr-history", baseURL, steamID)
	var mmrHistory []DeadlockMMR
	err := c.doRequest(url, &mmrHistory)
	return mmrHistory, err
}

func (c *Client) doRequest(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deadlock API returned non-200 status: %d for URL %s", resp.StatusCode, url)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
