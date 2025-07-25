package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/quenyu/deadlock-stats/internal/dto"
	"github.com/quenyu/deadlock-stats/internal/repositories"
	"go.uber.org/zap"
)

type PlayerProfileService struct {
	playerProfileRepository *repositories.PlayerProfilePostgresRepository
	userRepository          *repositories.UserRepository
	authService             *AuthService
	deadlockAPIClient       *deadlockapi.Client
	staticDataService       *StaticDataService
	redisClient             *redis.Client
	logger                  *zap.Logger
}

func NewPlayerProfileService(
	playerProfileRepository *repositories.PlayerProfilePostgresRepository,
	userRepository *repositories.UserRepository,
	authService *AuthService,
	deadlockAPIClient *deadlockapi.Client,
	staticDataService *StaticDataService,
	redisClient *redis.Client,
	logger *zap.Logger,
) *PlayerProfileService {
	return &PlayerProfileService{
		playerProfileRepository: playerProfileRepository,
		userRepository:          userRepository,
		authService:             authService,
		deadlockAPIClient:       deadlockAPIClient,
		staticDataService:       staticDataService,
		redisClient:             redisClient,
		logger:                  logger,
	}
}

func (s *PlayerProfileService) fetchFromCacheOrAPI(ctx context.Context, steamID string) (*dto.ExtendedPlayerProfile, error) {
	cacheKey := fmt.Sprintf("player-profile:%s", steamID)
	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var profile dto.ExtendedPlayerProfile
		if err := json.Unmarshal([]byte(val), &profile); err == nil {
			s.logger.Info("cache hit for player profile", zap.String("steamID", steamID))
			return &profile, nil
		}
	}

	s.logger.Info("cache miss for player profile", zap.String("steamID", steamID))

	partialProfile := s.tryGetPartialCache(ctx, steamID)

	fullProfile, err := s.fetchAndBuildProfile(ctx, steamID, cacheKey)
	if err != nil {
		if partialProfile != nil {
			s.logger.Warn("Using partial cached profile due to API error", zap.String("steamID", steamID), zap.Error(err))
			return partialProfile, nil
		}
		return nil, err
	}

	s.cacheProfile(ctx, cacheKey, fullProfile)

	s.updatePartialCache(ctx, steamID, fullProfile)

	return fullProfile, nil
}

func (s *PlayerProfileService) tryGetPartialCache(ctx context.Context, steamID string) *dto.ExtendedPlayerProfile {
	partialCacheKey := fmt.Sprintf("player-profile-partial:%s", steamID)

	val, err := s.redisClient.Get(ctx, partialCacheKey).Result()
	if err != nil {
		return nil
	}

	var profile dto.ExtendedPlayerProfile
	if err := json.Unmarshal([]byte(val), &profile); err == nil {
		if time.Since(profile.LastUpdatedAt) < time.Hour {
			s.logger.Info("using partial cache", zap.String("steamID", steamID))
			return &profile
		}
	}

	return nil
}

func (s *PlayerProfileService) updatePartialCache(ctx context.Context, steamID string, profile *dto.ExtendedPlayerProfile) {
	partialCacheKey := fmt.Sprintf("player-profile-partial:%s", steamID)

	partialProfile := &dto.ExtendedPlayerProfile{
		Card:          profile.Card,
		Nickname:      profile.Nickname,
		AvatarURL:     profile.AvatarURL,
		PlayerRank:    profile.PlayerRank,
		RankName:      profile.RankName,
		RankImage:     profile.RankImage,
		TotalMatches:  profile.TotalMatches,
		WinRate:       profile.WinRate,
		KDRatio:       profile.KDRatio,
		LastUpdatedAt: time.Now(),
	}

	data, err := json.Marshal(partialProfile)
	if err == nil {
		s.redisClient.Set(ctx, partialCacheKey, data, time.Hour).Err()
	}
}

