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

	// Transform all types
	if t.config.Verbose {
		fmt.Println("  Transforming types...")
	}
	invTypes := t.transformTypes(parseResult.Types)

	// Transform all groups
	if t.config.Verbose {
		fmt.Println("  Transforming groups...")
	}
	invGroups := t.transformGroups(parseResult.Groups)

	// Sort wormhole classes for consistent output
	if t.config.Verbose {
		fmt.Println("  Sorting wormhole classes...")
	}
	wormholeClasses := t.sortWormholeClasses(parseResult.WormholeClasses)

	// Transform system jumps with region/constellation lookup
	if t.config.Verbose {
		fmt.Println("  Transforming system jumps...")
	}
	systemJumps := t.transformSystemJumps(parseResult.SystemJumps, systems)

	result := &models.ConvertedData{
		Universe: &models.UniverseData{
			Regions:        regions,
			Constellations: constellations,
			SolarSystems:   systems,
		},
		InvTypes:        invTypes,
		InvGroups:       invGroups,
		WormholeClasses: wormholeClasses,
		SystemJumps:     systemJumps,
	}

	if t.config.Verbose {
		fmt.Printf("Transformation complete:\n")
		fmt.Printf("  Regions:         %d\n", len(result.Universe.Regions))
		fmt.Printf("  Constellations:  %d\n", len(result.Universe.Constellations))
		fmt.Printf("  Solar Systems:   %d\n", len(result.Universe.SolarSystems))
		fmt.Printf("  Types:           %d\n", len(result.InvTypes))
		fmt.Printf("  Groups:          %d\n", len(result.InvGroups))
		fmt.Printf("  Wormhole Classes: %d\n", len(result.WormholeClasses))
		fmt.Printf("  System Jumps:    %d\n", len(result.SystemJumps))
	}

	return result, nil
}

// transformSolarSystems applies security calculation to solar systems while preserving all fields.
func (t *Transformer) transformSolarSystems(systems []models.SolarSystem) []models.SolarSystem {
	result := make([]models.SolarSystem, len(systems))

	for i, sys := range systems {
		// Copy all fields, applying security transformation
		result[i] = sys
		result[i].Security = GetTrueSecurity(sys.Security)
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

// transformTypes converts SDE types to InvType format.
func (t *Transformer) transformTypes(types map[int64]models.SDEType) []models.InvType {
	result := make([]models.InvType, 0, len(types))

	for typeID, sdeType := range types {
		invType := models.InvType{
			TypeID:        typeID,
			GroupID:       sdeType.GroupID,
			TypeName:      sdeType.Name["en"],
			Description:   sdeType.Description["en"],
			Mass:          sdeType.Mass,
			Volume:        sdeType.Volume,
			Capacity:      sdeType.Capacity,
			PortionSize:   sdeType.PortionSize,
			RaceID:        models.Int64Ptr(sdeType.RaceID),
			BasePrice:     sdeType.BasePrice,
			Published:     sdeType.Published,
			MarketGroupID: models.Int64Ptr(sdeType.MarketGroupID),
			IconID:        models.Int64Ptr(sdeType.IconID),
			SoundID:       models.Int64Ptr(sdeType.SoundID),
			GraphicID:     models.Int64Ptr(sdeType.GraphicID),
		}
		result = append(result, invType)
	}

	// Sort by type ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].TypeID < result[j].TypeID
	})

	return result
}

// transformGroups converts SDE groups to InvGroup format.
func (t *Transformer) transformGroups(groups map[int64]models.SDEGroup) []models.InvGroup {
	result := make([]models.InvGroup, 0, len(groups))

	for groupID, sdeGroup := range groups {
		invGroup := models.InvGroup{
			GroupID:              groupID,
			CategoryID:           sdeGroup.CategoryID,
			GroupName:            sdeGroup.Name["en"],
			IconID:               models.Int64Ptr(sdeGroup.IconID),
			UseBasePrice:         sdeGroup.UseBasePrice,
			Anchored:             sdeGroup.Anchored,
			Anchorable:           sdeGroup.Anchorable,
			FittableNonSingleton: sdeGroup.FittableNonSingleton,
			Published:            sdeGroup.Published,
		}
		result = append(result, invGroup)
	}

	// Sort by group ID for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].GroupID < result[j].GroupID
	})

	return result
}

// transformSystemJumps enriches system jumps with region and constellation IDs.
func (t *Transformer) transformSystemJumps(jumps []models.SystemJump, systems []models.SolarSystem) []models.SystemJump {
	// Build lookup map for system -> region/constellation
	systemLookup := make(map[int64]models.SolarSystem, len(systems))
	for _, sys := range systems {
		systemLookup[sys.SolarSystemID] = sys
	}

	result := make([]models.SystemJump, 0, len(jumps))
	for _, jump := range jumps {
		fromSys, fromOK := systemLookup[jump.FromSolarSystemID]
		toSys, toOK := systemLookup[jump.ToSolarSystemID]

		enrichedJump := models.SystemJump{
			FromSolarSystemID: jump.FromSolarSystemID,
			ToSolarSystemID:   jump.ToSolarSystemID,
		}

		if fromOK {
			enrichedJump.FromRegionID = fromSys.RegionID
			enrichedJump.FromConstellationID = fromSys.ConstellationID
		}

		if toOK {
			enrichedJump.ToRegionID = toSys.RegionID
			enrichedJump.ToConstellationID = toSys.ConstellationID
		}

		result = append(result, enrichedJump)
	}

	// Sort by from system ID, then to system ID for consistent output
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
		InvTypes:        len(data.InvTypes),
		InvGroups:       len(data.InvGroups),
		SystemJumps:     len(data.SystemJumps),
		WormholeClasses: len(data.WormholeClasses),
	}

	// Validation thresholds based on known EVE universe size
	const (
		minSolarSystems   = 8000
		minRegions        = 100
		minConstellations = 1000
		minTypes          = 30000 // All types, not just ships
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

	if result.InvTypes < minTypes {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Type count (%d) is below expected minimum (%d)",
				result.InvTypes, minTypes))
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
