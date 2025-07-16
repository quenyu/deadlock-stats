package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/quenyu/deadlock-stats/internal/config"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/quenyu/deadlock-stats/internal/repository"
	"github.com/yohcop/openid-go"
)

type AuthService struct {
	userRepository *repository.UserRepository
	config         *config.Config
}

func NewAuthService(userRepository *repository.UserRepository, config *config.Config) *AuthService {
	return &AuthService{userRepository: userRepository, config: config}
}

func (s *AuthService) InitiateSteamAuth() (string, error) {
	callbackURL := s.config.Steam.Domain

	authURL, err := openid.RedirectURL("https://steamcommunity.com/openid", callbackURL, callbackURL) // callbackURL, realm string
	if err != nil {
		return "", err
	}
	return authURL, nil
}

func (s *AuthService) HandleSteamCallback(r *http.Request) (string, error) {
	claimedID, err := openid.Verify(r.URL.String(), openid.DiscoveryCache(nil), openid.NonceStore(nil))
	if err != nil {
		return "", err
	}

	steamIDRegex := regexp.MustCompile(`^https://steamcommunity\.com/openid/id/(\d+)$`)
	matches := steamIDRegex.FindStringSubmatch(claimedID)
	if len(matches) < 2 {
		return "", errors.New("cannot parse steam id from claimed id")
	}
	steamID := matches[1]

	user, err := s.userRepository.FindBySteamID(steamID)
	if err != nil {
		playerSummaries, err := s.GetPlayerSummaries(steamID)
		if err != nil {
			return "", err
		}

		user = &domain.User{
			SteamID:    steamID,
			Nickname:   playerSummaries.Nickname,
			AvatarURL:  playerSummaries.AvatarURL,
			ProfileURL: playerSummaries.ProfileURL,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err = s.userRepository.Create(user)
		if err != nil {
			return "", err
		}
	}

	// TODO: Generate JWT token for the user
	return "jwtTokenGoesHere", nil
}

func (s *AuthService) GetPlayerSummaries(steamID string) (*domain.User, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", s.config.Steam.SteamAPIKey, steamID)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Response struct {
			Players []struct {
				Nickname   string `json:"personaname"`
				AvatarURL  string `json:"avatarfull"`
				ProfileURL string `json:"profileurl"`
			} `json:"players"`
		} `json:"response"`
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Response.Players) == 0 {
		return nil, errors.New("no player data found for steam id")
	}
	playerData := resp.Response.Players[0]

	return &domain.User{
		Nickname:   playerData.Nickname,
		AvatarURL:  playerData.AvatarURL,
		ProfileURL: playerData.ProfileURL,
	}, nil
}
