package ghost

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/quenyu/deadlock-stats/internal/deadlock"
)

var (
	ErrTargetNotFound    = errors.New("target match not found in player history")
	ErrReferenceNotFound = errors.New("no comparable victory found on the same hero")
)

type Service struct {
	history deadlock.HistoryClient
}

func NewService(history deadlock.HistoryClient) *Service {
	return &Service{history: history}
}

type Report struct {
	Target         MatchSnapshot      `json:"target"`
	Reference      MatchSnapshot      `json:"reference"`
	Similarity     float64            `json:"similarity"`
	Confidence     string             `json:"confidence"`
	CandidateCount int                `json:"candidate_count"`
	Comparisons    []MetricComparison `json:"comparisons"`
	Mission        Mission            `json:"mission"`
	Limitations    []string           `json:"limitations"`
}

type MatchSnapshot struct {
	MatchID            int64   `json:"match_id"`
	HeroID             int     `json:"hero_id"`
	Won                bool    `json:"won"`
	DurationSeconds    int     `json:"duration_seconds"`
	NetWorth           int     `json:"net_worth"`
	SoulsPerMinute     float64 `json:"souls_per_minute"`
	KDA                float64 `json:"kda"`
	DeathsPer10Minutes float64 `json:"deaths_per_10_minutes"`
}

type MetricComparison struct {
	Key       string  `json:"key"`
	Label     string  `json:"label"`
	Target    float64 `json:"target"`
	Reference float64 `json:"reference"`
	DeltaPct  float64 `json:"delta_pct"`
}

type Mission struct {
	Metric      string  `json:"metric"`
	TargetValue float64 `json:"target_value"`
	Unit        string  `json:"unit"`
	Explanation string  `json:"explanation"`
}

func (s *Service) Build(ctx context.Context, accountID string, targetMatchID int64) (*Report, error) {
	matches, err := s.history.MatchHistory(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("load player history: %w", err)
	}

	var target *deadlock.Match
	for i := range matches {
		if matches[i].MatchID == targetMatchID {
			target = &matches[i]
			break
		}
	}
	if target == nil {
		return nil, ErrTargetNotFound
	}
	if target.MatchDurationS <= 0 {
		return nil, fmt.Errorf("target match has invalid duration: %d", target.MatchDurationS)
	}

	bestIndex := -1
	bestScore := -1.0
	candidateCount := 0

	for i := range matches {
		candidate := matches[i]
		if candidate.MatchID == target.MatchID || candidate.HeroID != target.HeroID || !candidate.Won() || candidate.MatchDurationS <= 0 {
			continue
		}

		candidateCount++
		score := similarity(*target, candidate)
		if score > bestScore {
			bestScore = score
			bestIndex = i
		}
	}

	if bestIndex < 0 {
		return nil, ErrReferenceNotFound
	}

	reference := matches[bestIndex]
	targetSnapshot := snapshot(*target)
	referenceSnapshot := snapshot(reference)

	return &Report{
		Target:         targetSnapshot,
		Reference:      referenceSnapshot,
		Similarity:     round(bestScore, 3),
		Confidence:     confidence(bestScore, candidateCount),
		CandidateCount: candidateCount,
		Comparisons:    comparisons(targetSnapshot, referenceSnapshot),
		Mission:        chooseMission(targetSnapshot, referenceSnapshot),
		Limitations: []string{
			"Reference is selected from the same player's available victories, not from a global rank cohort.",
			"Match-history data does not contain movement, item timing, lane assignment or decision intent.",
			"Differences are correlations and must not be presented as proven causes of a win or loss.",
		},
	}, nil
}

func snapshot(match deadlock.Match) MatchSnapshot {
	minutes := float64(match.MatchDurationS) / 60
	return MatchSnapshot{
		MatchID:            match.MatchID,
		HeroID:             match.HeroID,
		Won:                match.Won(),
		DurationSeconds:    match.MatchDurationS,
		NetWorth:           match.NetWorth,
		SoulsPerMinute:     round(safeDivide(float64(match.NetWorth), minutes), 1),
		KDA:                round(safeDivide(float64(match.PlayerKills+match.PlayerAssists), float64(max(match.PlayerDeaths, 1))), 2),
		DeathsPer10Minutes: round(safeDivide(float64(match.PlayerDeaths)*10, minutes), 2),
	}
}

