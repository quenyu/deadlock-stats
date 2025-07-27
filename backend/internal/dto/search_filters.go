package dto

import "fmt"

type SearchFilters struct {
	MinRank    int     `json:"min_rank"`
	MaxRank    int     `json:"max_rank"`
	MinMatches int     `json:"min_matches"`
	MaxMatches int     `json:"max_matches"`
	MinWinRate float64 `json:"min_win_rate"`
	MaxWinRate float64 `json:"max_win_rate"`
	MinKDRatio float64 `json:"min_kd_ratio"`
	MaxKDRatio float64 `json:"max_kd_ratio"`
	SortBy     string  `json:"sort_by"`    // "rank", "matches", "win_rate", "kd_ratio", "nickname", "created_at", "updated_at"
	SortOrder  string  `json:"sort_order"` // "asc", "desc"
}

func (f *SearchFilters) Validate() error {
	if f.MinRank > f.MaxRank && f.MaxRank != 0 {
		return fmt.Errorf("min_rank cannot be greater than max_rank")
	}
	if f.MinMatches > f.MaxMatches && f.MaxMatches != 0 {
		return fmt.Errorf("min_matches cannot be greater than max_matches")
	}
	if f.MinWinRate > f.MaxWinRate && f.MaxWinRate != 0 {
		return fmt.Errorf("min_win_rate cannot be greater than max_win_rate")
	}
	if f.MinKDRatio > f.MaxKDRatio && f.MaxKDRatio != 0 {
		return fmt.Errorf("min_kd_ratio cannot be greater than max_kd_ratio")
	}

	validSortBy := map[string]bool{
		"rank":       true,
		"matches":    true,
		"win_rate":   true,
		"kd_ratio":   true,
		"nickname":   true,
		"created_at": true,
		"updated_at": true,
	}
	if !validSortBy[f.SortBy] && f.SortBy != "" {
		return fmt.Errorf("invalid sort_by value: %s", f.SortBy)
	}

	validSortOrder := map[string]bool{
		"asc":  true,
		"desc": true,
	}
	if !validSortOrder[f.SortOrder] && f.SortOrder != "" {
		return fmt.Errorf("invalid sort_order value: %s", f.SortOrder)
	}

	return nil
}

func (f *SearchFilters) GetDefaultSortBy() string {
	if f.SortBy == "" {
		return "nickname"
	}
	return f.SortBy
}

func (f *SearchFilters) GetDefaultSortOrder() string {
	if f.SortOrder == "" {
		return "asc"
	}
	return f.SortOrder
}
