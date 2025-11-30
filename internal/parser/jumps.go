package parser

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// SDEStargateDestination represents the destination of a stargate.
type SDEStargateDestination struct {
	SolarSystemID int64 `yaml:"solarSystemID"`
	StargateID    int64 `yaml:"stargateID"`
}

// SDEMapStargate represents a stargate in the flat SDE format.
type SDEMapStargate struct {
	SolarSystemID int64                  `yaml:"solarSystemID"`
	Destination   SDEStargateDestination `yaml:"destination"`
	TypeID        int64                  `yaml:"typeID,omitempty"`
}

// ParseStargates parses the mapStargates.yaml file and extracts system jumps.
func (p *Parser) ParseStargates() ([]models.SystemJump, error) {
	path := p.filePath("mapStargates.yaml")

	// Parse the file as a map of stargate ID to stargate data
	rawStargates, err := yaml.ParseFileMap[int64, SDEMapStargate](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stargates file: %w", err)
	}

	// Create unique jump pairs (only count each connection once)
	// Use a map to deduplicate - store with smaller ID first
	jumpSet := make(map[[2]int64]struct{})

	for _, data := range rawStargates {
		fromSystem := data.SolarSystemID
		toSystem := data.Destination.SolarSystemID

		if fromSystem == 0 || toSystem == 0 {
			// Invalid data, skip
			continue
		}

		// Create ordered pair (smaller ID first) to avoid duplicates
		var pair [2]int64
		if fromSystem < toSystem {
			pair = [2]int64{fromSystem, toSystem}
		} else {
			pair = [2]int64{toSystem, fromSystem}
		}
		jumpSet[pair] = struct{}{}
	}

	// Convert to slice
	jumps := make([]models.SystemJump, 0, len(jumpSet))
	for pair := range jumpSet {
		jumps = append(jumps, models.SystemJump{
			FromSolarSystemID: pair[0],
			ToSolarSystemID:   pair[1],
		})
	}

	// Sort by from system ID, then to system ID for consistent output
	sort.Slice(jumps, func(i, j int) bool {
		if jumps[i].FromSolarSystemID != jumps[j].FromSolarSystemID {
			return jumps[i].FromSolarSystemID < jumps[j].FromSolarSystemID
		}
		return jumps[i].ToSolarSystemID < jumps[j].ToSolarSystemID
	})

	return jumps, nil
}
