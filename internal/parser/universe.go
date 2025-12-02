package parser

import (
	"fmt"
	"sort"

	"github.com/guarzo/wanderer-sde/internal/models"
	"github.com/guarzo/wanderer-sde/pkg/yaml"
)

// SDEMapRegion represents a region in the flat SDE format.
type SDEMapRegion struct {
	RegionID      int64             `yaml:"regionID"`
	Name          map[string]string `yaml:"name"`
	NameID        int64             `yaml:"nameID,omitempty"`
	DescriptionID int64             `yaml:"descriptionID,omitempty"`
	FactionID     int64             `yaml:"factionID,omitempty"`
	Center        []float64         `yaml:"center,omitempty"`
	Max           []float64         `yaml:"max,omitempty"`
	Min           []float64         `yaml:"min,omitempty"`
}

// SDEMapConstellation represents a constellation in the flat SDE format.
type SDEMapConstellation struct {
	ConstellationID int64             `yaml:"constellationID"`
	RegionID        int64             `yaml:"regionID"`
	Name            map[string]string `yaml:"name"`
	NameID          int64             `yaml:"nameID,omitempty"`
	FactionID       int64             `yaml:"factionID,omitempty"`
	Center          []float64         `yaml:"center,omitempty"`
	Max             []float64         `yaml:"max,omitempty"`
	Min             []float64         `yaml:"min,omitempty"`
	Radius          float64           `yaml:"radius,omitempty"`
}

// SDEMapSolarSystem represents a solar system in the flat SDE format.
type SDEMapSolarSystem struct {
	SolarSystemID   int64             `yaml:"solarSystemID"`
	ConstellationID int64             `yaml:"constellationID"`
	RegionID        int64             `yaml:"regionID"`
	Name            map[string]string `yaml:"name"`
	Security        float64           `yaml:"security"`
	SecurityClass   string            `yaml:"securityClass,omitempty"`
	SunTypeID       int64             `yaml:"sunTypeID,omitempty"`
	WormholeClassID int64             `yaml:"wormholeClassID,omitempty"`
	FactionID       int64             `yaml:"factionID,omitempty"`
	Border          bool              `yaml:"border,omitempty"`
	Corridor        bool              `yaml:"corridor,omitempty"`
	Fringe          bool              `yaml:"fringe,omitempty"`
	Hub             bool              `yaml:"hub,omitempty"`
	International   bool              `yaml:"international,omitempty"`
	Regional        bool              `yaml:"regional,omitempty"`
	Center          []float64         `yaml:"center,omitempty"`
	Max             []float64         `yaml:"max,omitempty"`
	Min             []float64         `yaml:"min,omitempty"`
	Luminosity      float64           `yaml:"luminosity,omitempty"`
	Radius          float64           `yaml:"radius,omitempty"`
}

// ParseRegions parses the mapRegions.yaml file.
func (p *Parser) ParseRegions() ([]models.Region, error) {
	path := p.filePath("mapRegions.yaml")

	// Parse the file as a map of region ID to region data
	rawRegions, err := yaml.ParseFileMap[int64, SDEMapRegion](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regions file: %w", err)
	}

	regions := make([]models.Region, 0, len(rawRegions))
	for id, data := range rawRegions {
		name := data.Name["en"]
		if name == "" {
			// Fall back to using ID-based name if no English name
			name = fmt.Sprintf("Region %d", id)
		}

		region := models.Region{
			RegionID:   id,
			RegionName: name,
			FactionID:  models.Int64Ptr(data.FactionID),
			Nebula:     0, // Not available in SDE
			Radius:     0, // Not directly available, could be calculated from min/max
		}

		// Extract coordinates from center array
		if len(data.Center) >= 3 {
			region.X = data.Center[0]
			region.Y = data.Center[1]
			region.Z = data.Center[2]
		}

		// Extract min coordinates
		if len(data.Min) >= 3 {
			region.XMin = data.Min[0]
			region.YMin = data.Min[1]
			region.ZMin = data.Min[2]
		}

		// Extract max coordinates
		if len(data.Max) >= 3 {
			region.XMax = data.Max[0]
			region.YMax = data.Max[1]
			region.ZMax = data.Max[2]
		}

		regions = append(regions, region)
	}

	// Sort by region ID for consistent output
	sort.Slice(regions, func(i, j int) bool {
		return regions[i].RegionID < regions[j].RegionID
	})

	return regions, nil
}