func (s *PlayerProfileService) fetchAndBuildProfile(ctx context.Context, steamID, cacheKey string) (*dto.ExtendedPlayerProfile, error) {
	start := time.Now()
	defer func() {
		s.logger.Debug("Profile building completed",
			zap.String("steamID", steamID),
			zap.Duration("buildTime", time.Since(start)))
	}()

	card, matches, heroStats, mmrHistory, profile, heroMMRHistory, err := s.fetchAllData(ctx, steamID)
	if err != nil {
		return nil, err
	}

	extendedProfile := s.buildExtendedProfile(card, matches, heroStats, mmrHistory, profile, heroMMRHistory)

	s.cacheProfile(ctx, cacheKey, extendedProfile)

	return extendedProfile, nil
}

func (s *PlayerProfileService) fetchAllData(ctx context.Context, steamID string) (
	*deadlockapi.DeadlockCard,
	[]deadlockapi.DeadlockMatch,
	[]domain.HeroStat,
	[]domain.DeadlockMMR,
	*domain.PlayerProfile,
	[]domain.HeroMMRHistory,
	error,
) {
	apiCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var card *deadlockapi.DeadlockCard
	var matches []deadlockapi.DeadlockMatch
	var heroStats []domain.HeroStat
	var mmrHistory []domain.DeadlockMMR
	var profile *domain.PlayerProfile
	var heroMMRHistory []domain.HeroMMRHistory

	var wg sync.WaitGroup
	errs := make(chan error, 5)

	wg.Add(5)
	go func() {
		defer wg.Done()
		profile, _ = s.playerProfileRepository.FindBySteamID(ctx, steamID)
	}()
	go func() {
		defer wg.Done()
		var err error
		card, err = s.deadlockAPIClient.FetchPlayerCard(steamID)
		if err != nil {
			errs <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		matches, err = s.deadlockAPIClient.FetchMatchHistory(steamID)
		if err != nil {
			errs <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		heroStats, err = s.deadlockAPIClient.FetchHeroStats(steamID)
		if err != nil {
			errs <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		mmrHistory, err = s.deadlockAPIClient.FetchMMRHistory(steamID)
		if err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	for err := range errs {
		s.logger.Warn("Failed to fetch some data from Deadlock API", zap.Error(err))
	}

	if card == nil || matches == nil || heroStats == nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("critical data could not be fetched (card, matches, or hero stats)")
	}

	if profile == nil {
		profile = &domain.PlayerProfile{SteamID: steamID}
	}

	var processWg sync.WaitGroup
	processWg.Add(2)

	go func() {
		defer processWg.Done()
		heroStats = s.processHeroStats(heroStats)
	}()

	go func() {
		defer processWg.Done()
		sort.SliceStable(mmrHistory, func(i, j int) bool {
			return mmrHistory[i].StartTime > mmrHistory[j].StartTime
		})
	}()

	processWg.Wait()

	heroMMRHistory = s.fetchHeroMMRHistoryWithTimeout(apiCtx, steamID, heroStats)

	return card, matches, heroStats, mmrHistory, profile, heroMMRHistory, nil
}

func (s *PlayerProfileService) fetchHeroMMRHistoryWithTimeout(ctx context.Context, steamID string, heroStats []domain.HeroStat) []domain.HeroMMRHistory {
	select {
	case <-ctx.Done():
		s.logger.Info("Skipping hero MMR history due to timeout", zap.String("steamID", steamID))
		return []domain.HeroMMRHistory{}
	default:
		return s.fetchHeroMMRHistory(steamID, heroStats)
	}
}

func (s *PlayerProfileService) processHeroStats(heroStats []domain.HeroStat) []domain.HeroStat {
	sort.Slice(heroStats, func(i, j int) bool {
		return heroStats[i].Matches > heroStats[j].Matches
	})

	s.enrichHeroStats(heroStats)
	return heroStats
}

func (s *PlayerProfileService) fetchHeroMMRHistory(steamID string, heroStats []domain.HeroStat) []domain.HeroMMRHistory {
	var topHeroes []domain.HeroStat
	if len(heroStats) > 5 {
		topHeroes = heroStats[:5]
	} else {
		topHeroes = heroStats
	}

	var heroWg sync.WaitGroup
	heroMMRHistoryChan := make(chan domain.HeroMMRHistory, len(topHeroes))
	heroErrs := make(chan error, len(topHeroes))

	for _, hero := range topHeroes {
		heroWg.Add(1)
		go func(hero domain.HeroStat) {
			defer heroWg.Done()
			history, err := s.deadlockAPIClient.FetchMMRHistoryByHero(steamID, hero.HeroID)
			if err != nil {
				heroErrs <- fmt.Errorf("failed to fetch MMR history for hero %d: %w", hero.HeroID, err)
				return
			}
			heroMMRHistoryChan <- domain.HeroMMRHistory{
				HeroID:   hero.HeroID,
				HeroName: hero.HeroName,
				History:  history,
			}
		}(hero)
	}

	heroWg.Wait()
	close(heroMMRHistoryChan)
	close(heroErrs)

	for err := range heroErrs {
		s.logger.Warn("Failed to fetch hero MMR history", zap.Error(err))
	}

	var heroMMRHistory []domain.HeroMMRHistory
	for history := range heroMMRHistoryChan {
		heroMMRHistory = append(heroMMRHistory, history)
	}

	return heroMMRHistory
}

func (s *PlayerProfileService) buildExtendedProfile(
	card *deadlockapi.DeadlockCard,
	matches []deadlockapi.DeadlockMatch,
	heroStats []domain.HeroStat,
	mmrHistory []domain.DeadlockMMR,
	profile *domain.PlayerProfile,
	heroMMRHistory []domain.HeroMMRHistory,
) *dto.ExtendedPlayerProfile {
	if card == nil {
		return nil
	}

	domainMatches := s.buildDomainMatches(matches, mmrHistory)
	s.calculateAndFillStats(profile, card, domainMatches, mmrHistory)

	featuredHeroes := s.enrichFeaturedHeroes(card)
	peakRank, peakRankName, peakRankImage := domain.FindPeakRank(mmrHistory, s.getRankNameAndSubRank, s.getRankImageURL)
	personalRecords := domain.CalculatePersonalRecords(domainMatches)
	avgStats := domain.CalculateAverageStats(domainMatches, len(matches))

	validMateStats := s.buildMateStats(card)

	dtoRecords := s.buildPersonalRecordsDTO(personalRecords)
	dtoMMRHistory := s.buildMMRHistoryDTO(mmrHistory)
	dtoHeroMMR := s.buildHeroMMRHistoryDTO(heroMMRHistory)

	return &dto.ExtendedPlayerProfile{
		Card:                card,
		MatchHistory:        domainMatches,
		HeroStats:           heroStats,
		MMRHistory:          dtoMMRHistory,
		TotalMatches:        profile.TotalMatches,
		WinRate:             profile.WinRate,
		KDRatio:             profile.KDRatio,
		PerformanceDynamics: profile.PerformanceDynamics,
		PlayerRank:          profile.PlayerRank,
		Nickname:            profile.Nickname,
		AvatarURL:           profile.AvatarURL,
		RankImage:           profile.RankImage,
		RankName:            profile.RankName,
		SubRank:             profile.SubRank,
		AvgSoulsPerMin:      profile.AvgSoulsPerMin,
		FeaturedHeroes:      featuredHeroes,
		PeakRank:            peakRank,
		PeakRankName:        peakRankName,
		PeakRankImage:       peakRankImage,
		PersonalRecords:     dtoRecords,
		AvgKillsPerMatch:    avgStats.AvgKills,
		AvgDeathsPerMatch:   avgStats.AvgDeaths,
		AvgAssistsPerMatch:  avgStats.AvgAssists,
		AvgMatchDuration:    avgStats.AvgDuration,
		MateStats:           validMateStats,
		HeroMMRHistory:      dtoHeroMMR,
		LastUpdatedAt:       time.Now(),
	}
}

func (s *PlayerProfileService) cacheProfile(ctx context.Context, cacheKey string, profile *dto.ExtendedPlayerProfile) {
	data, err := json.Marshal(profile)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, data, time.Hour).Err()
	}
}

func (s *PlayerProfileService) buildPersonalRecordsDTO(personalRecords domain.PersonalRecords) domain.PersonalRecords {
	return domain.PersonalRecords{
		MaxKills:           personalRecords.MaxKills,
		MaxAssists:         personalRecords.MaxAssists,
		MaxNetWorth:        personalRecords.MaxNetWorth,
		BestKDA:            personalRecords.BestKDA,
		MaxKillsMatchID:    personalRecords.MaxKillsMatchID,
		MaxAssistsMatchID:  personalRecords.MaxAssistsMatchID,
		MaxNetWorthMatchID: personalRecords.MaxNetWorthMatchID,
		BestKDAMatchID:     personalRecords.BestKDAMatchID,
	}
}

func (s *PlayerProfileService) buildMMRHistoryDTO(mmrHistory []domain.DeadlockMMR) []domain.DeadlockMMR {
	dtoMMRHistory := make([]domain.DeadlockMMR, len(mmrHistory))
	for i, mmr := range mmrHistory {
		dtoMMRHistory[i] = domain.DeadlockMMR{
			MatchID:      mmr.MatchID,
			Rank:         mmr.Rank,
			StartTime:    mmr.StartTime,
			PlayerScore:  mmr.PlayerScore,
			Division:     mmr.Division,
			DivisionTier: mmr.DivisionTier,
		}
	}
	return dtoMMRHistory
}

func (s *PlayerProfileService) buildHeroMMRHistoryDTO(heroMMRHistory []domain.HeroMMRHistory) []domain.HeroMMRHistory {
	dtoHeroMMR := make([]domain.HeroMMRHistory, len(heroMMRHistory))
	for i, hero := range heroMMRHistory {
		history := make([]domain.DeadlockMMR, len(hero.History))
		for j, h := range hero.History {
			history[j] = domain.DeadlockMMR{
				MatchID:      h.MatchID,
				Rank:         h.Rank,
				StartTime:    h.StartTime,
				PlayerScore:  h.PlayerScore,
				Division:     h.Division,
				DivisionTier: h.DivisionTier,
			}
		}
		dtoHeroMMR[i] = domain.HeroMMRHistory{
			HeroID:   hero.HeroID,
			HeroName: hero.HeroName,
			History:  history,
		}
	}
	return dtoHeroMMR
}

func (s *PlayerProfileService) buildMateStats(card *deadlockapi.DeadlockCard) []domain.MateStat {
	// This method would contain the logic for fetching mate stats
	// For now, returning empty slice as placeholder
	return []domain.MateStat{}
}

// GetExtendedPlayerProfile retrieves comprehensive player profile data including match history,
// hero statistics, MMR history, and personal records. Data is cached in Redis for performance.
func (s *PlayerProfileService) GetExtendedPlayerProfile(ctx context.Context, steamID string) (*dto.ExtendedPlayerProfile, error) {
	start := time.Now()
	defer func() {
		s.logger.Info("Profile loading completed",
			zap.String("steamID", steamID),
			zap.Duration("totalTime", time.Since(start)))
	}()

	profile, err := s.fetchFromCacheOrAPI(ctx, steamID)
	if err != nil {
		s.logger.Error("Failed to load profile",
			zap.String("steamID", steamID),
			zap.Error(err),
			zap.Duration("totalTime", time.Since(start)))
		return nil, err
	}

	s.logger.Info("Profile loaded successfully",
		zap.String("steamID", steamID),
		zap.Duration("totalTime", time.Since(start)),
		zap.Int("totalMatches", profile.TotalMatches),
		zap.Int("heroStatsCount", len(profile.HeroStats)))

	return profile, nil
}

func (s *PlayerProfileService) buildDomainMatches(matches []deadlockapi.DeadlockMatch, mmrHistory []domain.DeadlockMMR) []domain.Match {
	rankMap := s.createRankMap(mmrHistory)
	s.sortMatchesByTime(matches)

	domainMatches := s.convertToDomainMatches(matches, rankMap)
	s.enrichMatchesWithHeroData(domainMatches)
	s.enrichMatchesWithRankData(domainMatches, rankMap)
	s.calculateRankChanges(domainMatches)

	return domainMatches
}

func (s *PlayerProfileService) createRankMap(mmrHistory []domain.DeadlockMMR) map[int64]domain.DeadlockMMR {
	rankMap := make(map[int64]domain.DeadlockMMR)
	for _, r := range mmrHistory {
		rankMap[r.MatchID] = r
	}
	return rankMap
}

func (s *PlayerProfileService) sortMatchesByTime(matches []deadlockapi.DeadlockMatch) {
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].StartTime > matches[j].StartTime
	})
}

