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
}

// SDEMapConstellation represents a constellation in the flat SDE format.
type SDEMapConstellation struct {
	ConstellationID int64             `yaml:"constellationID"`
	RegionID        int64             `yaml:"regionID"`
	Name            map[string]string `yaml:"name"`
	NameID          int64             `yaml:"nameID,omitempty"`
	FactionID       int64             `yaml:"factionID,omitempty"`
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

		regions = append(regions, models.Region{
			RegionID:   id,
			RegionName: name,
		})
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

		constellations = append(constellations, models.Constellation{
			ConstellationID:   id,
			ConstellationName: name,
			RegionID:          data.RegionID,
		})
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

		systems = append(systems, models.SolarSystem{
			SolarSystemID:   id,
			RegionID:        data.RegionID,
			ConstellationID: data.ConstellationID,
			SolarSystemName: name,
			SunTypeID:       data.SunTypeID,
			Security:        data.Security,
		})
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
