// Copyright (c) 2022 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

// StatisticsConfig holds configuration for statistic-based matchmaking
type StatisticsConfig struct {
	// Statistics is the list of valid stat codes (e.g., ["mmr_ryu", "mmr_ken", "rank_score"])
	Statistics []string `json:"statistics"`

	// SelectedStatKey is the attribute key players use to specify which stat to use
	// Default: "selected_stat"
	SelectedStatKey string `json:"selected_stat_key"`

	// EnrichedKey is the ticket attribute key where the selected stat value is stored after enrichment
	// Default: "mmr"
	EnrichedKey string `json:"enriched_key"`

	// DefaultValue is the value to use if player doesn't have the selected stat
	// If 0, validation will fail for missing stat
	DefaultValue float64 `json:"default_value"`
}

// GetSelectedStatKey returns the key for selected stat, defaulting to "selected_stat"
func (c StatisticsConfig) GetSelectedStatKey() string {
	if c.SelectedStatKey == "" {
		return "selected_stat"
	}

	return c.SelectedStatKey
}

// GetEnrichedKey returns the enriched key, defaulting to "mmr"
func (c StatisticsConfig) GetEnrichedKey() string {
	if c.EnrichedKey == "" {
		return "mmr"
	}

	return c.EnrichedKey
}

// IsValidStat checks if a stat code is in the allowed list
func (c StatisticsConfig) IsValidStat(statCode string) bool {
	for _, validStat := range c.Statistics {
		if validStat == statCode {
			return true
		}
	}

	return false
}

// GameRules defines the matchmaking rules parsed from JSON
type GameRules struct {
	Statistics StatisticsConfig `json:"statistics_config"`
}
