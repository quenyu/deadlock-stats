package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/quenyu/deadlock-stats/internal/config"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/quenyu/deadlock-stats/internal/repositories"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepository *repositories.UserRepository
	config         *config.Config
	logger         *zap.Logger
}

func NewAuthService(userRepository *repositories.UserRepository, config *config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		config:         config,
		logger:         logger.Named("AuthService"),
	}
}

func (s *AuthService) InitiateSteamAuth() (string, error) {
	parsedURL, err := url.Parse(s.config.Steam.Domain)
	if err != nil {
		s.logger.Error("Failed to parse steam domain from config", zap.Error(err))
		return "", err
	}
	realm := parsedURL.Scheme + "://" + parsedURL.Host

	params := url.Values{}
	params.Add("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Add("openid.mode", "checkid_setup")
	params.Add("openid.return_to", s.config.Steam.Domain + "/api/v1/auth/steam/callback")
	params.Add("openid.realm", realm)
	params.Add("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Add("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")

	return "https://steamcommunity.com/openid/login?" + params.Encode(), nil
}

func (s *AuthService) HandleSteamCallback(r *http.Request) (string, error) {
	s.logger.Info("Handling Steam callback")

	params := r.URL.Query()
	params.Set("openid.mode", "check_authentication")

	reqURL := "https://steamcommunity.com/openid/login"
	resp, err := http.PostForm(reqURL, params)
	if err != nil {
		s.logger.Error("Failed to make verification request to Steam", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body from Steam", zap.Error(err))
		return "", err
	}

	if !strings.Contains(string(body), "is_valid:true") {
		err := errors.New("authentication failed: is_valid is not true")
		s.logger.Error("Steam authentication failed", zap.Error(err), zap.String("response", string(body)))
		return "", err
	}

	claimedID := params.Get("openid.claimed_id")
	steamIDRegex := regexp.MustCompile(`^https://steamcommunity\.com/openid/id/(\d+)$`)
	matches := steamIDRegex.FindStringSubmatch(claimedID)
	if len(matches) < 2 {
		err := errors.New("cannot parse steam id from claimed id")
		s.logger.Error("Failed to parse SteamID", zap.String("claimedID", claimedID))
		return "", err
	}
	steamID := matches[1]
	s.logger.Info("SteamID verified", zap.String("steamID", steamID))

	user, err := s.userRepository.FindBySteamID(steamID)
	if err != nil {
		s.logger.Warn("User not found, creating new user", zap.String("steamID", steamID), zap.Error(err))
		playerSummaries, err := s.GetPlayerSummaries(steamID)
		if err != nil {
			s.logger.Error("Failed to get player summaries from Steam API", zap.String("steamID", steamID), zap.Error(err))
			return "", err
		}
		s.logger.Info("Successfully got player summaries", zap.String("nickname", playerSummaries.Nickname))

		user = &domain.User{
			ID:         uuid.New(),
			SteamID:    steamID,
			Nickname:   playerSummaries.Nickname,
			AvatarURL:  playerSummaries.AvatarURL,
			ProfileURL: playerSummaries.ProfileURL,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err = s.userRepository.Create(user)
		if err != nil {
			s.logger.Error("Failed to create user in database", zap.Error(err))
			return "", err
		}
		s.logger.Info("Successfully created new user", zap.String("userID", user.ID.String()))
	} else {
		s.logger.Info("Found existing user", zap.String("userID", user.ID.String()))
	}

	jwtToken, err := s.GenerateJWTToken(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate JWT token", zap.Error(err))
		return "", err
	}
	s.logger.Info("Successfully generated JWT token")

	return jwtToken, nil
}

func (s *AuthService) GetPlayerSummaries(steamID string) (*domain.User, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", s.config.Steam.SteamAPIKey, steamID)

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

func (s *AuthService) GenerateJWTToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.config.JWT.Expiration).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *AuthService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepository.FindByID(id)
}
