package repositories

import (
	"context"
	"errors"
	"fmt"
	"sort"

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
	var profile domain.PlayerProfile
	var user domain.User

	if err := r.db.WithContext(ctx).Where("steam_id = ?", steamID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	result := r.db.WithContext(ctx).Table("users").
		Select(`
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
		`).
		Joins("LEFT JOIN player_stats ps ON users.id = ps.user_id").
		Where("users.steam_id = ?", steamID).
		First(&profile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	var trendMatches []domain.Match
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
		Limit(25).
		Find(&trendMatches).Error

	if err != nil {
		return nil, err
	}

	if len(trendMatches) > 5 {
		profile.RecentMatches = trendMatches[:5]
	} else {
		profile.RecentMatches = trendMatches
	}

	profile.PerformanceDynamics = calculatePerformanceDynamics(trendMatches)

	var heroStats []domain.HeroStat
	err = r.db.WithContext(ctx).Table("player_match_stats").
		Select(`
			hero_name,
			COUNT(*) as matches,
			SUM(CASE WHEN result = 'Win' THEN 1 ELSE 0 END)::float / COUNT(*) as win_rate,
			(SUM(kills) + SUM(assists))::float / GREATEST(1, SUM(deaths))::float as kda
		`).
		Where("user_id = ?", user.ID).
		Group("hero_name").
		Order("matches DESC").
		Limit(5).
		Find(&heroStats).Error

	if err != nil {
		return nil, err
	}
	profile.HeroStats = heroStats

	return &profile, nil
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

func (r *PlayerProfilePostgresRepository) UpdateProfile(ctx context.Context, profile *domain.PlayerProfile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		if err := tx.Exec(statsQuery,
			profile.SteamID, profile.KDRatio, profile.WinRate, profile.AvgMatchesPerDay, profile.FavoriteHero,
			profile.PlayerRank, profile.TotalMatches, profile.TotalKills, profile.TotalDeaths, profile.TotalAssists,
			profile.MaxKillsInMatch, profile.AvgDamagePerMatch, profile.AvgObjectivesPerMatch, profile.AvgSoulsPerMin,
		).Error; err != nil {
			return err
		}

		for _, match := range profile.RecentMatches {
			matchQuery := `INSERT INTO matches (id, map_name, duration_minutes, match_time) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`
			if err := tx.Exec(matchQuery, match.ID, "Unknown Map", match.DurationMinutes, match.MatchTime).Error; err != nil {
				return err
			}

			pmsQuery := `
				INSERT INTO player_match_stats (user_id, match_id, hero_name, kills, deaths, assists, result, player_rank_change, player_rank_after_match)
				SELECT u.id, $2, $3, $4, $5, $6, $7, $8, $9
				FROM users u WHERE u.steam_id = $1
				ON CONFLICT (user_id, match_id) DO NOTHING
			`
			if err := tx.Exec(pmsQuery,
				profile.SteamID, match.ID, match.HeroName, match.Kills, match.Deaths, match.Assists,
				match.Result, match.PlayerRankChange, match.PlayerRankAfterMatch,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func calculatePerformanceDynamics(matches []domain.Match) domain.PerformanceDynamics {
	var dynamics domain.PerformanceDynamics
	if len(matches) < 2 {
		return dynamics
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].MatchTime.Before(matches[j].MatchTime)
	})

	// Rank Trend
	if len(matches) > 0 {
		rankStart := matches[0].PlayerRankAfterMatch - matches[0].PlayerRankChange
		rankEnd := matches[len(matches)-1].PlayerRankAfterMatch
		rankDiff := rankEnd - rankStart
		dynamics.Rank.Value = fmt.Sprintf("%+d Rank", rankDiff)
		dynamics.Rank.Trend = getTrend(float64(rankDiff))
		for _, m := range matches {
			dynamics.Rank.Sparkline = append(dynamics.Rank.Sparkline, float64(m.PlayerRankAfterMatch))
		}
	}

	// Win/Loss Trend
	netWins := 0
	for _, m := range matches {
		if m.Result == "Win" {
			netWins++
		} else {
			netWins--
		}
	}
	dynamics.WinLoss.Value = fmt.Sprintf("%+d WINS", netWins)
	dynamics.WinLoss.Trend = getTrend(float64(netWins))
	cumulativeWins := 0
	for _, m := range matches {
		if m.Result == "Win" {
			cumulativeWins++
		}
		dynamics.WinLoss.Sparkline = append(dynamics.WinLoss.Sparkline, float64(cumulativeWins))
	}

	kdaValues := []float64{}
	for _, m := range matches {
		kda := (float64(m.Kills) + float64(m.Assists)) / float64(m.Deaths)
		if m.Deaths == 0 {
			kda = float64(m.Kills) + float64(m.Assists)
		}
		kdaValues = append(kdaValues, kda)
	}
	avgKDA := kdaValues[len(kdaValues)-1]
	dynamics.KDA.Value = fmt.Sprintf("%.2f KDA", avgKDA)
	dynamics.KDA.Trend = getTrend(kdaValues[len(kdaValues)-1] - kdaValues[0])
	dynamics.KDA.Sparkline = kdaValues

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
