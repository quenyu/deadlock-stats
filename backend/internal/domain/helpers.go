package domain

import (
	"fmt"
	"sort"
)

var rankScores = []float64{
	0, 11, 12, 13, 14, 15, 16, 21, 22, 23, 24, 25, 26, 31, 32, 33, 34, 35, 36, 41, 42, 43, 44, 45, 46, 51, 52, 53, 54, 55, 56, 61, 62, 63, 64, 65, 66, 71, 72, 73, 74, 75, 76, 81, 82, 83, 84, 85, 86, 91, 92, 93, 94, 95, 96, 101, 102, 103, 104, 105, 106, 111, 112, 113, 114, 115, 116,
}

func GetRankFromScore(playerScore float64) int {
	index := int(playerScore)

	if index < 0 {
		index = 0
	}
	if index >= len(rankScores) {
		index = len(rankScores) - 1
	}

	return int(rankScores[index])
}

func CalculatePerformanceDynamics(matches []Match) PerformanceDynamics {
	var dynamics PerformanceDynamics
	if len(matches) < 2 {
		return dynamics
	}

	sortedMatches := make([]Match, len(matches))
	copy(sortedMatches, matches)

	sort.SliceStable(sortedMatches, func(i, j int) bool {
		return sortedMatches[i].MatchTime.Before(sortedMatches[j].MatchTime)
	})

	firstRankedIdx := -1
	lastRankedIdx := -1
	for i, m := range sortedMatches {
		if m.PlayerRankAfterMatch != 0 {
			if firstRankedIdx == -1 {
				firstRankedIdx = i
			}
			lastRankedIdx = i
		}
	}

	if firstRankedIdx == -1 || lastRankedIdx == firstRankedIdx {
		return dynamics
	}

	firstRankedMatch := sortedMatches[firstRankedIdx]
	lastRankedMatch := sortedMatches[lastRankedIdx]

	rankEnd := lastRankedMatch.PlayerRankAfterMatch
	var rankStart int

	if firstRankedMatch.PlayerRankChange == 0 {
		rankStart = firstRankedMatch.PlayerRankAfterMatch
	} else {
		rankStart = firstRankedMatch.PlayerRankAfterMatch - firstRankedMatch.PlayerRankChange
	}

	if rankEnd != 0 && rankStart != 0 {
		rankDiff := rankEnd - rankStart
		dynamics.Rank.Value = fmt.Sprintf("%+d Rank", rankDiff)
		dynamics.Rank.Trend = GetTrend(float64(rankDiff))
		for i := firstRankedIdx; i <= lastRankedIdx; i++ {
			dynamics.Rank.Sparkline = append(dynamics.Rank.Sparkline, float64(sortedMatches[i].PlayerRankAfterMatch))
		}
	}

	var winLossTrend Trend
	var kdaTrend Trend

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
		winLossTrend.Trend = GetTrend(float64(wins - losses))

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

func GetTrend(value float64) string {
	if value > 0 {
		return "up"
	}
	if value < 0 {
		return "down"
	}
	return "stable"
}

func MapMatchResult(result int) string {
	if result == 1 {
		return "Win"
	}
	return "Loss"
}

func FindPeakRank(mmrHistory []DeadlockMMR, getRankNameAndSubRank func(int) (string, int, string), getRankImageURL func(int, int) string) (int, string, string) {
	if len(mmrHistory) == 0 {
		return 0, "", ""
	}

	var peakRank int
	for _, mmr := range mmrHistory {
		rankScore := GetRankFromScore(mmr.PlayerScore)
		if rankScore > peakRank {
			peakRank = rankScore
		}
	}

	if peakRank > 0 {
		tier := peakRank / 10
		subTier := peakRank % 10
		rankName, _, _ := getRankNameAndSubRank(tier)
		rankImage := getRankImageURL(tier, subTier)
		return peakRank, rankName, rankImage
	}

	return 0, "Unranked", ""
}

func CalculatePersonalRecords(matches []Match) PersonalRecords {
	var records PersonalRecords

	for _, match := range matches {
		if match.PlayerKills > records.MaxKills {
			records.MaxKills = match.PlayerKills
			records.MaxKillsMatchID = match.ID
		}
		if match.PlayerAssists > records.MaxAssists {
			records.MaxAssists = match.PlayerAssists
			records.MaxAssistsMatchID = match.ID
		}

		if match.NetWorth > records.MaxNetWorth {
			records.MaxNetWorth = match.NetWorth
			records.MaxNetWorthMatchID = match.ID
		}

		kda := 0.0
		if match.PlayerDeaths > 0 {
			kda = float64(match.PlayerKills+match.PlayerAssists) / float64(match.PlayerDeaths)
		} else {
			kda = float64(match.PlayerKills + match.PlayerAssists)
		}

		if kda > records.BestKDA {
			records.BestKDA = kda
			records.BestKDAMatchID = match.ID
		}
	}

	return records
}

type AverageStats struct {
	AvgKills    float64
	AvgDeaths   float64
	AvgAssists  float64
	AvgDuration float64
}

func CalculateAverageStats(matches []Match, totalMatches int) AverageStats {
	var stats AverageStats

	if totalMatches == 0 {
		return stats
	}

	var totalKills, totalDeaths, totalAssists, totalDuration int
	for _, match := range matches {
		totalKills += match.PlayerKills
		totalDeaths += match.PlayerDeaths
		totalAssists += match.PlayerAssists
		totalDuration += match.MatchDurationS
	}

	stats.AvgKills = float64(totalKills) / float64(totalMatches)
	stats.AvgDeaths = float64(totalDeaths) / float64(totalMatches)
	stats.AvgAssists = float64(totalAssists) / float64(totalMatches)
	stats.AvgDuration = float64(totalDuration) / float64(totalMatches) / 60 // minutes

	return stats
}
