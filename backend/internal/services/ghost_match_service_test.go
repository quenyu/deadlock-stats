package services

import (
	"errors"
	"testing"

	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestGhostMatch_SelectsClosestSuccessfulSameHero(t *testing.T) {
	service := NewGhostMatchService()
	matches := []domain.Match{
		{ID: "loss", HeroID: 7, PlayerKills: 5, PlayerDeaths: 8, PlayerAssists: 9, NetWorth: 21000, MatchDurationS: 1800, Result: "loss", PlayerRankAfterMatch: 36},
		{ID: "close-win", HeroID: 7, PlayerKills: 8, PlayerDeaths: 5, PlayerAssists: 12, NetWorth: 27000, MatchDurationS: 1860, Result: "win", PlayerRankAfterMatch: 37},
		{ID: "short-win", HeroID: 7, PlayerKills: 12, PlayerDeaths: 1, PlayerAssists: 10, NetWorth: 18000, MatchDurationS: 900, Result: "win", PlayerRankAfterMatch: 55},
		{ID: "other-hero", HeroID: 9, PlayerKills: 8, PlayerDeaths: 4, PlayerAssists: 8, NetWorth: 26000, MatchDurationS: 1800, Result: "win", PlayerRankAfterMatch: 36},
	}

	report, err := service.BuildReport(matches, "loss")
	require.NoError(t, err)
	require.Equal(t, "close-win", report.ReferenceMatch.ID)
	require.Greater(t, report.Similarity, 0.6)
	require.Equal(t, "raise_economy_tempo", report.Mission.Code)
	require.Len(t, report.Comparisons, 3)
}

func TestGhostMatch_FallsBackToTeamResult(t *testing.T) {
	service := NewGhostMatchService()
	matches := []domain.Match{
		{ID: "target", HeroID: 1, PlayerTeam: 1, MatchResult: 2, MatchDurationS: 1200, NetWorth: 15000},
		{ID: "win", HeroID: 1, PlayerTeam: 1, MatchResult: 1, MatchDurationS: 1200, NetWorth: 16000},
	}

	report, err := service.BuildReport(matches, "target")
	require.NoError(t, err)
	require.Equal(t, "win", report.ReferenceMatch.ID)
}

func TestGhostMatch_ReturnsErrorWithoutComparableWin(t *testing.T) {
	service := NewGhostMatchService()
	matches := []domain.Match{{ID: "target", HeroID: 1, Result: "loss"}}

	_, err := service.BuildReport(matches, "target")
	require.True(t, errors.Is(err, ErrGhostReferenceAbsent))
}

func TestGhostMatch_ReturnsErrorForUnknownTarget(t *testing.T) {
	service := NewGhostMatchService()
	_, err := service.BuildReport(nil, "missing")
	require.True(t, errors.Is(err, ErrGhostTargetNotFound))
}
