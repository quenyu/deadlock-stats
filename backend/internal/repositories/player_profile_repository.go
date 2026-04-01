package repositories

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"gorm.io/gorm"
)

type PlayerProfilePostgresRepository struct {
	db *gorm.DB
}

func NewPlayerProfilePostgresRepository(db *gorm.DB) *PlayerProfilePostgresRepository {
	return &PlayerProfilePostgresRepository{db: db}
}

func (r *PlayerProfilePostgresRepository) FindBySteamID(ctx context.Context, steamID string) (*domain.PlayerProfile, error) {
	user, err := r.findUserBySteamID(ctx, steamID)
	if err != nil {
		return nil, err
	}

	profile, err := r.buildPlayerProfile(ctx, steamID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return profile, nil
	}

	trendMatches, err := r.fetchTrendMatches(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	profile.RecentMatches = r.getRecentMatches(trendMatches)
	profile.PerformanceDynamics = calculatePerformanceDynamics(trendMatches)

	heroStats, err := r.fetchHeroStats(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	profile.HeroStats = heroStats

	return profile, nil
}

func (r *PlayerProfilePostgresRepository) findUserBySteamID(ctx context.Context, steamID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("steam_id = ?", steamID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PlayerProfilePostgresRepository) buildPlayerProfile(ctx context.Context, steamID string) (*domain.PlayerProfile, error) {
	var profile domain.PlayerProfile

	result := r.db.WithContext(ctx).Table("users").
		Select(r.getPlayerProfileSelectQuery()).
		Joins("LEFT JOIN player_stats ps ON users.id = ps.user_id").
		Where("users.steam_id = ?", steamID).
		First(&profile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &profile, nil
}

func (r *PlayerProfilePostgresRepository) getPlayerProfileSelectQuery() string {
	return `
		users.steam_id,
		users.nickname,
		users.avatar_url,
		users.profile_url,
		users.created_at,
		ps.kd_ratio,
		ps.win_rate,
		ps.avg_matches_per_day,
		ps.favorite_hero,
		ps.last_updated_at,
		ps.player_rank,
		ps.total_matches,
		ps.total_kills,
		ps.total_deaths,
		ps.total_assists,
		ps.max_kills_in_match,
		ps.avg_damage_per_match,
		ps.avg_objectives_per_match,
		ps.avg_souls_per_min
	`
}

func (r *PlayerProfilePostgresRepository) fetchTrendMatches(ctx context.Context, userID uuid.UUID) ([]domain.Match, error) {
	var trendMatches []domain.Match

	err := r.db.WithContext(ctx).Table("player_match_stats as pms").
		Select(r.getTrendMatchesSelectQuery()).
		Joins("JOIN matches as m ON pms.match_id = m.id").
		Where("pms.user_id = ?", userID).
		Order("m.match_time DESC").
		Limit(25).
		Find(&trendMatches).Error

	return trendMatches, err
}

func (r *PlayerProfilePostgresRepository) getTrendMatchesSelectQuery() string {
	return `
		pms.match_id as id,
		pms.hero_name,
		pms.result,
		pms.kills,
		pms.deaths,
		pms.assists,
		m.duration_minutes,
		pms.player_rank_change,
		pms.player_rank_after_match,
		m.match_time
	`
}

func (r *PlayerProfilePostgresRepository) getRecentMatches(trendMatches []domain.Match) []domain.Match {
	if len(trendMatches) > 5 {
		return trendMatches[:5]
	}
	return trendMatches
}

func (r *PlayerProfilePostgresRepository) fetchHeroStats(ctx context.Context, userID uuid.UUID) ([]domain.HeroStat, error) {
	var heroStats []domain.HeroStat

	err := r.db.WithContext(ctx).Table("player_match_stats").
		Select(`
			hero_name,
			COUNT(*) as matches,
			SUM(CASE WHEN result = 'Win' THEN 1 ELSE 0 END)::float / COUNT(*) as win_rate,
			(SUM(kills) + SUM(assists))::float / GREATEST(1, SUM(deaths))::float as kda
		`).
		Where("user_id = ?", userID).
		Group("hero_name").
		Order("matches DESC").
		Limit(5).
		Find(&heroStats).Error

	return heroStats, err
}

func (r *PlayerProfilePostgresRepository) FindRecentMatchesBySteamID(ctx context.Context, steamID string, limit int) ([]domain.Match, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("steam_id = ?", steamID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	var matches []domain.Match
	err := r.db.WithContext(ctx).Table("player_match_stats as pms").
		Select(`
			pms.match_id as id,
			pms.hero_name,
			pms.result,
			pms.kills,
			pms.deaths,
			pms.assists,
			m.duration_minutes,
			pms.player_rank_change,
			pms.player_rank_after_match,
			m.match_time
		`).
		Joins("JOIN matches as m ON pms.match_id = m.id").
		Where("pms.user_id = ?", user.ID).
		Order("m.match_time DESC").
		Limit(limit).
		Find(&matches).Error

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *PlayerProfilePostgresRepository) SearchByNickname(ctx context.Context, query string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).
		Where("nickname ILIKE ?", fmt.Sprintf("%%%s%%", query)).
		Limit(10).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PlayerProfilePostgresRepository) SearchByNicknamePartial(ctx context.Context, query string, limit int) ([]domain.User, error) {
	var users []domain.User

	err := r.db.WithContext(ctx).
		Where("nickname ILIKE ?", fmt.Sprintf("%s%%", query)).
		Or("nickname ILIKE ?", fmt.Sprintf("%%%s%%", query)).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PlayerProfilePostgresRepository) SearchBySteamIDPartial(ctx context.Context, query string, limit int) ([]domain.User, error) {
	var users []domain.User

	err := r.db.WithContext(ctx).
		Where("steam_id LIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PlayerProfilePostgresRepository) GetPopularPlayers(ctx context.Context, limit int) ([]domain.User, error) {
	var users []domain.User

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PlayerProfilePostgresRepository) GetRecentlyActivePlayers(ctx context.Context, limit int) ([]domain.User, error) {
	var users []domain.User

	err := r.db.WithContext(ctx).
		Order("updated_at DESC").
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PlayerProfilePostgresRepository) UpdateProfile(ctx context.Context, profile *domain.PlayerProfile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.updatePlayerStats(tx, profile); err != nil {
			return err
		}

		return r.updateMatches(tx, profile)
	})
}

func (r *PlayerProfilePostgresRepository) updatePlayerStats(tx *gorm.DB, profile *domain.PlayerProfile) error {
	statsQuery := `
		INSERT INTO player_stats (user_id, kd_ratio, win_rate, avg_matches_per_day, favorite_hero, player_rank, total_matches, total_kills, total_deaths, total_assists, max_kills_in_match, avg_damage_per_match, avg_objectives_per_match, avg_souls_per_min, last_updated_at)
		SELECT u.id, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW()
		FROM users u WHERE u.steam_id = $1
		ON CONFLICT (user_id) DO UPDATE SET
			kd_ratio = EXCLUDED.kd_ratio,
			win_rate = EXCLUDED.win_rate,
			avg_matches_per_day = EXCLUDED.avg_matches_per_day,
			favorite_hero = EXCLUDED.favorite_hero,
			player_rank = EXCLUDED.player_rank,
			total_matches = EXCLUDED.total_matches,
			total_kills = EXCLUDED.total_kills,
			total_deaths = EXCLUDED.total_deaths,
			total_assists = EXCLUDED.total_assists,
			max_kills_in_match = EXCLUDED.max_kills_in_match,
			avg_damage_per_match = EXCLUDED.avg_damage_per_match,
			avg_objectives_per_match = EXCLUDED.avg_objectives_per_match,
			avg_souls_per_min = EXCLUDED.avg_souls_per_min,
			last_updated_at = NOW();
	`

	return tx.Exec(statsQuery,
		profile.SteamID, profile.KDRatio, profile.WinRate, profile.AvgMatchesPerDay, profile.FavoriteHero,
		profile.PlayerRank, profile.TotalMatches, profile.TotalKills, profile.TotalDeaths, profile.TotalAssists,
		profile.MaxKillsInMatch, profile.AvgDamagePerMatch, profile.AvgObjectivesPerMatch, profile.AvgSoulsPerMin,
	).Error
}

func (r *PlayerProfilePostgresRepository) updateMatches(tx *gorm.DB, profile *domain.PlayerProfile) error {
	for _, match := range profile.RecentMatches {
		if err := r.insertMatch(tx, match); err != nil {
			return err
		}

		if err := r.insertPlayerMatchStats(tx, profile.SteamID, match); err != nil {
			return err
		}
	}
	return nil
}

func (r *PlayerProfilePostgresRepository) insertMatch(tx *gorm.DB, match domain.Match) error {
	matchQuery := `
		INSERT INTO matches (id, map_name, duration_minutes, match_time) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (id) DO NOTHING
	`
	return tx.Exec(matchQuery, match.ID, "Unknown Map", match.DurationMinutes, match.MatchTime).Error
}

func (r *PlayerProfilePostgresRepository) insertPlayerMatchStats(tx *gorm.DB, steamID string, match domain.Match) error {
	pmsQuery := `
		INSERT INTO player_match_stats (user_id, match_id, hero_name, kills, deaths, assists, result, player_rank_change, player_rank_after_match)
		SELECT u.id, $2, $3, $4, $5, $6, $7, $8, $9
		FROM users u WHERE u.steam_id = $1
		ON CONFLICT (user_id, match_id) DO NOTHING
	`

	return tx.Exec(pmsQuery,
		steamID, match.ID, match.HeroName, match.Kills, match.Deaths, match.Assists,
		match.Result, match.PlayerRankChange, match.PlayerRankAfterMatch,
	).Error
}

func calculatePerformanceDynamics(matches []domain.Match) domain.PerformanceDynamics {
	var dynamics domain.PerformanceDynamics
	if len(matches) < 2 {
		return dynamics
	}

	sortedMatches := sortMatchesByTime(matches)

	dynamics.Rank = calculateRankDynamics(sortedMatches)
	dynamics.WinLoss = calculateWinLossDynamics(sortedMatches)
	dynamics.KDA = calculateKDADynamics(sortedMatches)

	return dynamics
}

func sortMatchesByTime(matches []domain.Match) []domain.Match {
	sorted := make([]domain.Match, len(matches))
	copy(sorted, matches)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].MatchTime.Before(sorted[j].MatchTime)
	})

	return sorted
}

