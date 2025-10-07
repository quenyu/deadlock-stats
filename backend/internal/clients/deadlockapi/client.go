package deadlockapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/quenyu/deadlock-stats/internal/domain"
)

const baseURL = "https://api.deadlock-api.com/v1"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
		},
	}
}

func NewClientWithCustomTimeout(timeout time.Duration) *Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

func (c *Client) FetchMatchHistory(steamID string) ([]DeadlockMatch, error) {
	url := fmt.Sprintf("%s/players/%s/match-history", baseURL, steamID)

	var apiMatches []DeadlockMatch
	if err := c.doRequestWithRetry(url, &apiMatches, 2); err != nil {
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

	return c.convertToDomainHeroStats(apiHeroStats), nil
}

func (c *Client) convertToDomainHeroStats(apiHeroStats []HeroStatAPI) []domain.HeroStat {
	domainHeroStats := make([]domain.HeroStat, len(apiHeroStats))

	for i, apiStat := range apiHeroStats {
		domainHeroStats[i] = c.convertSingleHeroStat(apiStat)
	}

	return domainHeroStats
}

func (c *Client) convertSingleHeroStat(apiStat HeroStatAPI) domain.HeroStat {
	matchesCount := c.calculateMatchesCount(apiStat)
	winRate := c.calculateWinRate(apiStat.Wins, matchesCount)
	kda := c.calculateKDA(apiStat.Kills, apiStat.Deaths, apiStat.Assists)

	return domain.HeroStat{
		HeroID:     apiStat.HeroID,
		HeroName:   fmt.Sprintf("HeroID_%d", apiStat.HeroID),
		Matches:    matchesCount,
		WinRate:    winRate,
		KDA:        kda,
		HeroAvatar: "",
	}
}

func (c *Client) calculateMatchesCount(apiStat HeroStatAPI) int {
	if apiStat.MatchesPlayed > 0 {
		return apiStat.MatchesPlayed
	}
	return len(apiStat.Matches)
}

func (c *Client) calculateWinRate(wins, matchesCount int) float64 {
	if matchesCount > 0 {
		return (float64(wins) / float64(matchesCount)) * 100
	}
	return 0.0
}

func (c *Client) calculateKDA(kills, deaths, assists int) float64 {
	if deaths > 0 {
		return float64(kills+assists) / float64(deaths)
	}
	return float64(kills + assists)
}

func (c *Client) FetchMMRHistory(steamID string) ([]domain.DeadlockMMR, error) {
	url := fmt.Sprintf("%s/players/%s/mmr-history", baseURL, steamID)
	var mmrHistory []domain.DeadlockMMR
	err := c.doRequest(url, &mmrHistory)
	return mmrHistory, err
}

func (c *Client) FetchMateStats(steamID string) ([]domain.MateStatAPI, error) {
	url := fmt.Sprintf("%s/players/%s/mate-stats", baseURL, steamID)
	var mateStats []domain.MateStatAPI
	err := c.doRequest(url, &mateStats)
	if err != nil {
		return nil, err
	}
	return mateStats, nil
}

func (c *Client) FetchMMRHistoryByHero(steamID string, heroID int) ([]domain.DeadlockMMR, error) {
	url := fmt.Sprintf("%s/players/%s/mmr-history/%d", baseURL, steamID, heroID)
	var mmrHistory []domain.DeadlockMMR
	err := c.doRequest(url, &mmrHistory)
	if err != nil {
		return nil, err
	}
	return mmrHistory, nil
}

func (c *Client) FetchLiteProfile(steamID string) (*domain.PlayerProfile, error) {
	url := fmt.Sprintf("%s/players/%s/profile", baseURL, steamID)
	var profile domain.PlayerProfile
	err := c.doRequest(url, &profile)
	return &profile, err
}

func (c *Client) FetchSteamProfileSearch(query string) ([]domain.SteamProfileSearch, error) {
	encodedQuery := url.QueryEscape(query)
	url := fmt.Sprintf("%s/players/steam-search?search_query=%s", baseURL, encodedQuery)
	var profileSearch []domain.SteamProfileSearch
	err := c.doRequestWithRetry(url, &profileSearch, 2)
	return profileSearch, err
}

func (c *Client) doRequest(url string, target interface{}) error {
	req, err := c.createRequest(url)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.validateResponse(resp, url); err != nil {
		return err
	}

	return c.decodeResponse(resp, target)
}

func (c *Client) doRequestWithRetry(url string, target interface{}, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := c.doRequest(url, target)
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt == maxRetries {
			break
		}

		backoff := time.Duration(100*(1<<attempt)) * time.Millisecond
		time.Sleep(backoff)
	}

	return fmt.Errorf("request failed after %d attempts: %w", maxRetries+1, lastErr)
}

func (c *Client) createRequest(url string) (*http.Request, error) {
	return http.NewRequest("GET", url, nil)
}

func (c *Client) executeRequest(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) validateResponse(resp *http.Response, url string) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deadlock API returned non-200 status: %d for URL %s", resp.StatusCode, url)
	}
	return nil
}

func (c *Client) decodeResponse(resp *http.Response, target interface{}) error {
	return json.NewDecoder(resp.Body).Decode(target)
}