// ParseConstellations parses the mapConstellations.yaml file.
func (p *Parser) ParseConstellations() ([]models.Constellation, error) {
	path := p.filePath("mapConstellations.yaml")

	// Parse the file as a map of constellation ID to constellation data
	rawConstellations, err := yaml.ParseFileMap[int64, SDEMapConstellation](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse constellations file: %w", err)
	}

	constellations := make([]models.Constellation, 0, len(rawConstellations))
	for id, data := range rawConstellations {
		name := data.Name["en"]
		if name == "" {
			name = fmt.Sprintf("Constellation %d", id)
		}

		constellation := models.Constellation{
			RegionID:          data.RegionID,
			ConstellationID:   id,
			ConstellationName: name,
			FactionID:         models.Int64Ptr(data.FactionID),
			Radius:            data.Radius,
		}

		// Extract coordinates from center array
		if len(data.Center) >= 3 {
			constellation.X = data.Center[0]
			constellation.Y = data.Center[1]
			constellation.Z = data.Center[2]
		}

		// Extract min coordinates
		if len(data.Min) >= 3 {
			constellation.XMin = data.Min[0]
			constellation.YMin = data.Min[1]
			constellation.ZMin = data.Min[2]
		}

		// Extract max coordinates
		if len(data.Max) >= 3 {
			constellation.XMax = data.Max[0]
			constellation.YMax = data.Max[1]
			constellation.ZMax = data.Max[2]
		}

		constellations = append(constellations, constellation)
	}

	// Sort by constellation ID for consistent output
	sort.Slice(constellations, func(i, j int) bool {
		return constellations[i].ConstellationID < constellations[j].ConstellationID
	})

	return constellations, nil
}

// ParseSolarSystems parses the mapSolarSystems.yaml file.
func (p *Parser) ParseSolarSystems() ([]models.SolarSystem, error) {
	path := p.filePath("mapSolarSystems.yaml")

	// Parse the file as a map of system ID to system data
	rawSystems, err := yaml.ParseFileMap[int64, SDEMapSolarSystem](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse solar systems file: %w", err)
	}

	systems := make([]models.SolarSystem, 0, len(rawSystems))
	for id, data := range rawSystems {
		name := data.Name["en"]
		if name == "" {
			name = fmt.Sprintf("System %d", id)
		}

		system := models.SolarSystem{
			RegionID:        data.RegionID,
			ConstellationID: data.ConstellationID,
			SolarSystemID:   id,
			SolarSystemName: name,
			Luminosity:      data.Luminosity,
			Border:          data.Border,
			Fringe:          data.Fringe,
			Corridor:        data.Corridor,
			Hub:             data.Hub,
			International:   data.International,
			Regional:        data.Regional,
			Constellation:   "None", // Always "None" - legacy field
			Security:        data.Security,
			FactionID:       models.Int64Ptr(data.FactionID),
			Radius:          data.Radius,
			SunTypeID:       models.Int64Ptr(data.SunTypeID),
			SecurityClass:   data.SecurityClass,
		}

		// Extract coordinates from center array
		if len(data.Center) >= 3 {
			system.X = data.Center[0]
			system.Y = data.Center[1]
			system.Z = data.Center[2]
		}

		// Extract min coordinates
		if len(data.Min) >= 3 {
			system.XMin = data.Min[0]
			system.YMin = data.Min[1]
			system.ZMin = data.Min[2]
		}

		// Extract max coordinates
		if len(data.Max) >= 3 {
			system.XMax = data.Max[0]
			system.YMax = data.Max[1]
			system.ZMax = data.Max[2]
		}

		systems = append(systems, system)
	}

	// Sort by solar system ID for consistent output
	sort.Slice(systems, func(i, j int) bool {
		return systems[i].SolarSystemID < systems[j].SolarSystemID
	})

	return systems, nil
}

// ExtractWormholeClasses extracts wormhole class information from solar systems.
// In the new SDE format, wormhole class is stored directly on the solar system.
func (p *Parser) ExtractWormholeClasses(systems []models.SolarSystem) []models.WormholeClassLocation {
	// Re-parse systems to get wormhole class info (not stored in Wanderer model)
	path := p.filePath("mapSolarSystems.yaml")
	rawSystems, err := yaml.ParseFileMap[int64, SDEMapSolarSystem](path)
	if err != nil {
		// If we can't parse, return empty slice
		return nil
	}

	var wormholeClasses []models.WormholeClassLocation
	for id, data := range rawSystems {
		if data.WormholeClassID != 0 {
			wormholeClasses = append(wormholeClasses, models.WormholeClassLocation{
				LocationID:      id,
				WormholeClassID: data.WormholeClassID,
			})
		}
	}

	// Sort by location ID for consistent output
	sort.Slice(wormholeClasses, func(i, j int) bool {
		return wormholeClasses[i].LocationID < wormholeClasses[j].LocationID
	})

	return wormholeClasses
}