func calculateRankDynamics(matches []domain.Match) domain.Trend {
	var rank domain.Trend

	if len(matches) > 0 {
		rankStart := matches[0].PlayerRankAfterMatch - matches[0].PlayerRankChange
		rankEnd := matches[len(matches)-1].PlayerRankAfterMatch
		rankDiff := rankEnd - rankStart

		rank.Value = fmt.Sprintf("%+d Rank", rankDiff)
		rank.Trend = getTrend(float64(rankDiff))

		for _, m := range matches {
			rank.Sparkline = append(rank.Sparkline, float64(m.PlayerRankAfterMatch))
		}
	}

	return rank
}

func calculateWinLossDynamics(matches []domain.Match) domain.Trend {
	var winLoss domain.Trend

	netWins := calculateNetWins(matches)
	winLoss.Value = fmt.Sprintf("%+d WINS", netWins)
	winLoss.Trend = getTrend(float64(netWins))
	winLoss.Sparkline = calculateCumulativeWins(matches)

	return winLoss
}

func calculateNetWins(matches []domain.Match) int {
	netWins := 0
	for _, m := range matches {
		if m.Result == "Win" {
			netWins++
		} else {
			netWins--
		}
	}
	return netWins
}

func calculateCumulativeWins(matches []domain.Match) []float64 {
	var sparkline []float64
	cumulativeWins := 0

	for _, m := range matches {
		if m.Result == "Win" {
			cumulativeWins++
		}
		sparkline = append(sparkline, float64(cumulativeWins))
	}

	return sparkline
}

func calculateKDADynamics(matches []domain.Match) domain.Trend {
	var kda domain.Trend

	kdaValues := calculateKDAValues(matches)
	if len(kdaValues) > 0 {
		avgKDA := kdaValues[len(kdaValues)-1]
		kda.Value = fmt.Sprintf("%.2f KDA", avgKDA)
		kda.Trend = getTrend(kdaValues[len(kdaValues)-1] - kdaValues[0])
		kda.Sparkline = kdaValues
	}

	return kda
}

func calculateKDAValues(matches []domain.Match) []float64 {
	var kdaValues []float64

	for _, m := range matches {
		kda := calculateKDA(m.Kills, m.Deaths, m.Assists)
		kdaValues = append(kdaValues, kda)
	}

	return kdaValues
}

func calculateKDA(kills, deaths, assists int) float64 {
	if deaths == 0 {
		return float64(kills) + float64(assists)
	}
	return (float64(kills) + float64(assists)) / float64(deaths)
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