func (s *PlayerProfileService) convertToDomainMatches(matches []deadlockapi.DeadlockMatch, rankMap map[int64]domain.DeadlockMMR) []domain.Match {
	domainMatches := make([]domain.Match, len(matches))

	for i, match := range matches {
		domainMatches[i] = domain.Match{
			ID:             strconv.FormatInt(match.MatchID, 10),
			HeroID:         match.HeroID,
			PlayerKills:    match.PlayerKills,
			PlayerDeaths:   match.PlayerDeaths,
			PlayerAssists:  match.PlayerAssists,
			NetWorth:       match.NetWorth,
			MatchDurationS: match.MatchDurationS,
			MatchResult:    match.MatchResult,
			StartTime:      match.StartTime,
			MatchTime:      time.Unix(match.StartTime, 0),
			Result:         domain.MapMatchResult(match.MatchResult),
		}
	}

	return domainMatches
}

func (s *PlayerProfileService) enrichMatchesWithHeroData(domainMatches []domain.Match) {
	for i := range domainMatches {
		if hero, ok := s.staticDataService.HeroesByHeroID[domainMatches[i].HeroID]; ok {
			domainMatches[i].HeroName = hero.Name
			if hero.Images.IconHeroCard != nil {
				domainMatches[i].HeroAvatar = *hero.Images.IconHeroCard
			}
		}
	}
}

