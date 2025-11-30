// Package transformer provides data transformation logic for converting SDE data
// to Wanderer's format.
package transformer

import "math"

// TruncateToTwoDigits truncates a float to 2 decimal places.
func TruncateToTwoDigits(value float64) float64 {
	return math.Floor(value*100) / 100
}

// GetTrueSecurity calculates the display security status matching Wanderer's logic.
// EVE's security status display rounds in specific ways that this function replicates.
//
// The logic handles these cases:
// - Values between 0.0 and 0.05 (exclusive) round up to 0.1
// - Other values are truncated to 2 decimals, then rounded based on the second decimal
func GetTrueSecurity(security float64) float64 {
	// Special case: very low positive security rounds up to 0.1
	if security > 0.0 && security < 0.05 {
		return math.Ceil(security*10) / 10
	}

	// Truncate to 2 decimal places
	truncated := TruncateToTwoDigits(security)

	// Get the floor at 1 decimal place
	floor := math.Floor(truncated*10) / 10

	// Calculate the difference (second decimal digit contribution)
	diff := math.Round((truncated-floor)*100) / 100

	// If the second decimal is less than 0.05, round down; otherwise round up
	if diff < 0.05 {
		return floor
	}
	return math.Ceil(truncated*10) / 10
}

// RoundSecurity rounds security to one decimal place for display.
func RoundSecurity(security float64) float64 {
	return math.Round(security*10) / 10
}
