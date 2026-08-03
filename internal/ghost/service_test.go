package ghost

import (
	"context"
	"errors"
	"testing"

	"github.com/quenyu/deadlock-stats/internal/deadlock"
)

type fakeHistory struct {
	matches []deadlock.Match
	err     error
}

func (f fakeHistory) MatchHistory(context.Context, string) ([]deadlock.Match, error) {
	return f.matches, f.err
}

func TestBuildSelectsClosestVictoryOnSameHero(t *testing.T) {
	matches := []deadlock.Match{
		{MatchID: 100, HeroID: 7, PlayerTeam: 0, MatchResult: 1, MatchDurationS: 1800, NetWorth: 21000, PlayerKills: 5, PlayerDeaths: 8, PlayerAssists: 9},
		// Team 0 victory. This catches the old bug where result == 0 was called a loss.
		{MatchID: 200, HeroID: 7, PlayerTeam: 0, MatchResult: 0, MatchDurationS: 1850, NetWorth: 24500, PlayerKills: 8, PlayerDeaths: 5, PlayerAssists: 11},
		{MatchID: 300, HeroID: 7, PlayerTeam: 1, MatchResult: 1, MatchDurationS: 3100, NetWorth: 51000, PlayerKills: 20, PlayerDeaths: 3, PlayerAssists: 20},
		{MatchID: 400, HeroID: 8, PlayerTeam: 0, MatchResult: 0, MatchDurationS: 1810, NetWorth: 24000},
		{MatchID: 500, HeroID: 7, PlayerTeam: 1, MatchResult: 0, MatchDurationS: 1810, NetWorth: 24000},
	}

	report, err := NewService(fakeHistory{matches: matches}).Build(context.Background(), "123", 100)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if report.Reference.MatchID != 200 {
		t.Fatalf("reference match = %d, want 200", report.Reference.MatchID)
	}
	if report.Target.Won {
		t.Fatal("target match must be classified as a loss")
	}
	if !report.Reference.Won {
		t.Fatal("team 0 reference must be classified as a win")
	}
	if report.CandidateCount != 2 {
		t.Fatalf("candidate count = %d, want 2", report.CandidateCount)
	}
}

func TestBuildReturnsNoReferenceWhenPlayerHasNoSameHeroVictory(t *testing.T) {
	matches := []deadlock.Match{
		{MatchID: 100, HeroID: 7, PlayerTeam: 0, MatchResult: 1, MatchDurationS: 1800},
		{MatchID: 200, HeroID: 8, PlayerTeam: 0, MatchResult: 0, MatchDurationS: 1800},
	}

	_, err := NewService(fakeHistory{matches: matches}).Build(context.Background(), "123", 100)
	if !errors.Is(err, ErrReferenceNotFound) {
		t.Fatalf("error = %v, want ErrReferenceNotFound", err)
	}
}

func TestBuildCreatesEconomyMissionForLargeEconomyDeficit(t *testing.T) {
	matches := []deadlock.Match{
		{MatchID: 100, HeroID: 7, PlayerTeam: 0, MatchResult: 1, MatchDurationS: 1800, NetWorth: 15000, PlayerKills: 7, PlayerDeaths: 5, PlayerAssists: 10},
		{MatchID: 200, HeroID: 7, PlayerTeam: 0, MatchResult: 0, MatchDurationS: 1800, NetWorth: 30000, PlayerKills: 7, PlayerDeaths: 5, PlayerAssists: 10},
	}

	report, err := NewService(fakeHistory{matches: matches}).Build(context.Background(), "123", 100)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if report.Mission.Metric != "souls_per_minute" {
		t.Fatalf("mission metric = %q, want souls_per_minute", report.Mission.Metric)
	}
	if report.Mission.TargetValue <= report.Target.SoulsPerMinute {
		t.Fatalf("mission target = %.2f, must exceed current %.2f", report.Mission.TargetValue, report.Target.SoulsPerMinute)
	}
}

func TestBuildPropagatesHistoryError(t *testing.T) {
	upstreamErr := errors.New("upstream unavailable")
	_, err := NewService(fakeHistory{err: upstreamErr}).Build(context.Background(), "123", 100)
	if !errors.Is(err, upstreamErr) {
		t.Fatalf("error = %v, want wrapped upstream error", err)
	}
}