func (s *PlayerProfileService) enrichMatchesWithRankData(domainMatches []domain.Match, rankMap map[int64]domain.DeadlockMMR) {
	for i := range domainMatches {
		matchID, _ := strconv.ParseInt(domainMatches[i].ID, 10, 64)
		if mmr, ok := rankMap[matchID]; ok {
			calculatedRank := domain.GetRankFromScore(mmr.PlayerScore)
			domainMatches[i].PlayerRankAfterMatch = calculatedRank
			tier := calculatedRank / 10
			subTier := calculatedRank % 10
			domainMatches[i].RankName, _, domainMatches[i].RankImage = s.getRankNameAndSubRank(tier)
			domainMatches[i].SubRank = subTier
		}
	}
}

func (s *PlayerProfileService) calculateRankChanges(domainMatches []domain.Match) {
	for i := 0; i < len(domainMatches); i++ {
		if domainMatches[i].PlayerRankAfterMatch == 0 {
			continue
		}

		rankBefore := s.findPreviousRank(domainMatches, i)
		if rankBefore > 0 {
			domainMatches[i].PlayerRankChange = domainMatches[i].PlayerRankAfterMatch - rankBefore
		}
	}
}

func (s *PlayerProfileService) findPreviousRank(domainMatches []domain.Match, currentIndex int) int {
	for j := currentIndex + 1; j < len(domainMatches); j++ {
		if domainMatches[j].PlayerRankAfterMatch != 0 {
			return domainMatches[j].PlayerRankAfterMatch
		}
	}
	return 0
}

