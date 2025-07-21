package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/quenyu/deadlock-stats/internal/dto"
	"github.com/quenyu/deadlock-stats/internal/repositories"
	"go.uber.org/zap"
)

type PlayerProfileService struct {
	playerProfileRepository *repositories.PlayerProfilePostgresRepository
	deadlockAPIClient       *deadlockapi.Client
	staticDataService       *StaticDataService
	redisClient             *redis.Client
	logger                  *zap.Logger
}

func NewPlayerProfileService(
	playerProfileRepository *repositories.PlayerProfilePostgresRepository,
	deadlockAPIClient *deadlockapi.Client,
	staticDataService *StaticDataService,
	redisClient *redis.Client,
	logger *zap.Logger,
) *PlayerProfileService {
	return &PlayerProfileService{
		playerProfileRepository: playerProfileRepository,
		deadlockAPIClient:       deadlockAPIClient,
		staticDataService:       staticDataService,
		redisClient:             redisClient,
		logger:                  logger,
	}
}

func (s *PlayerProfileService) GetExtendedPlayerProfile(ctx context.Context, steamID string) (*dto.ExtendedPlayerProfile, error) {
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

	var card *deadlockapi.DeadlockCard
	var matches []deadlockapi.DeadlockMatch
	var heroStats []domain.HeroStat
	var mmrHistory []deadlockapi.DeadlockMMR
	var profile *domain.PlayerProfile

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
		if err != nil {
			return nil, err
		}
	}

	if profile == nil {
		return nil, fmt.Errorf("player not found and card is missing")
	}

	domainMatches := s.buildDomainMatches(matches, mmrHistory)
	s.enrichHeroStats(heroStats)
	s.calculateAndFillStats(profile, card, domainMatches, mmrHistory)

	extendedProfile := &dto.ExtendedPlayerProfile{
		Card:                card,
		MatchHistory:        domainMatches,
		HeroStats:           heroStats,
		MMRHistory:          mmrHistory,
		TotalMatches:        profile.TotalMatches,
		WinRate:             profile.WinRate,
		KDRatio:             profile.KDRatio,
		PerformanceDynamics: profile.PerformanceDynamics,
		Nickname:            profile.Nickname,
		AvatarURL:           profile.AvatarURL,
		RankImage:           profile.RankImage,
		AvgSoulsPerMin:      profile.AvgSoulsPerMin,
	}

	data, err := json.Marshal(extendedProfile)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, data, time.Hour).Err()
	}

	return extendedProfile, nil
}

func (s *PlayerProfileService) buildDomainMatches(matches []deadlockapi.DeadlockMatch, mmrHistory []deadlockapi.DeadlockMMR) []domain.Match {
	rankMap := make(map[int64]deadlockapi.DeadlockMMR)
	for _, r := range mmrHistory {
		rankMap[r.MatchID] = r
	}

	domainMatches := make([]domain.Match, len(matches))
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].StartTime > matches[j].StartTime
	})

	var lastRank int
	if len(mmrHistory) > 0 {
		sort.SliceStable(mmrHistory, func(i, j int) bool {
			return mmrHistory[i].StartTime > mmrHistory[j].StartTime
		})
		lastRank = mmrHistory[0].Rank
	}

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
			Result:         mapMatchResult(match.MatchResult),
		}

		if hero, ok := s.staticDataService.HeroesByHeroID[match.HeroID]; ok {
			domainMatches[i].HeroName = hero.Name
			if hero.Images.IconHeroCard != nil {
				domainMatches[i].HeroAvatar = *hero.Images.IconHeroCard
			}
		}

		mmr, ok := rankMap[match.MatchID]
		if ok {
			calculatedRank := getRankFromScore(mmr.PlayerScore)
			tier := calculatedRank / 10
			subTier := calculatedRank % 10

			domainMatches[i].PlayerRankAfterMatch = calculatedRank
			if lastRank != 0 {
				domainMatches[i].PlayerRankChange = calculatedRank - lastRank
			}
			domainMatches[i].RankName, _, domainMatches[i].RankImage = s.getRankNameAndSubRank(tier)
			domainMatches[i].SubRank = subTier
			lastRank = calculatedRank
		}
	}
	return domainMatches
}

func (s *PlayerProfileService) enrichHeroStats(heroStats []domain.HeroStat) {
	for i := range heroStats {
		if hero, ok := s.staticDataService.HeroesByHeroID[heroStats[i].HeroID]; ok {
			heroStats[i].HeroName = hero.Name
			if hero.Images.IconHeroCard != nil {
				heroStats[i].HeroAvatar = *hero.Images.IconHeroCard
			}
		}
	}
}

func (s *PlayerProfileService) GetRecentMatches(ctx context.Context, steamID string) ([]domain.Match, error) {
	return s.playerProfileRepository.FindRecentMatchesBySteamID(ctx, steamID, 5)
}

func (s *PlayerProfileService) SearchPlayers(ctx context.Context, query string) ([]domain.User, error) {
	return s.playerProfileRepository.SearchByNickname(ctx, query)
}