func similarity(a, b deadlock.Match) float64 {
	aSnapshot := snapshot(a)
	bSnapshot := snapshot(b)

	duration := closeness(float64(a.MatchDurationS), float64(b.MatchDurationS))
	economy := closeness(aSnapshot.SoulsPerMinute, bSnapshot.SoulsPerMinute)
	kda := closeness(aSnapshot.KDA, bSnapshot.KDA)

	return clamp(duration*0.45+economy*0.40+kda*0.15, 0, 1)
}

func comparisons(target, reference MatchSnapshot) []MetricComparison {
	return []MetricComparison{
		metric("souls_per_minute", "Souls per minute", target.SoulsPerMinute, reference.SoulsPerMinute),
		metric("kda", "KDA", target.KDA, reference.KDA),
		metric("deaths_per_10_minutes", "Deaths per 10 minutes", target.DeathsPer10Minutes, reference.DeathsPer10Minutes),
		metric("duration_minutes", "Match duration", float64(target.DurationSeconds)/60, float64(reference.DurationSeconds)/60),
	}
}

func metric(key, label string, target, reference float64) MetricComparison {
	return MetricComparison{
		Key:       key,
		Label:     label,
		Target:    round(target, 2),
		Reference: round(reference, 2),
		DeltaPct:  round(safeDivide(target-reference, math.Abs(reference))*100, 1),
	}
}

func chooseMission(target, reference MatchSnapshot) Mission {
	economyDeficit := positiveRatio(reference.SoulsPerMinute-target.SoulsPerMinute, reference.SoulsPerMinute)
	kdaDeficit := positiveRatio(reference.KDA-target.KDA, reference.KDA)
	deathExcess := positiveRatio(target.DeathsPer10Minutes-reference.DeathsPer10Minutes, target.DeathsPer10Minutes)

	switch {
	case economyDeficit >= kdaDeficit && economyDeficit >= deathExcess && economyDeficit > 0:
		goal := math.Max(target.SoulsPerMinute*1.08, reference.SoulsPerMinute*0.90)
		return Mission{
			Metric:      "souls_per_minute",
			TargetValue: round(goal, 0),
			Unit:        "souls/min",
			Explanation: "Raise economy pace toward the comparable victory. The target is deliberately below the full reference value for a realistic next-match step.",
		}
	case deathExcess >= kdaDeficit && deathExcess > 0:
		goal := math.Max(0, (target.DeathsPer10Minutes+reference.DeathsPer10Minutes)/2)
		return Mission{
			Metric:      "deaths_per_10_minutes",
			TargetValue: round(goal, 2),
			Unit:        "deaths/10 min",
			Explanation: "Reduce death frequency halfway toward the comparable victory.",
		}
	case kdaDeficit > 0:
		goal := math.Min(reference.KDA, target.KDA*1.12)
		return Mission{
			Metric:      "kda",
			TargetValue: round(goal, 2),
			Unit:        "KDA",
			Explanation: "Improve combined kill and assist contribution relative to deaths by a modest measurable step.",
		}
	default:
		return Mission{
			Metric:      "souls_per_minute",
			TargetValue: round(target.SoulsPerMinute, 0),
			Unit:        "souls/min",
			Explanation: "No major deficit is visible in the available summary metrics. Repeat the current economy pace before adding deeper claims.",
		}
	}
}

func confidence(score float64, candidates int) string {
	switch {
	case score >= 0.85 && candidates >= 3:
		return "high"
	case score >= 0.70:
		return "medium"
	default:
		return "low"
	}
}

func closeness(a, b float64) float64 {
	denominator := math.Max(math.Abs(a), math.Abs(b))
	if denominator == 0 {
		return 1
	}
	return clamp(1-math.Abs(a-b)/denominator, 0, 1)
}

func positiveRatio(value, denominator float64) float64 {
	if value <= 0 || denominator <= 0 {
		return 0
	}
	return value / denominator
}

func safeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func clamp(value, low, high float64) float64 {
	return math.Max(low, math.Min(high, value))
}

func round(value float64, digits int) float64 {
	factor := math.Pow(10, float64(digits))
	return math.Round(value*factor) / factor
}