func (s *PlayerProfileService) enrichHeroStats(heroStats []domain.HeroStat) {
	for i := range heroStats {
		if hero, ok := s.staticDataService.HeroesByHeroID[heroStats[i].HeroID]; ok {
			if hero.Name != "" && !strings.HasPrefix(hero.Name, "HeroID_") {
				heroStats[i].HeroName = hero.Name
			} else {
				heroStats[i].HeroName = fmt.Sprintf("Hero %d", heroStats[i].HeroID)
				s.logger.Warn("Bad hero name from static data, using fallback",
					zap.Int("heroID", heroStats[i].HeroID),
					zap.String("originalName", hero.Name))
			}

			if hero.Images.IconHeroCard != nil {
				heroStats[i].HeroAvatar = *hero.Images.IconHeroCard
			}
		} else {
			heroStats[i].HeroName = fmt.Sprintf("Hero %d", heroStats[i].HeroID)
			s.logger.Error("Hero not found in static data", zap.Int("heroID", heroStats[i].HeroID))
		}
	}
}

func (s *PlayerProfileService) GetRecentMatches(ctx context.Context, steamID string) ([]domain.Match, error) {
	return s.playerProfileRepository.FindRecentMatchesBySteamID(ctx, steamID, 5)
}

func (s *PlayerProfileService) SearchPlayers(ctx context.Context, query string, searchType string) ([]domain.User, error) {
	switch searchType {
	case "steamid":
		return s.searchBySteamID(ctx, query)
	case "nickname":
		return s.searchByNickname(ctx, query)
	default:
		return nil, fmt.Errorf("invalid search type: %s", searchType)
	}
}

