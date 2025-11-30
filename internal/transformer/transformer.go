package transformer

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/internal/parser"
)

// Transformer handles transformation of parsed SDE data to Wanderer format.
type Transformer struct {
	config *config.Config
}

// New creates a new Transformer with the given configuration.
func New(cfg *config.Config) *Transformer {
	return &Transformer{
		config: cfg,
	}
}

// Transform converts parsed SDE data into Wanderer's output format.
func (t *Transformer) Transform(parseResult *parser.ParseResult) (*models.ConvertedData, error) {
	if t.config.Verbose {
		fmt.Println("Transforming SDE data...")
	}

	// Transform solar systems with security calculation
	if t.config.Verbose {
		fmt.Println("  Transforming solar systems...")
	}
	systems := t.transformSolarSystems(parseResult.SolarSystems)

	// Sort regions for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting regions...")
	}
	regions := t.sortRegions(parseResult.Regions)

	// Sort constellations for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting constellations...")
	}
	constellations := t.sortConstellations(parseResult.Constellations)

	// Filter to ship types only
	if t.config.Verbose {
		fmt.Println("  Filtering ship types...")
	}
	shipTypes := FilterShipTypes(parseResult.Types, parseResult.Groups)

	// Filter to ship groups only
	if t.config.Verbose {
		fmt.Println("  Filtering ship groups...")
	}
	shipGroups := FilterShipGroups(parseResult.Groups)

	// Sort wormhole classes for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting wormhole classes...")
	}
	wormholeClasses := t.sortWormholeClasses(parseResult.WormholeClasses)

	// Sort system jumps for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting system jumps...")
	}
	systemJumps := t.sortSystemJumps(parseResult.SystemJumps)

	result := &models.ConvertedData{
		Universe: &models.UniverseData{
			Regions:        regions,
			Constellations: constellations,
			SolarSystems:   systems,
		},
		ShipTypes:       shipTypes,
		ItemGroups:      shipGroups,
		WormholeClasses: wormholeClasses,
		SystemJumps:     systemJumps,
	}

	if t.config.Verbose {
		fmt.Printf("Transformation complete:\n")
		fmt.Printf("  Regions:         %d\n", len(result.Universe.Regions))
		fmt.Printf("  Constellations:  %d\n", len(result.Universe.Constellations))
		fmt.Printf("  Solar Systems:   %d\n", len(result.Universe.SolarSystems))
		fmt.Printf("  Ship Types:      %d\n", len(result.ShipTypes))
		fmt.Printf("  Ship Groups:     %d\n", len(result.ItemGroups))
		fmt.Printf("  Wormhole Classes: %d\n", len(result.WormholeClasses))
		fmt.Printf("  System Jumps:    %d\n", len(result.SystemJumps))
	}

	return result, nil
}

// transformSolarSystems applies security calculation to solar systems.
func (t *Transformer) transformSolarSystems(systems []models.SolarSystem) []models.SolarSystem {
	result := make([]models.SolarSystem, len(systems))

	for i, sys := range systems {
		result[i] = models.SolarSystem{
			SolarSystemID:   sys.SolarSystemID,
			RegionID:        sys.RegionID,
			ConstellationID: sys.ConstellationID,
			SolarSystemName: sys.SolarSystemName,
			SunTypeID:       sys.SunTypeID,
			Security:        GetTrueSecurity(sys.Security),
		}
	}

	// Sort by system ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].SolarSystemID < result[j].SolarSystemID
	})

	return result
}

// sortRegions returns regions sorted by ID.
func (t *Transformer) sortRegions(regions []models.Region) []models.Region {
	result := make([]models.Region, len(regions))
	copy(result, regions)

	sort.Slice(result, func(i, j int) bool {
		return result[i].RegionID < result[j].RegionID
	})

	return result
}

// sortConstellations returns constellations sorted by ID.
func (t *Transformer) sortConstellations(constellations []models.Constellation) []models.Constellation {
	result := make([]models.Constellation, len(constellations))
	copy(result, constellations)

	sort.Slice(result, func(i, j int) bool {
		return result[i].ConstellationID < result[j].ConstellationID
	})

	return result
}

// sortWormholeClasses returns wormhole classes sorted by location ID.
func (t *Transformer) sortWormholeClasses(classes []models.WormholeClassLocation) []models.WormholeClassLocation {
	result := make([]models.WormholeClassLocation, len(classes))
	copy(result, classes)

	sort.Slice(result, func(i, j int) bool {
		return result[i].LocationID < result[j].LocationID
	})

	return result
}

// sortSystemJumps returns system jumps sorted by from system ID, then to system ID.
func (t *Transformer) sortSystemJumps(jumps []models.SystemJump) []models.SystemJump {
	result := make([]models.SystemJump, len(jumps))
	copy(result, jumps)

	sort.Slice(result, func(i, j int) bool {
		if result[i].FromSolarSystemID != result[j].FromSolarSystemID {
			return result[i].FromSolarSystemID < result[j].FromSolarSystemID
		}
		return result[i].ToSolarSystemID < result[j].ToSolarSystemID
	})

	return result
}

// Validate performs validation checks on the converted data.
func (t *Transformer) Validate(data *models.ConvertedData) *models.ValidationResult {
	result := &models.ValidationResult{
		SolarSystems:    len(data.Universe.SolarSystems),
		Regions:         len(data.Universe.Regions),
		Constellations:  len(data.Universe.Constellations),
		ShipTypes:       len(data.ShipTypes),
		ItemGroups:      len(data.ItemGroups),
		SystemJumps:     len(data.SystemJumps),
		WormholeClasses: len(data.WormholeClasses),
	}

	// Validation thresholds based on known EVE universe size
	const (
		minSolarSystems   = 8000
		minRegions        = 100
		minConstellations = 1000
		minShipTypes      = 400
		minSystemJumps    = 10000
	)

	// Check minimum counts
	if result.SolarSystems < minSolarSystems {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Solar system count (%d) is below expected minimum (%d)",
				result.SolarSystems, minSolarSystems))
	}

	if result.Regions < minRegions {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Region count (%d) is below expected minimum (%d)",
				result.Regions, minRegions))
	}

	if result.Constellations < minConstellations {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Constellation count (%d) is below expected minimum (%d)",
				result.Constellations, minConstellations))
	}

	if result.ShipTypes < minShipTypes {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Ship type count (%d) is below expected minimum (%d)",
				result.ShipTypes, minShipTypes))
	}

	if result.SystemJumps < minSystemJumps {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("System jump count (%d) is below expected minimum (%d)",
				result.SystemJumps, minSystemJumps))
	}

	// Check for empty required data
	if result.SolarSystems == 0 {
		result.Errors = append(result.Errors, "No solar systems found")
	}

	if result.Regions == 0 {
		result.Errors = append(result.Errors, "No regions found")
	}

	if result.Constellations == 0 {
		result.Errors = append(result.Errors, "No constellations found")
	}

	return result
}
