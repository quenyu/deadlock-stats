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
	"github.com/quenyu/deadlock-stats/internal/dto"
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
	realm, err := s.extractRealmFromConfig()
	if err != nil {
		return "", err
	}

	params := s.buildSteamAuthParams(realm)
	return s.buildSteamAuthURL(params), nil
}

func (s *AuthService) extractRealmFromConfig() (string, error) {
	parsedURL, err := url.Parse(s.config.Steam.RedirectURL)
	if err != nil {
		s.logger.Error("Failed to parse steam domain from config", zap.Error(err))
		return "", err
	}
	return parsedURL.Scheme + "://" + parsedURL.Host, nil
}

func (s *AuthService) buildSteamAuthParams(realm string) url.Values {
	params := url.Values{}
	params.Add("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Add("openid.mode", "checkid_setup")
	params.Add("openid.return_to", s.config.Steam.RedirectURL+"/api/v1/auth/steam/callback")
	params.Add("openid.realm", realm)
	params.Add("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Add("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
	return params
}

func (s *AuthService) buildSteamAuthURL(params url.Values) string {
	return "https://steamcommunity.com/openid/login?" + params.Encode()
}

func (s *AuthService) HandleSteamCallback(r *http.Request) (string, error) {
	s.logger.Info("Handling Steam callback")

	if err := s.verifySteamAuthentication(r); err != nil {
		return "", err
	}

	steamID, err := s.extractSteamID(r)
	if err != nil {
		return "", err
	}

	user, err := s.findOrCreateUser(steamID)
	if err != nil {
		return "", err
	}

	return s.generateAndReturnToken(user.ID)
}

func (s *AuthService) verifySteamAuthentication(r *http.Request) error {
	params := r.URL.Query()
	params.Set("openid.mode", "check_authentication")

	reqURL := "https://steamcommunity.com/openid/login"
	resp, err := http.PostForm(reqURL, params)
	if err != nil {
		s.logger.Error("Failed to make verification request to Steam", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body from Steam", zap.Error(err))
		return err
	}

	if !strings.Contains(string(body), "is_valid:true") {
		err := errors.New("authentication failed: is_valid is not true")
		s.logger.Error("Steam authentication failed", zap.Error(err), zap.String("response", string(body)))
		return err
	}

	return nil
}

func (s *AuthService) extractSteamID(r *http.Request) (string, error) {
	params := r.URL.Query()
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
	return steamID, nil
}

func (s *AuthService) findOrCreateUser(steamID string) (*domain.User, error) {
	user, err := s.userRepository.FindBySteamID(steamID)
	if err != nil {
		return s.createNewUser(steamID)
	}

	s.logger.Info("Found existing user", zap.String("userID", user.ID.String()))
	return user, nil
}

func (s *AuthService) createNewUser(steamID string) (*domain.User, error) {
	s.logger.Warn("User not found, creating new user", zap.String("steamID", steamID))

	playerSummaries, err := s.GetPlayerSummaries(steamID)
	if err != nil {
		s.logger.Error("Failed to get player summaries from Steam API", zap.String("steamID", steamID), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Successfully got player summaries", zap.String("nickname", playerSummaries.Nickname))

	user := &domain.User{
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
		return nil, err
	}

	s.logger.Info("Successfully created new user", zap.String("userID", user.ID.String()))
	return user, nil
}

func (s *AuthService) generateAndReturnToken(userID uuid.UUID) (string, error) {
	jwtToken, err := s.GenerateJWTToken(userID)
	if err != nil {
		s.logger.Error("Failed to generate JWT token", zap.Error(err))
		return "", err
	}

	s.logger.Info("Successfully generated JWT token")
	return jwtToken, nil
}

func (s *AuthService) ResolveVanityURL(vanityURL string) (string, error) {
	url := s.buildSteamAPIURL("ResolveVanityURL", map[string]string{"vanityurl": vanityURL})

	body, err := s.makeSteamAPIRequest(url)
	if err != nil {
		return "", err
	}

	return s.parseVanityURLResponse(body)
}

func (s *AuthService) GetPlayerSummaries(steamID string) (*domain.User, error) {
	url := s.buildSteamAPIURL("GetPlayerSummaries", map[string]string{"steamids": steamID})

	body, err := s.makeSteamAPIRequest(url)
	if err != nil {
		return nil, err
	}

	return s.parsePlayerSummariesResponse(body, steamID)
}

func (s *AuthService) buildSteamAPIURL(method string, params map[string]string) string {
	baseURL := fmt.Sprintf("https://api.steampowered.com/ISteamUser/%s/v0001/", method)

	urlParams := url.Values{}
	urlParams.Add("key", s.config.Steam.APIKey)

	for key, value := range params {
		urlParams.Add(key, value)
	}

	return baseURL + "?" + urlParams.Encode()
}

func (s *AuthService) makeSteamAPIRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *AuthService) parsePlayerSummariesResponse(body []byte, steamID string) (*domain.User, error) {
	var resp struct {
		Response struct {
			Players []struct {
				Nickname   string `json:"personaname"`
				AvatarURL  string `json:"avatarfull"`
				ProfileURL string `json:"profileurl"`
			} `json:"players"`
		} `json:"response"`
	}

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Response.Players) > 0 {
		s.logger.Info("Steam API GetPlayerSummaries successful",
			zap.String("steamID", steamID),
			zap.String("personaname", resp.Response.Players[0].Nickname),
			zap.String("avatarfull", resp.Response.Players[0].AvatarURL),
		)
	} else {
		s.logger.Warn("Steam API GetPlayerSummaries returned zero players", zap.String("steamID", steamID))
	}

	if len(resp.Response.Players) == 0 {
		return nil, errors.New("no player data found for steam id")
	}

	playerData := resp.Response.Players[0]
	return s.createUserFromSteamData(playerData), nil
}

func (s *AuthService) parseVanityURLResponse(body []byte) (string, error) {
	var resp dto.ResolveVanityURLResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	if resp.Response.Success != 1 {
		return "", errors.New("vanity URL could not be resolved")
	}

	return resp.Response.SteamID, nil
}

func (s *AuthService) createUserFromSteamData(playerData struct {
	Nickname   string `json:"personaname"`
	AvatarURL  string `json:"avatarfull"`
	ProfileURL string `json:"profileurl"`
}) *domain.User {
	return &domain.User{
		Nickname:   playerData.Nickname,
		AvatarURL:  playerData.AvatarURL,
		ProfileURL: playerData.ProfileURL,
	}
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