func (s *PlayerProfileService) searchBySteamID(ctx context.Context, query string) ([]domain.User, error) {
	if _, err := strconv.ParseInt(query, 10, 64); err == nil {
		users, err := s.findAndCreateUser(query)
		if err != nil {
			s.logger.Error("Error finding user by SteamID64", zap.Error(err), zap.String("steamID", query))
			return []domain.User{}, nil
		}
		if len(users) > 0 {
			return users, nil
		}
	}

	resolvedSteamID, err := s.authService.ResolveVanityURL(query)
	if err == nil && resolvedSteamID != "" {
		users, err := s.findAndCreateUser(resolvedSteamID)
		if err != nil {
			s.logger.Error("Error finding user by resolved vanity URL", zap.Error(err), zap.String("vanityURL", query), zap.String("resolvedID", resolvedSteamID))
			return []domain.User{}, nil
		}
		return users, nil
	}

	return []domain.User{}, nil
}

func (s *PlayerProfileService) findAndCreateUser(steamID string) ([]domain.User, error) {
	localUser, _ := s.userRepository.FindBySteamID(steamID)
	if localUser != nil {
		return []domain.User{*localUser}, nil
	}

	steamUser, err := s.authService.GetPlayerSummaries(steamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player summaries for steamID %s: %w", steamID, err)
	}

	if steamUser != nil {
		s.logger.Info("Data from GetPlayerSummaries to be saved",
			zap.String("steamID", steamID),
			zap.String("nickname", steamUser.Nickname),
			zap.String("avatar", steamUser.AvatarURL),
		)

		newUser := s.createUserFromSteamData(steamID, steamUser)
		if err := s.userRepository.FindOrCreate(&newUser); err != nil {
			return nil, fmt.Errorf("failed to create user for steamID %s: %w", steamID, err)
		}
		return []domain.User{newUser}, nil
	}

	return []domain.User{}, nil
}

