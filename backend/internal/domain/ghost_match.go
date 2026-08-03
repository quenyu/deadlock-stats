package domain

// GhostMatchReport compares one target match with a similar successful match.
// The first MVP intentionally uses only metrics available in regular match history.
type GhostMatchReport struct {
	TargetMatch    Match             `json:"target_match"`
	ReferenceMatch Match             `json:"reference_match"`
	Similarity     float64           `json:"similarity"`
	Confidence     string            `json:"confidence"`
	Comparisons    []MetricComparison `json:"comparisons"`
	Mission        TrainingMission   `json:"mission"`
	Limitations    []string          `json:"limitations"`
}

type MetricComparison struct {
	Metric         string  `json:"metric"`
	TargetValue    float64 `json:"target_value"`
	ReferenceValue float64 `json:"reference_value"`
	Delta          float64 `json:"delta"`
	Unit           string  `json:"unit"`
	Interpretation string  `json:"interpretation"`
}

type TrainingMission struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Target      float64 `json:"target"`
	Unit        string  `json:"unit"`
}