func (s *PlayerProfileService) calculateAndFillStats(profile *domain.PlayerProfile, card *deadlockapi.DeadlockCard, matches []domain.Match, mmrHistory []deadlockapi.DeadlockMMR) {
	var wins, totalKills, totalDeaths, totalAssists, totalSouls, totalDuration int
	for _, match := range matches {
		if match.MatchResult == 1 {
			wins++
		}
		totalKills += match.PlayerKills
		totalDeaths += match.PlayerDeaths
		totalAssists += match.PlayerAssists
		totalSouls += match.NetWorth
		totalDuration += match.MatchDurationS
	}

	profile.TotalMatches = len(matches)
	if profile.TotalMatches > 0 {
		profile.WinRate = (float64(wins) / float64(profile.TotalMatches)) * 100
	}

	if totalDeaths > 0 {
		profile.KDRatio = float64(totalKills+totalAssists) / float64(totalDeaths)
	} else if totalKills > 0 || totalAssists > 0 {
		profile.KDRatio = float64(totalKills + totalAssists)
	}

	if totalDuration > 0 {
		profile.AvgSoulsPerMin = (float64(totalSouls) * 60) / float64(totalDuration)
	} else {
		profile.AvgSoulsPerMin = 0
	}

	if card != nil && card.RankedRank != nil {
		profile.PlayerRank = *card.RankedRank
	} else if len(mmrHistory) > 0 {
		profile.PlayerRank = mmrHistory[0].Rank
	}

	if len(matches) > 0 {
		profile.PerformanceDynamics = calculatePerformanceDynamics(matches)
		if len(mmrHistory) > 0 {
			mostRecentMMR := mmrHistory[0]
			calculatedRank := getRankFromScore(mostRecentMMR.PlayerScore)
			tier := calculatedRank / 10
			subTier := calculatedRank % 10

			profile.RankName, _, profile.RankImage = s.getRankNameAndSubRank(tier)
			profile.SubRank = subTier
		}
		profile.LastMatchTime = matches[0].MatchTime
	}
	profile.LastUpdatedAt = time.Now()
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

var rankScores = []float64{
	0, 11, 12, 13, 14, 15, 16, 21, 22, 23, 24, 25, 26, 31, 32, 33, 34, 35, 36, 41, 42, 43, 44, 45, 46, 51, 52, 53, 54, 55, 56, 61, 62, 63, 64, 65, 66, 71, 72, 73, 74, 75, 76, 81, 82, 83, 84, 85, 86, 91, 92, 93, 94, 95, 96, 101, 102, 103, 104, 105, 106, 111, 112, 113, 114, 115, 116,
}

func getRankFromScore(playerScore float64) int {
	closestIndex := int(math.Round(playerScore))

	if closestIndex < 0 || closestIndex >= len(rankScores) {
		if len(rankScores) > 0 {
			return int(rankScores[0])
		}
		return 0
	}

	return int(rankScores[closestIndex])
}

func calculatePerformanceDynamics(matches []domain.Match) domain.PerformanceDynamics {
	var dynamics domain.PerformanceDynamics
	if len(matches) < 2 {
		return dynamics
	}

	sortedMatches := make([]domain.Match, len(matches))
	copy(sortedMatches, matches)

	sort.SliceStable(sortedMatches, func(i, j int) bool {
		return sortedMatches[i].MatchTime.Before(sortedMatches[j].MatchTime)
	})

	if len(sortedMatches) > 0 {
		rankStart := sortedMatches[0].PlayerRankAfterMatch - sortedMatches[0].PlayerRankChange
		rankEnd := sortedMatches[len(sortedMatches)-1].PlayerRankAfterMatch
		rankDiff := rankEnd - rankStart
		dynamics.Rank.Value = fmt.Sprintf("%+d Rank", rankDiff)
		dynamics.Rank.Trend = getTrend(float64(rankDiff))
		for _, m := range sortedMatches {
			dynamics.Rank.Sparkline = append(dynamics.Rank.Sparkline, float64(m.PlayerRankAfterMatch))
		}
	}

	var winLossTrend domain.Trend
	var kdaTrend domain.Trend

	var wins, losses int
	for _, m := range sortedMatches {
		if m.Result == "Win" {
			wins++
			winLossTrend.Sparkline = append(winLossTrend.Sparkline, float64(1))
		} else {
			losses++
			winLossTrend.Sparkline = append(winLossTrend.Sparkline, float64(0))
		}

		var perMatchKDA float64
		if m.PlayerDeaths > 0 {
			perMatchKDA = float64(m.PlayerKills+m.PlayerAssists) / float64(m.PlayerDeaths)
		} else {
			perMatchKDA = float64(m.PlayerKills + m.PlayerAssists)
		}
		kdaTrend.Sparkline = append(kdaTrend.Sparkline, perMatchKDA)
	}

	if len(sortedMatches) > 0 {
		winLossTrend.Value = fmt.Sprintf("%d/%d", wins, losses)
		winLossTrend.Trend = getTrend(float64(wins - losses))

		var totalKills, totalDeaths, totalAssists float64
		for _, m := range sortedMatches {
			totalKills += float64(m.PlayerKills)
			totalDeaths += float64(m.PlayerDeaths)
			totalAssists += float64(m.PlayerAssists)
		}
		var kda float64
		if totalDeaths > 0 {
			kda = (totalKills + totalAssists) / totalDeaths
		} else {
			kda = totalKills + totalAssists
		}

		kdaTrend.Value = fmt.Sprintf("%.2f KDA", kda)
		kdaTrend.Trend = "stable"
	}

	dynamics.WinLoss = winLossTrend
	dynamics.KDA = kdaTrend

	return dynamics
}

func getTrend(value float64) string {
	if value > 0 {
		return "up"
	}
	if value < 0 {
		return "down"
	}
	return "stable"
}

func mapMatchResult(result int) string {
	if result == 1 {
		return "Win"
	}
	return "Loss"
}