func (s *PlayerProfileService) createUserFromSteamData(steamID string, steamUser *domain.User) domain.User {
	return domain.User{
		ID:         uuid.New(),
		SteamID:    steamID,
		Nickname:   steamUser.Nickname,
		AvatarURL:  steamUser.AvatarURL,
		ProfileURL: steamUser.ProfileURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (s *PlayerProfileService) searchByNickname(ctx context.Context, query string) ([]domain.User, error) {
	if len(query) < 3 {
		return []domain.User{}, nil
	}

	localResults, apiResults, err := s.fetchSearchResults(ctx, query)
	if err != nil {
		return []domain.User{}, err
	}

	combinedUsers := s.combineSearchResults(localResults, apiResults)
	return s.sortAndReturnResults(combinedUsers), nil
}

func (s *PlayerProfileService) fetchSearchResults(ctx context.Context, query string) ([]domain.User, []domain.SteamProfileSearch, error) {
	var localResults []domain.User
	var apiResults []domain.SteamProfileSearch
	var wg sync.WaitGroup
	errs := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		res, err := s.playerProfileRepository.SearchByNickname(ctx, query)
		if err != nil {
			errs <- err
			return
		}
		localResults = res
	}()

	go func() {
		defer wg.Done()
		searchResult, err := s.deadlockAPIClient.FetchSteamProfileSearch(query)
		if err != nil {
			s.logger.Warn("Failed to fetch from Deadlock API search", zap.Error(err))
			return
		}

		const maxAPIResults = 25
		if len(searchResult) > maxAPIResults {
			apiResults = searchResult[:maxAPIResults]
		} else {
			apiResults = searchResult
		}
	}()

	wg.Wait()
	close(errs)

	for err := range errs {
		s.logger.Error("Error during player search", zap.Error(err))
	}

	return localResults, apiResults, nil
}

func (s *PlayerProfileService) combineSearchResults(localResults []domain.User, apiResults []domain.SteamProfileSearch) map[string]domain.User {
	combinedUsers := make(map[string]domain.User)

	for _, u := range localResults {
		combinedUsers[u.SteamID] = u
	}

	for _, apiPlayer := range apiResults {
		if apiPlayer.Personaname == "" {
			continue
		}

		steamID := strconv.Itoa(apiPlayer.AccountID)
		if _, exists := combinedUsers[steamID]; !exists {
			user := s.createUserFromAPISearch(apiPlayer)
			err := s.userRepository.FindOrCreate(&user)
			if err != nil {
				s.logger.Error("Failed to create user from API search", zap.Error(err))
				continue
			}
			combinedUsers[steamID] = user
		}
	}

	return combinedUsers
}

func (s *PlayerProfileService) createUserFromAPISearch(apiPlayer domain.SteamProfileSearch) domain.User {
	steamID := strconv.Itoa(apiPlayer.AccountID)
	return domain.User{
		ID:         uuid.New(),
		SteamID:    steamID,
		Nickname:   apiPlayer.Personaname,
		AvatarURL:  apiPlayer.Avatar,
		ProfileURL: apiPlayer.Profileurl,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (s *PlayerProfileService) sortAndReturnResults(combinedUsers map[string]domain.User) []domain.User {
	finalResults := make([]domain.User, 0, len(combinedUsers))
	for _, user := range combinedUsers {
		finalResults = append(finalResults, user)
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Nickname < finalResults[j].Nickname
	})

	return finalResults
}

func (s *PlayerProfileService) calculateAndFillStats(profile *domain.PlayerProfile, card *deadlockapi.DeadlockCard, matches []domain.Match, mmrHistory []domain.DeadlockMMR) {
	if card == nil {
		return
	}

	stats := s.calculateMatchStats(matches)
	s.fillBasicStats(profile, card, matches, stats)
	s.fillRankInfo(profile, card, mmrHistory)
	s.fillPerformanceData(profile, matches)

	profile.LastUpdatedAt = time.Now()
}

func (s *PlayerProfileService) calculateMatchStats(matches []domain.Match) struct {
	wins, totalKills, totalDeaths, totalAssists, totalSouls, totalDuration int
} {
	var stats struct {
		wins, totalKills, totalDeaths, totalAssists, totalSouls, totalDuration int
	}

	for _, match := range matches {
		if match.MatchResult == 1 {
			stats.wins++
		}
		stats.totalKills += match.PlayerKills
		stats.totalDeaths += match.PlayerDeaths
		stats.totalAssists += match.PlayerAssists
		stats.totalSouls += match.NetWorth
		stats.totalDuration += match.MatchDurationS
	}

	return stats
}

func (s *PlayerProfileService) fillBasicStats(profile *domain.PlayerProfile, card *deadlockapi.DeadlockCard, matches []domain.Match, stats struct {
	wins, totalKills, totalDeaths, totalAssists, totalSouls, totalDuration int
}) {
	profile.TotalMatches = len(matches)
	if profile.TotalMatches > 0 {
		profile.WinRate = (float64(stats.wins) / float64(profile.TotalMatches)) * 100
	}

	if stats.totalDeaths > 0 {
		profile.KDRatio = float64(stats.totalKills+stats.totalAssists) / float64(stats.totalDeaths)
	} else if stats.totalKills > 0 || stats.totalAssists > 0 {
		profile.KDRatio = float64(stats.totalKills + stats.totalAssists)
	}

	if stats.totalDuration > 0 {
		profile.AvgSoulsPerMin = (float64(stats.totalSouls) * 60) / float64(stats.totalDuration)
	} else {
		profile.AvgSoulsPerMin = 0
	}
}

func (s *PlayerProfileService) fillRankInfo(profile *domain.PlayerProfile, card *deadlockapi.DeadlockCard, mmrHistory []domain.DeadlockMMR) {
	finalRank := 0
	if len(mmrHistory) > 0 {
		finalRank = domain.GetRankFromScore(mmrHistory[0].PlayerScore)
	} else if card.RankedRank != nil {
		finalRank = *card.RankedRank
	}

	profile.PlayerRank = finalRank
	if finalRank > 0 {
		tier := finalRank / 10
		subTier := finalRank % 10
		profile.RankName, _, _ = s.getRankNameAndSubRank(tier)
		profile.RankImage = s.getRankImageURL(tier, subTier)
		profile.SubRank = subTier
	} else {
		profile.RankName = "Unranked"
		profile.SubRank = 0
	}
}

func (s *PlayerProfileService) fillPerformanceData(profile *domain.PlayerProfile, matches []domain.Match) {
	if len(matches) > 0 {
		profile.PerformanceDynamics = domain.CalculatePerformanceDynamics(matches)
		profile.LastMatchTime = matches[0].MatchTime
	}
}

func (s *PlayerProfileService) getRankNameAndSubRank(tier int) (string, int, string) {
	if r, ok := s.staticDataService.Ranks[tier]; ok {
		var img string
		if r.Images.Large != nil {
			img = *r.Images.Large
		}
		return r.Name, 0, img
	}
	return "Unranked", 0, ""
}

func (s *PlayerProfileService) getRankImageURL(tier, subTier int) string {
	rank, found := s.staticDataService.Ranks[tier]
	if !found {
		return ""
	}

	var imageURL *string
	switch subTier {
	case 1:
		imageURL = rank.Images.LargeSubrank1
	case 2:
		imageURL = rank.Images.LargeSubrank2
	case 3:
		imageURL = rank.Images.LargeSubrank3
	case 4:
		imageURL = rank.Images.LargeSubrank4
	case 5:
		imageURL = rank.Images.LargeSubrank5
	case 6:
		imageURL = rank.Images.LargeSubrank6
	default:
		imageURL = rank.Images.Large
	}

	if imageURL != nil {
		return *imageURL
	}
	return ""
}

func (s *PlayerProfileService) enrichFeaturedHeroes(card *deadlockapi.DeadlockCard) []domain.FeaturedHero {
	if card == nil || len(card.Slots) == 0 {
		return []domain.FeaturedHero{}
	}

	var featuredHeroes []domain.FeaturedHero
	for _, slot := range card.Slots {
		if slot.Hero.ID != 0 {
			hero := s.createFeaturedHero(slot)
			featuredHeroes = append(featuredHeroes, hero)
		}
	}

	return featuredHeroes
}

func (s *PlayerProfileService) createFeaturedHero(slot deadlockapi.DeadlockCardSlot) domain.FeaturedHero {
	hero := domain.FeaturedHero{
		HeroID: slot.Hero.ID,
	}

	s.enrichHeroWithStaticData(&hero)
	s.enrichHeroWithSlotData(&hero, slot)

	return hero
}

func (s *PlayerProfileService) enrichHeroWithStaticData(hero *domain.FeaturedHero) {
	if heroData, ok := s.staticDataService.HeroesByHeroID[hero.HeroID]; ok {
		hero.HeroName = heroData.Name
		if heroData.Images.IconHeroCard != nil {
			hero.HeroImage = *heroData.Images.IconHeroCard
		}
	}

	if hero.HeroName == "" {
		hero.HeroName = fmt.Sprintf("Hero %d", hero.HeroID)
	}
}

func (s *PlayerProfileService) enrichHeroWithSlotData(hero *domain.FeaturedHero, slot deadlockapi.DeadlockCardSlot) {
	if slot.Hero.Kills > 0 {
		hero.Kills = slot.Hero.Kills
	}
	if slot.Hero.Wins > 0 {
		hero.Wins = slot.Hero.Wins
	}
	if slot.Stat != nil {
		hero.StatID = slot.Stat.StatID
		hero.StatScore = slot.Stat.StatScore
	}
}
