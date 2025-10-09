package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"sync"

	"github.com/google/uuid"
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/quenyu/deadlock-stats/internal/dto"
	"github.com/quenyu/deadlock-stats/internal/repositories"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type steamSearchCacheEntry struct {
	results   []domain.SteamProfileSearch
	timestamp time.Time
}

const steamSearchCacheTTL = 5 * time.Minute

type PlayerSearchService struct {
	playerProfileRepository *repositories.PlayerProfilePostgresRepository
	userRepository          *repositories.UserRepository
	authService             *AuthService
	deadlockAPIClient       *deadlockapi.Client
	logger                  *zap.Logger

	steamSearchCache      map[string]steamSearchCacheEntry
	steamSearchCacheMutex sync.Mutex

	redisClient *redis.Client
	steamAPIKey string
}

func NewPlayerSearchService(
	playerProfileRepository *repositories.PlayerProfilePostgresRepository,
	userRepository *repositories.UserRepository,
	authService *AuthService,
	deadlockAPIClient *deadlockapi.Client,
	redisClient *redis.Client,
	steamAPIKey string,
	logger *zap.Logger,
) *PlayerSearchService {
	return &PlayerSearchService{
		playerProfileRepository: playerProfileRepository,
		userRepository:          userRepository,
		authService:             authService,
		deadlockAPIClient:       deadlockAPIClient,
		redisClient:             redisClient,
		steamAPIKey:             steamAPIKey,
		logger:                  logger,
		steamSearchCache:        make(map[string]steamSearchCacheEntry),
	}
}

func (s *PlayerSearchService) SearchPlayers(ctx context.Context, query string, searchType string, page, pageSize int) (*dto.SearchResult, error) {
	if len(query) < 2 {
		return &dto.SearchResult{
			Results:    []dto.UserSearchResult{},
			TotalCount: 0,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: 0,
		}, nil
	}

	if strings.HasPrefix(query, "https://steamcommunity.com/id/") {
		vanity := strings.TrimPrefix(query, "https://steamcommunity.com/id/")
		vanity = strings.TrimSuffix(vanity, "/")
		steamID, err := resolveVanityURL(vanity, s.steamAPIKey)
		if err == nil && steamID != "" {
			query = steamID
		}
	}

	var results []dto.UserSearchResult
	var err error

	switch searchType {
	case "steamid":
		results, err = s.searchBySteamID(ctx, query)
	case "nickname":
		results, err = s.searchByNickname(ctx, query)
	default:
		localResults, apiResults, err := s.fetchSearchResults(ctx, query)
		if err != nil {
			return nil, err
		}
		results = s.combineSearchResults(localResults, apiResults)
	}

	if err != nil {
		return nil, err
	}

	return s.applyPagination(results, page, pageSize), nil
}

