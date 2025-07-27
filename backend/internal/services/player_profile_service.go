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
			s.logger.Info("Cache hit for player profile", zap.String("steamID", steamID))
			return &profile, nil
		}
	}

	s.logger.Info("Cache miss for player profile", zap.String("steamID", steamID))

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
		s.logger.Debug("No partial cache found", zap.String("steamID", steamID))
		return nil
	}

	var profile dto.ExtendedPlayerProfile
	if err := json.Unmarshal([]byte(val), &profile); err == nil {
		if time.Since(profile.LastUpdatedAt) < time.Hour {
			s.logger.Info("Using partial cache", zap.String("steamID", steamID))
			return &profile
		} else {
			s.logger.Debug("Partial cache expired", zap.String("steamID", steamID))
		}
	} else {
		s.logger.Warn("Failed to unmarshal partial cache", zap.String("steamID", steamID), zap.Error(err))
	}

	return nil
}

func (s *PlayerProfileService) updatePartialCache(ctx context.Context, steamID string, profile *dto.ExtendedPlayerProfile) {
	partialCacheKey := fmt.Sprintf("player-profile-partial:%s", steamID)

	partialProfile := &dto.ExtendedPlayerProfile{
		Nickname:            profile.Nickname,
		AvatarURL:           profile.AvatarURL,
		PlayerRank:          profile.PlayerRank,
		RankName:            profile.RankName,
		RankImage:           profile.RankImage,
		TotalMatches:        profile.TotalMatches,
		WinRate:             profile.WinRate,
		KDRatio:             profile.KDRatio,
		PerformanceDynamics: profile.PerformanceDynamics,
		LastUpdatedAt:       time.Now(),
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

	matches, heroStats, mmrHistory, profile, heroMMRHistory, err := s.fetchAllData(ctx, steamID)
	if err != nil {
		return nil, err
	}

	extendedProfile := s.buildExtendedProfile(matches, heroStats, mmrHistory, profile, heroMMRHistory)

	s.cacheProfile(ctx, cacheKey, extendedProfile)

	return extendedProfile, nil
}

func (s *PlayerProfileService) fetchAllData(ctx context.Context, steamID string) (
	[]deadlockapi.DeadlockMatch,
	[]domain.HeroStat,
	[]domain.DeadlockMMR,
	*domain.PlayerProfile,
	[]domain.HeroMMRHistory,
	error,
) {
	apiCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var matches []deadlockapi.DeadlockMatch
	var heroStats []domain.HeroStat
	var mmrHistory []domain.DeadlockMMR
	var profile *domain.PlayerProfile
	var heroMMRHistory []domain.HeroMMRHistory

	var wg sync.WaitGroup
	errs := make(chan error, 4)

	wg.Add(4)
	go func() {
		defer wg.Done()
		profile, _ = s.playerProfileRepository.FindBySteamID(ctx, steamID)
	}()
	go func() {
		defer wg.Done()
		var err error
		s.logger.Info("Fetching match history from API", zap.String("steamID", steamID))
		matches, err = s.deadlockAPIClient.FetchMatchHistory(steamID)
		if err != nil {
			s.logger.Error("Failed to fetch match history", zap.String("steamID", steamID), zap.Error(err))
			errs <- err
			matches = []deadlockapi.DeadlockMatch{}
		} else {
			s.logger.Info("Successfully fetched match history", zap.String("steamID", steamID), zap.Int("matchesCount", len(matches)))
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		heroStats, err = s.deadlockAPIClient.FetchHeroStats(steamID)
		if err != nil {
			errs <- err
			heroStats = []domain.HeroStat{}
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		mmrHistory, err = s.deadlockAPIClient.FetchMMRHistory(steamID)
		if err != nil {
			errs <- err
			mmrHistory = []domain.DeadlockMMR{}
		}
	}()

	wg.Wait()
	close(errs)

	for err := range errs {
		s.logger.Warn("Failed to fetch some data from Deadlock API", zap.Error(err))
	}

	if len(matches) == 0 && len(heroStats) == 0 {
		return nil, nil, nil, nil, nil, fmt.Errorf("critical data could not be fetched (matches or hero stats)")
	}

	if profile == nil {
		steamProfiles, err := s.deadlockAPIClient.FetchSteamProfileSearch(steamID)
		if err != nil || len(steamProfiles) == 0 {
			profile = &domain.PlayerProfile{
				SteamID:    steamID,
				Nickname:   "Unknown Player",
				AvatarURL:  "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb.jpg",
				ProfileURL: fmt.Sprintf("https://steamcommunity.com/profiles/%s", steamID),
			}
		} else {
			apiProfile := steamProfiles[0]
			profile = &domain.PlayerProfile{
				SteamID:    steamID,
				Nickname:   apiProfile.Personaname,
				AvatarURL:  apiProfile.Avatar,
				ProfileURL: apiProfile.Profileurl,
			}
		}
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

	return matches, heroStats, mmrHistory, profile, heroMMRHistory, nil
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
	matches []deadlockapi.DeadlockMatch,
	heroStats []domain.HeroStat,
	mmrHistory []domain.DeadlockMMR,
	profile *domain.PlayerProfile,
	heroMMRHistory []domain.HeroMMRHistory,
) *dto.ExtendedPlayerProfile {
	s.logger.Info("Building extended profile",
		zap.Int("apiMatchesCount", len(matches)),
		zap.Int("heroStatsCount", len(heroStats)),
		zap.Int("mmrHistoryCount", len(mmrHistory)))

	domainMatches := s.buildDomainMatches(matches, mmrHistory)
	s.logger.Info("Domain matches built", zap.Int("domainMatchesCount", len(domainMatches)))

	s.calculateAndFillStats(profile, domainMatches, mmrHistory)

	featuredHeroes := []domain.FeaturedHero{}
	peakRank, peakRankName, peakRankImage := domain.FindPeakRank(mmrHistory, s.getRankNameAndSubRank, s.getRankImageURL)
	personalRecords := domain.CalculatePersonalRecords(domainMatches)
	avgStats := domain.CalculateAverageStats(domainMatches, len(matches))

	validMateStats := []domain.MateStat{}

	dtoRecords := s.buildPersonalRecordsDTO(personalRecords)
	dtoMMRHistory := s.buildMMRHistoryDTO(mmrHistory)
	dtoHeroMMR := s.buildHeroMMRHistoryDTO(heroMMRHistory)

	s.logger.Info("Building DTO with performance dynamics")

	return &dto.ExtendedPlayerProfile{
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
		s.logger.Debug("Cached profile", zap.String("cacheKey", cacheKey))
	} else {
		s.logger.Warn("Failed to marshal profile for caching", zap.Error(err))
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

func (s *PlayerProfileService) buildMateStats() []domain.MateStat {
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
	s.logger.Info("Building domain matches", zap.Int("apiMatchesCount", len(matches)), zap.Int("mmrHistoryCount", len(mmrHistory)))

	rankMap := s.createRankMap(mmrHistory)
	s.sortMatchesByTime(matches)

	domainMatches := s.convertToDomainMatches(matches, rankMap)
	s.logger.Info("Converted to domain matches", zap.Int("domainMatchesCount", len(domainMatches)))

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
			PlayerTeam:     match.PlayerTeam,
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
	matches, err := s.playerProfileRepository.FindRecentMatchesBySteamID(ctx, steamID, 5)
	if err != nil {
		if err.Error() == "user not found" {
			return []domain.Match{}, nil
		}
		return nil, err
	}
	return matches, nil
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

func (s *PlayerProfileService) SearchPlayersWithAutocomplete(ctx context.Context, query string, limit int) ([]domain.User, error) {
	if len(query) < 2 {
		return []domain.User{}, nil
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	if _, err := strconv.ParseInt(query, 10, 64); err == nil {
		user, err := s.authService.SearchPlayersBySteamID(query)
		if err != nil {
			s.logger.Error("Error finding user by SteamID64", zap.Error(err), zap.String("steamID", query))
			return []domain.User{}, nil
		}
		if user != nil {
			return []domain.User{*user}, nil
		}
	}

	resolvedSteamID, err := s.authService.ResolveVanityURL(query)
	if err == nil && resolvedSteamID != "" {
		user, err := s.authService.SearchPlayersBySteamID(resolvedSteamID)
		if err != nil {
			s.logger.Error("Error finding user by resolved vanity URL", zap.Error(err), zap.String("vanityURL", query), zap.String("resolvedID", resolvedSteamID))
			return []domain.User{}, nil
		}
		if user != nil {
			return []domain.User{*user}, nil
		}
	}

	nicknameResults, err := s.playerProfileRepository.SearchByNicknamePartial(ctx, query, limit)
	if err != nil {
		s.logger.Error("Error searching by partial nickname", zap.Error(err))
		return []domain.User{}, err
	}

	steamIDResults, err := s.playerProfileRepository.SearchBySteamIDPartial(ctx, query, limit)
	if err != nil {
		s.logger.Error("Error searching by partial Steam ID", zap.Error(err))
		return []domain.User{}, err
	}

	allResults := make([]domain.User, 0, len(nicknameResults)+len(steamIDResults))
	seenSteamIDs := make(map[string]bool)

	for _, user := range nicknameResults {
		if s.IsValidUser(user) {
			allResults = append(allResults, user)
			seenSteamIDs[user.SteamID] = true
		}
	}

	for _, user := range steamIDResults {
		if s.IsValidUser(user) && !seenSteamIDs[user.SteamID] {
			allResults = append(allResults, user)
			seenSteamIDs[user.SteamID] = true
		}
	}

	validResults := allResults

	sort.Slice(validResults, func(i, j int) bool {
		queryLower := strings.ToLower(query)
		nickI := strings.ToLower(validResults[i].Nickname)
		nickJ := strings.ToLower(validResults[j].Nickname)

		startsWithI := strings.HasPrefix(nickI, queryLower)
		startsWithJ := strings.HasPrefix(nickJ, queryLower)

		if startsWithI && !startsWithJ {
			return true
		}
		if !startsWithI && startsWithJ {
			return false
		}

		if len(nickI) != len(nickJ) {
			return len(nickI) < len(nickJ)
		}

		return nickI < nickJ
	})

	if len(validResults) > limit {
		validResults = validResults[:limit]
	}

	return validResults, nil
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
		if u.SteamID != "" && u.Nickname != "" {
			combinedUsers[u.SteamID] = u
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
			user := s.createUserFromAPISearch(apiPlayer, steamID64)

			err := s.userRepository.FindOrCreate(&user)
			if err != nil {
				s.logger.Error("Failed to create user from API search",
					zap.Error(err),
					zap.String("steamID", steamID64),
					zap.String("nickname", user.Nickname))
				continue
			}

			combinedUsers[steamID64] = user
			s.logger.Debug("Created user from API search",
				zap.String("steamID", steamID64),
				zap.String("nickname", user.Nickname))
		}
	}

	return combinedUsers
}

func (s *PlayerProfileService) isValidAPIPlayer(apiPlayer domain.SteamProfileSearch) bool {
	if apiPlayer.AccountID <= 0 {
		return false
	}

	if apiPlayer.Personaname == "" {
		return false
	}

	if apiPlayer.Avatar == "" || apiPlayer.Avatar == "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb.jpg" {
		return true
	}

	return true
}

func (s *PlayerProfileService) convertAccountIDToSteamID64(accountID int) string {
	steamID64 := int64(accountID) + 76561197960265728
	return strconv.FormatInt(steamID64, 10)
}

func (s *PlayerProfileService) createUserFromAPISearch(apiPlayer domain.SteamProfileSearch, steamID64 string) domain.User {
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

func (s *PlayerProfileService) sortAndReturnResults(combinedUsers map[string]domain.User) []domain.User {
	finalResults := make([]domain.User, 0, len(combinedUsers))

	for _, user := range combinedUsers {
		if s.IsValidUser(user) {
			finalResults = append(finalResults, user)
		} else {
			s.logger.Debug("Filtering out invalid user from search results",
				zap.String("steamID", user.SteamID),
				zap.String("nickname", user.Nickname),
				zap.String("avatarURL", user.AvatarURL))
		}
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Nickname < finalResults[j].Nickname
	})

	s.logger.Info("Search results filtered and sorted",
		zap.Int("totalFound", len(combinedUsers)),
		zap.Int("validResults", len(finalResults)))

	return finalResults
}

func (s *PlayerProfileService) IsValidUser(user domain.User) bool {
	if user.SteamID == "" {
		return false
	}

	if user.Nickname == "" {
		return false
	}

	trimmedNickname := strings.TrimSpace(user.Nickname)
	if trimmedNickname == "" || len(trimmedNickname) < 2 {
		return false
	}

	if user.AvatarURL == "" {
		return false
	}

	return true
}

func (s *PlayerProfileService) calculateAndFillStats(profile *domain.PlayerProfile, matches []domain.Match, mmrHistory []domain.DeadlockMMR) {
	stats := s.calculateMatchStats(matches)
	s.fillBasicStats(profile, matches, stats)
	s.fillRankInfo(profile, mmrHistory)
	s.fillPerformanceData(profile, matches)

	profile.LastUpdatedAt = time.Now()
}

func (s *PlayerProfileService) calculateMatchStats(matches []domain.Match) struct {
	wins, losses, totalKills, totalDeaths, totalAssists, totalSouls, totalDuration int
} {
	var stats struct {
		wins, losses, totalKills, totalDeaths, totalAssists, totalSouls, totalDuration int
	}

	for _, match := range matches {
		if match.PlayerTeam == 1 {
			if match.MatchResult == 1 {
				stats.wins++
			} else {
				stats.losses++
			}
		} else {
			if match.MatchResult == 0 {
				stats.wins++
			} else {
				stats.losses++
			}
		}
		stats.totalKills += match.PlayerKills
		stats.totalDeaths += match.PlayerDeaths
		stats.totalAssists += match.PlayerAssists
		stats.totalSouls += match.NetWorth
		stats.totalDuration += match.MatchDurationS
	}

	return stats
}

func (s *PlayerProfileService) fillBasicStats(profile *domain.PlayerProfile, matches []domain.Match, stats struct {
	wins, losses, totalKills, totalDeaths, totalAssists, totalSouls, totalDuration int
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

func (s *PlayerProfileService) fillRankInfo(profile *domain.PlayerProfile, mmrHistory []domain.DeadlockMMR) {
	finalRank := 0
	if len(mmrHistory) > 0 {
		finalRank = domain.GetRankFromScore(mmrHistory[0].PlayerScore)
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
	s.logger.Info("Filling performance data", zap.Int("matchesCount", len(matches)))

	if len(matches) > 0 {
		s.logger.Info("Calculating performance dynamics from matches")
		profile.PerformanceDynamics = domain.CalculatePerformanceDynamics(matches)
		profile.LastMatchTime = matches[0].MatchTime
		s.logger.Info("Performance dynamics calculated",
			zap.String("winLossTrend", profile.PerformanceDynamics.WinLoss.Trend),
			zap.String("winLossValue", profile.PerformanceDynamics.WinLoss.Value))
	} else {
		s.logger.Warn("No matches found, initializing empty performance dynamics")
		profile.PerformanceDynamics = domain.PerformanceDynamics{
			WinLoss: domain.Trend{Trend: "stable", Value: "0/0", Sparkline: []float64{}},
			KDA:     domain.Trend{Trend: "stable", Value: "0.00 KDA", Sparkline: []float64{}},
			Rank:    domain.Trend{Trend: "stable", Value: "0 Rank", Sparkline: []float64{}},
		}
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

func (s *PlayerProfileService) enrichFeaturedHeroes() []domain.FeaturedHero {
	// Убираем логику, связанную с карточкой
	return []domain.FeaturedHero{}
}

func (s *PlayerProfileService) SearchPlayersWithFilters(ctx context.Context, query string, filters dto.SearchFilters, limit int) ([]domain.User, error) {
	if len(query) < 2 {
		return []domain.User{}, nil
	}

	users, err := s.playerProfileRepository.SearchByNicknamePartial(ctx, query, limit)
	if err != nil {
		return []domain.User{}, err
	}

	validResults := make([]domain.User, 0, len(users))
	for _, user := range users {
		if s.IsValidUser(user) {
			validResults = append(validResults, user)
		}
	}

	s.sortUsersByFilters(validResults, filters)

	return validResults, nil
}

func (s *PlayerProfileService) sortUsersByFilters(users []domain.User, filters dto.SearchFilters) {
	sort.Slice(users, func(i, j int) bool {
		switch filters.SortBy {
		case "created_at":
			if filters.SortOrder == "desc" {
				return users[i].CreatedAt.After(users[j].CreatedAt)
			}
			return users[i].CreatedAt.Before(users[j].CreatedAt)
		case "updated_at":
			if filters.SortOrder == "desc" {
				return users[i].UpdatedAt.After(users[j].UpdatedAt)
			}
			return users[i].UpdatedAt.Before(users[j].UpdatedAt)
		default: // "nickname"
			if filters.SortOrder == "desc" {
				return users[i].Nickname > users[j].Nickname
			}
			return users[i].Nickname < users[j].Nickname
		}
	})
}

func (s *PlayerProfileService) GetPopularPlayers(ctx context.Context, limit int) ([]domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	users, err := s.playerProfileRepository.GetPopularPlayers(ctx, limit)
	if err != nil {
		s.logger.Error("Error getting popular players", zap.Error(err))
		return []domain.User{}, err
	}

	validResults := make([]domain.User, 0, len(users))
	for _, user := range users {
		if s.IsValidUser(user) {
			validResults = append(validResults, user)
		}
	}

	s.logger.Info("Retrieved popular players",
		zap.Int("requested", limit),
		zap.Int("found", len(validResults)))

	return validResults, nil
}

func (s *PlayerProfileService) GetRecentlyActivePlayers(ctx context.Context, limit int) ([]domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	users, err := s.playerProfileRepository.GetRecentlyActivePlayers(ctx, limit)
	if err != nil {
		s.logger.Error("Error getting recently active players", zap.Error(err))
		return []domain.User{}, err
	}

	validResults := make([]domain.User, 0, len(users))
	for _, user := range users {
		if s.IsValidUser(user) {
			validResults = append(validResults, user)
		}
	}

	s.logger.Info("Retrieved recently active players",
		zap.Int("requested", limit),
		zap.Int("found", len(validResults)))

	return validResults, nil
}
