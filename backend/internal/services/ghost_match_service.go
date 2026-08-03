package services

import (
	"errors"
	"math"
	"sort"
	"strings"

	"github.com/quenyu/deadlock-stats/internal/domain"
)

var (
	ErrGhostTargetNotFound  = errors.New("target match not found")
	ErrGhostReferenceAbsent = errors.New("no comparable successful match found")
)

type GhostMatchService struct{}

func NewGhostMatchService() *GhostMatchService { return &GhostMatchService{} }

// BuildReport selects a successful match with the same hero and the closest
// duration/performance profile. In the first MVP candidates come from the same
// player's history. A global cohort can later be supplied without changing the
// report contract.
func (s *GhostMatchService) BuildReport(matches []domain.Match, targetID string) (*domain.GhostMatchReport, error) {
	target, ok := findMatch(matches, targetID)
	if !ok {
		return nil, ErrGhostTargetNotFound
	}

	candidates := make([]scoredMatch, 0)
	for _, candidate := range matches {
		if candidate.ID == target.ID || candidate.HeroID != target.HeroID || !isWin(candidate) {
			continue
		}
		score := similarity(target, candidate)
		candidates = append(candidates, scoredMatch{match: candidate, score: score})
	}
	if len(candidates) == 0 {
		return nil, ErrGhostReferenceAbsent
	}

	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].score > candidates[j].score })
	reference := candidates[0]
	comparisons := compareMatches(target, reference.match)

	return &domain.GhostMatchReport{
		TargetMatch:     target,
		ReferenceMatch:  reference.match,
		Similarity:      round(reference.score, 3),
		Confidence:      confidence(reference.score, len(candidates)),
		Comparisons:     comparisons,
		Mission:         chooseMission(target, reference.match),
		Limitations: []string{
			"Сравнение основано на итоговых метриках истории матчей, а не на намерениях игрока.",
			"Первая версия подбирает успешный матч из истории того же игрока.",
			"Для анализа перемещений, драк и таймингов покупок нужны расширенные metadata/demo-данные.",
		},
	}, nil
}

type scoredMatch struct {
	match domain.Match
	score float64
}

func findMatch(matches []domain.Match, id string) (domain.Match, bool) {
	for _, match := range matches {
		if match.ID == id {
			return match, true
		}
	}
	return domain.Match{}, false
}

func isWin(match domain.Match) bool {
	result := strings.ToLower(strings.TrimSpace(match.Result))
	if result == "win" || result == "victory" || result == "победа" {
		return true
	}
	if result == "loss" || result == "defeat" || result == "поражение" {
		return false
	}
	return match.PlayerTeam != 0 && match.MatchResult == match.PlayerTeam
}

func similarity(a, b domain.Match) float64 {
	duration := closeness(float64(a.MatchDurationS), float64(b.MatchDurationS), 900)
	netWorthPerMinute := closeness(rate(a.NetWorth, a.MatchDurationS), rate(b.NetWorth, b.MatchDurationS), 500)
	kda := closeness(kdaValue(a), kdaValue(b), 3)
	rank := closeness(float64(a.PlayerRankAfterMatch), float64(b.PlayerRankAfterMatch), 12)

	// Duration and economy are the strongest currently available signals.
	return clamp(0.35*duration+0.35*netWorthPerMinute+0.20*kda+0.10*rank, 0, 1)
}

func compareMatches(target, reference domain.Match) []domain.MetricComparison {
	targetNPM := rate(target.NetWorth, target.MatchDurationS)
	refNPM := rate(reference.NetWorth, reference.MatchDurationS)
	targetKDA := kdaValue(target)
	refKDA := kdaValue(reference)
	targetDeathsPer10 := perTen(target.PlayerDeaths, target.MatchDurationS)
	refDeathsPer10 := perTen(reference.PlayerDeaths, reference.MatchDurationS)

	return []domain.MetricComparison{
		metric("net_worth_per_minute", targetNPM, refNPM, "souls/min", lowerHigher(targetNPM, refNPM, "Темп экономики ниже успешного аналога.", "Темп экономики не хуже успешного аналога.")),
		metric("kda", targetKDA, refKDA, "ratio", lowerHigher(targetKDA, refKDA, "Участие в убийствах относительно смертей ниже.", "KDA не хуже успешного аналога.")),
		metric("deaths_per_10_minutes", targetDeathsPer10, refDeathsPer10, "deaths/10m", higherWorse(targetDeathsPer10, refDeathsPer10, "Частота смертей выше успешного аналога.", "Частота смертей не выше успешного аналога.")),
	}
}

func chooseMission(target, reference domain.Match) domain.TrainingMission {
	targetNPM := rate(target.NetWorth, target.MatchDurationS)
	refNPM := rate(reference.NetWorth, reference.MatchDurationS)
	if targetNPM < refNPM*0.92 {
		goal := math.Max(targetNPM*1.08, refNPM*0.95)
		return domain.TrainingMission{
			Code: "raise_economy_tempo",
			Title: "Поднять темп экономики",
			Description: "В следующем матче на этом герое достигни указанного темпа душ в минуту. Это проверяемая цель; она не утверждает, что экономика была единственной причиной поражения.",
			Target: round(goal, 0), Unit: "souls/min",
		}
	}

	targetDeaths := perTen(target.PlayerDeaths, target.MatchDurationS)
	refDeaths := perTen(reference.PlayerDeaths, reference.MatchDurationS)
	if targetDeaths > refDeaths+0.5 {
		return domain.TrainingMission{
			Code: "reduce_death_rate",
			Title: "Снизить частоту смертей",
			Description: "В следующем матче удерживай число смертей на 10 минут не выше цели.",
			Target: round(math.Max(refDeaths, targetDeaths-0.75), 1), Unit: "deaths/10m",
		}
	}

	return domain.TrainingMission{
		Code: "match_reference_kda",
		Title: "Повторить уровень участия",
		Description: "В следующем матче достигни KDA не ниже успешного аналога.",
		Target: round(kdaValue(reference), 2), Unit: "KDA",
	}
}

func metric(name string, target, reference float64, unit, interpretation string) domain.MetricComparison {
	return domain.MetricComparison{Metric: name, TargetValue: round(target, 2), ReferenceValue: round(reference, 2), Delta: round(target-reference, 2), Unit: unit, Interpretation: interpretation}
}

func rate(value, seconds int) float64 {
	if seconds <= 0 { return 0 }
	return float64(value) / (float64(seconds) / 60)
}

func perTen(value, seconds int) float64 {
	if seconds <= 0 { return 0 }
	return float64(value) / (float64(seconds) / 600)
}

func kdaValue(match domain.Match) float64 {
	deaths := match.PlayerDeaths
	if deaths < 1 { deaths = 1 }
	return float64(match.PlayerKills+match.PlayerAssists) / float64(deaths)
}

func closeness(a, b, scale float64) float64 {
	if scale <= 0 { return 0 }
	return clamp(1-math.Abs(a-b)/scale, 0, 1)
}

func confidence(score float64, candidates int) string {
	if score >= 0.82 && candidates >= 3 { return "high" }
	if score >= 0.62 { return "medium" }
	return "low"
}

func lowerHigher(target, reference float64, lower, otherwise string) string {
	if target < reference { return lower }
	return otherwise
}

func higherWorse(target, reference float64, higher, otherwise string) string {
	if target > reference { return higher }
	return otherwise
}

func clamp(value, min, max float64) float64 {
	if value < min { return min }
	if value > max { return max }
	return value
}

func round(value float64, precision int) float64 {
	factor := math.Pow10(precision)
	return math.Round(value*factor) / factor
}