func (s *PlayerSearchService) SearchPlayersWithAutocomplete(ctx context.Context, query string, limit int) ([]dto.UserSearchResult, error) {
	if len(query) < 2 {
		return []dto.UserSearchResult{}, nil
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	if strings.HasPrefix(query, "https://steamcommunity.com/id/") {
		vanity := strings.TrimPrefix(query, "https://steamcommunity.com/id/")
		vanity = strings.TrimSuffix(vanity, "/")
		s.logger.Info("Attempting to resolve vanity URL", zap.String("vanity", vanity))
		steamID, err := resolveVanityURL(vanity, s.steamAPIKey)
		if err == nil && steamID != "" {
			s.logger.Info("Successfully resolved vanity URL", zap.String("vanity", vanity), zap.String("steamID", steamID))
			query = steamID
		} else {
			s.logger.Warn("Failed to resolve vanity URL", zap.String("vanity", vanity), zap.Error(err))
		}
	}

	if !s.isValidSteamID(query) && !strings.Contains(query, " ") {
		s.logger.Info("Attempting to resolve potential vanity URL", zap.String("query", query))
		steamID, err := resolveVanityURL(query, s.steamAPIKey)
		if err == nil && steamID != "" {
			s.logger.Info("Successfully resolved potential vanity URL", zap.String("query", query), zap.String("steamID", steamID))
			query = steamID
		} else {
			s.logger.Debug("Failed to resolve potential vanity URL", zap.String("query", query), zap.Error(err))
		}
	}

	s.logger.Info("Searching players with autocomplete", zap.String("query", query), zap.Int("limit", limit))

	nicknameResults, err := s.playerProfileRepository.SearchByNicknamePartial(ctx, query, limit)
	if err != nil {
		s.logger.Error("Error searching by partial nickname", zap.Error(err))
		return []dto.UserSearchResult{}, err
	}

	steamIDResults, err := s.playerProfileRepository.SearchBySteamIDPartial(ctx, query, limit)
	if err != nil {
		s.logger.Error("Error searching by partial Steam ID", zap.Error(err))
		return []dto.UserSearchResult{}, err
	}

	localResults := append(nicknameResults, steamIDResults...)

	apiResults, err := s.deadlockAPIClient.FetchSteamProfileSearch(query)
	if err != nil {
		s.logger.Warn("Failed to fetch from Deadlock API search", zap.Error(err))
		apiResults = []domain.SteamProfileSearch{}
	}

	combined := s.combineSearchResults(localResults, apiResults)

	if len(combined) > limit {
		combined = combined[:limit]
	}
	return combined, nil
}

func (s *PlayerSearchService) SearchPlayersWithFilters(ctx context.Context, query string, filters dto.SearchFilters, page, pageSize int) (*dto.SearchResult, error) {
	if len(query) < 2 {
		return &dto.SearchResult{
			Results:    []dto.UserSearchResult{},
			TotalCount: 0,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: 0,
		}, nil
	}

	if strings.HasPrefix(query, "https://steamcommunity.com/id/") {
		vanity := strings.TrimPrefix(query, "https://steamcommunity.com/id/")
		vanity = strings.TrimSuffix(vanity, "/")
		steamID, err := resolveVanityURL(vanity, s.steamAPIKey)
		if err == nil && steamID != "" {
			query = steamID
		}
	}

	var results []dto.UserSearchResult
	var err error

	switch filters.GetSearchType() {
	case "steamid":
		results, err = s.searchBySteamID(ctx, query)
	case "nickname":
		results, err = s.searchByNickname(ctx, query)
	default:
		localResults, apiResults, err := s.fetchSearchResults(ctx, query)
		if err != nil {
			return nil, err
		}
		results = s.combineSearchResults(localResults, apiResults)
	}

	if err != nil {
		return nil, err
	}

	s.sortUsersByFilters(results, filters)
	return s.applyPagination(results, page, pageSize), nil
}

func (s *PlayerSearchService) GetPopularPlayers(ctx context.Context, page, pageSize int) (*dto.SearchResult, error) {
	users, err := s.playerProfileRepository.GetPopularPlayers(ctx, 1000)
	if err != nil {
		s.logger.Error("Error getting popular players", zap.Error(err))
		return nil, err
	}

	validResults := make([]dto.UserSearchResult, 0, len(users))
	for _, user := range users {
		if user.SteamID != "" && user.Nickname != "" && user.AvatarURL != "" {
			validResults = append(validResults, dto.UserSearchResult{
				ID:         user.ID.String(),
				SteamID:    user.SteamID,
				Nickname:   user.Nickname,
				AvatarURL:  user.AvatarURL,
				ProfileURL: user.ProfileURL,
				CreatedAt:  &user.CreatedAt,
				UpdatedAt:  &user.UpdatedAt,
			})
		}
	}

	s.logger.Info("Retrieved popular players",
		zap.Int("found", len(validResults)))

	return s.applyPagination(validResults, page, pageSize), nil
}

func (s *PlayerSearchService) GetRecentlyActivePlayers(ctx context.Context, page, pageSize int) (*dto.SearchResult, error) {
	users, err := s.playerProfileRepository.GetRecentlyActivePlayers(ctx, 1000)
	if err != nil {
		s.logger.Error("Error getting recently active players", zap.Error(err))
		return nil, err
	}

	validResults := make([]dto.UserSearchResult, 0, len(users))
	for _, user := range users {
		if user.SteamID != "" && user.Nickname != "" && user.AvatarURL != "" {
			validResults = append(validResults, dto.UserSearchResult{
				ID:         user.ID.String(),
				SteamID:    user.SteamID,
				Nickname:   user.Nickname,
				AvatarURL:  user.AvatarURL,
				ProfileURL: user.ProfileURL,
				CreatedAt:  &user.CreatedAt,
				UpdatedAt:  &user.UpdatedAt,
			})
		}
	}

	s.logger.Info("Retrieved recently active players",
		zap.Int("found", len(validResults)))

	return s.applyPagination(validResults, page, pageSize), nil
}

func (s *PlayerSearchService) IsValidUser(user dto.UserSearchResult) bool {
	return user.SteamID != "" && user.Nickname != "" && user.AvatarURL != ""
}

func (s *PlayerSearchService) searchBySteamID(ctx context.Context, query string) ([]dto.UserSearchResult, error) {
	if !s.isValidSteamID(query) {
		return []dto.UserSearchResult{}, nil
	}

	user, err := s.authService.SearchPlayersBySteamID(query)
	if err != nil {
		s.logger.Error("Error finding user by SteamID", zap.Error(err), zap.String("steamID", query))
		return []dto.UserSearchResult{}, err
	}

	if user != nil {
		return []dto.UserSearchResult{{
			ID:         user.ID.String(),
			SteamID:    user.SteamID,
			Nickname:   user.Nickname,
			AvatarURL:  user.AvatarURL,
			ProfileURL: user.ProfileURL,
			CreatedAt:  &user.CreatedAt,
			UpdatedAt:  &user.UpdatedAt,
		}}, nil
	}

	return []dto.UserSearchResult{}, nil
}

func (s *PlayerSearchService) searchByNickname(ctx context.Context, query string) ([]dto.UserSearchResult, error) {
	localResults, err := s.playerProfileRepository.SearchByNickname(ctx, query)
	if err != nil {
		s.logger.Error("Error searching by nickname", zap.Error(err))
		return []dto.UserSearchResult{}, err
	}

	validResults := make([]dto.UserSearchResult, 0, len(localResults))
	for _, user := range localResults {
		if user.SteamID != "" && user.Nickname != "" && user.AvatarURL != "" {
			validResults = append(validResults, dto.UserSearchResult{
				ID:         user.ID.String(),
				SteamID:    user.SteamID,
				Nickname:   user.Nickname,
				AvatarURL:  user.AvatarURL,
				ProfileURL: user.ProfileURL,
				CreatedAt:  &user.CreatedAt,
				UpdatedAt:  &user.UpdatedAt,
			})
		}
	}

	return validResults, nil
}

func (s *PlayerSearchService) fetchSearchResults(ctx context.Context, query string) ([]domain.User, []domain.SteamProfileSearch, error) {
	localResults, err := s.playerProfileRepository.SearchByNickname(ctx, query)
	if err != nil {
		s.logger.Error("Error searching local database", zap.Error(err))
		return []domain.User{}, []domain.SteamProfileSearch{}, err
	}

	s.steamSearchCacheMutex.Lock()
	cacheEntry, found := s.steamSearchCache[query]
	if found && time.Since(cacheEntry.timestamp) < steamSearchCacheTTL {
		s.steamSearchCacheMutex.Unlock()
		return localResults, cacheEntry.results, nil
	}
	s.steamSearchCacheMutex.Unlock()

	apiResults, err := s.deadlockAPIClient.FetchSteamProfileSearch(query)
	if err != nil {
		s.logger.Warn("Failed to fetch from Deadlock API search", zap.Error(err))
		apiResults = []domain.SteamProfileSearch{}
	}

	s.steamSearchCacheMutex.Lock()
	s.steamSearchCache[query] = steamSearchCacheEntry{
		results:   apiResults,
		timestamp: time.Now(),
	}
	s.steamSearchCacheMutex.Unlock()

	return localResults, apiResults, nil
}

func (s *PlayerSearchService) combineSearchResults(localResults []domain.User, apiResults []domain.SteamProfileSearch) []dto.UserSearchResult {
	combinedUsers := make(map[string]dto.UserSearchResult)

	for _, u := range localResults {
		combinedUsers[u.SteamID] = dto.UserSearchResult{
			ID:         u.ID.String(),
			SteamID:    u.SteamID,
			Nickname:   u.Nickname,
			AvatarURL:  u.AvatarURL,
			ProfileURL: u.ProfileURL,
			CreatedAt:  &u.CreatedAt,
			UpdatedAt:  &u.UpdatedAt,
		}
	}

	for _, apiPlayer := range apiResults {
		if !s.isValidAPIPlayer(apiPlayer) {
			s.logger.Debug("Skipping invalid API player",
				zap.Int("accountID", apiPlayer.AccountID),
				zap.String("personaname", apiPlayer.Personaname),
				zap.String("avatar", apiPlayer.Avatar))
			continue
		}

		steamID64 := s.convertAccountIDToSteamID64(apiPlayer.AccountID)
		if steamID64 == "" {
			s.logger.Warn("Failed to convert AccountID to SteamID64",
				zap.Int("accountID", apiPlayer.AccountID))
			continue
		}

		if _, exists := combinedUsers[steamID64]; !exists {
			combinedUsers[steamID64] = dto.UserSearchResult{
				ID:                  "",
				SteamID:             steamID64,
				Nickname:            apiPlayer.Personaname,
				AvatarURL:           apiPlayer.Avatar,
				ProfileURL:          apiPlayer.Profileurl,
				AccountID:           apiPlayer.AccountID,
				CountryCode:         apiPlayer.CountryCode,
				LastUpdated:         apiPlayer.LastUpdated,
				Realname:            apiPlayer.Realname,
				IsDeadlockPlayer:    true,
				DeadlockStatusKnown: true,
			}
		}
	}

	result := make([]dto.UserSearchResult, 0, len(combinedUsers))
	for _, user := range combinedUsers {
		if !user.IsDeadlockPlayer {
			isActive, known := s.hasDeadlockActivity(user.SteamID, false)
			user.IsDeadlockPlayer = isActive
			user.DeadlockStatusKnown = known
		}
		result = append(result, user)
	}

	return result
}

func (s *PlayerSearchService) hasDeadlockActivity(steamID string, isFromDeadlockAPI bool) (bool, bool) {
	if isFromDeadlockAPI {
		return true, true
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("deadlock-activity:%s", steamID)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		return cached == "true", true
	}

	var count int64
	err = s.userRepository.GetDB().Table("users").
		Joins("JOIN player_stats ps ON users.id = ps.user_id").
		Where("users.steam_id = ? AND ps.total_matches > 0", steamID).
		Count(&count).Error

	if err != nil {
		s.logger.Warn("Error checking local Deadlock activity",
			zap.String("steamID", steamID),
			zap.Error(err))
		return false, false
	}

	if count > 0 {
		s.redisClient.Set(ctx, cacheKey, "true", 6*time.Hour)
		return true, true
	}

	var userExists int64
	err = s.userRepository.GetDB().Table("users").
		Where("steam_id = ?", steamID).
		Count(&userExists).Error

	if err == nil && userExists > 0 {
		s.redisClient.Set(ctx, cacheKey, "false", 6*time.Hour)
		return false, true
	}

	return false, true
}

func (s *PlayerSearchService) isValidAPIPlayer(apiPlayer domain.SteamProfileSearch) bool {
	if apiPlayer.AccountID <= 0 {
		return false
	}

	if apiPlayer.Personaname == "" {
		return false
	}

	return true
}

func (s *PlayerSearchService) convertAccountIDToSteamID64(accountID int) string {
	steamID64 := int64(accountID) + 76561197960265728
	return strconv.FormatInt(steamID64, 10)
}

func (s *PlayerSearchService) createUserFromAPISearch(apiPlayer domain.SteamProfileSearch, steamID64 string) domain.User {
	avatarURL := apiPlayer.Avatar
	if avatarURL == "" {
		avatarURL = "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb.jpg"
	}

	profileURL := apiPlayer.Profileurl
	if profileURL == "" {
		profileURL = fmt.Sprintf("https://steamcommunity.com/profiles/%s", steamID64)
	}

	return domain.User{
		ID:         uuid.New(),
		SteamID:    steamID64,
		Nickname:   apiPlayer.Personaname,
		AvatarURL:  avatarURL,
		ProfileURL: profileURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (s *PlayerSearchService) findAndCreateUser(steamID string) ([]domain.User, error) {
	steamUser, err := s.authService.GetPlayerSummaries(steamID)
	if err != nil {
		s.logger.Error("Error getting player summaries", zap.Error(err), zap.String("steamID", steamID))
		return []domain.User{}, err
	}

	if steamUser == nil {
		return []domain.User{}, nil
	}

	user := s.createUserFromSteamData(steamID, steamUser)
	err = s.userRepository.FindOrCreate(&user)
	if err != nil {
		s.logger.Error("Error creating user", zap.Error(err), zap.String("steamID", steamID))
		return []domain.User{}, err
	}

	return []domain.User{user}, nil
}

func (s *PlayerSearchService) createUserFromSteamData(steamID string, steamUser *domain.User) domain.User {
	return domain.User{
		SteamID:    steamID,
		Nickname:   steamUser.Nickname,
		AvatarURL:  steamUser.AvatarURL,
		ProfileURL: steamUser.ProfileURL,
	}
}

func (s *PlayerSearchService) isValidSteamID(steamID string) bool {
	if len(steamID) != 17 {
		return false
	}
	for _, char := range steamID {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (s *PlayerSearchService) sortUsersByFilters(users []dto.UserSearchResult, filters dto.SearchFilters) {
	sort.Slice(users, func(i, j int) bool {
		switch filters.GetDefaultSortBy() {
		case "created_at":
			if users[i].CreatedAt == nil && users[j].CreatedAt == nil {
				return false
			}
			if users[i].CreatedAt == nil {
				return filters.GetDefaultSortOrder() == "asc"
			}
			if users[j].CreatedAt == nil {
				return filters.GetDefaultSortOrder() == "desc"
			}
			if filters.GetDefaultSortOrder() == "desc" {
				return users[i].CreatedAt.After(*users[j].CreatedAt)
			}
			return users[i].CreatedAt.Before(*users[j].CreatedAt)
		case "updated_at":
			if users[i].UpdatedAt == nil && users[j].UpdatedAt == nil {
				return false
			}
			if users[i].UpdatedAt == nil {
				return filters.GetDefaultSortOrder() == "asc"
			}
			if users[j].UpdatedAt == nil {
				return filters.GetDefaultSortOrder() == "desc"
			}
			if filters.GetDefaultSortOrder() == "desc" {
				return users[i].UpdatedAt.After(*users[j].UpdatedAt)
			}
			return users[i].UpdatedAt.Before(*users[j].UpdatedAt)
		default: // "nickname"
			if filters.GetDefaultSortOrder() == "desc" {
				return users[i].Nickname > users[j].Nickname
			}
			return users[i].Nickname < users[j].Nickname
		}
	})
}

func (s *PlayerSearchService) applyPagination(results []dto.UserSearchResult, page, pageSize int) *dto.SearchResult {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	totalCount := len(results)
	totalPages := (totalCount + pageSize - 1) / pageSize

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= totalCount {
		return &dto.SearchResult{
			Results:    []dto.UserSearchResult{},
			TotalCount: totalCount,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		}
	}

	if end > totalCount {
		end = totalCount
	}

	return &dto.SearchResult{
		Results:    results[start:end],
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func resolveVanityURL(vanity string, apiKey string) (string, error) {
	if apiKey == "" {
		return "", errors.New("steam API key is empty")
	}

	url := "https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=" + apiKey + "&vanityurl=" + vanity
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			SteamID string `json:"steamid"`
			Success int    `json:"success"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Response.Success != 1 {
		return "", errors.New("not found")
	}
	return result.Response.SteamID, nil
}
