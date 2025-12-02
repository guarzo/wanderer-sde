package transformer

import "github.com/guarzo/wanderer-sde/internal/models"

// CalculateRegionBounds calculates min/max coordinates for regions
// based on their constituent solar systems.
func CalculateRegionBounds(regions []models.Region, systems []models.SolarSystem) {
	// Build map of regionID -> systems
	regionSystems := make(map[int64][]models.SolarSystem)
	for _, sys := range systems {
		regionSystems[sys.RegionID] = append(regionSystems[sys.RegionID], sys)
	}

	// Calculate bounds for each region
	for i := range regions {
		sysList := regionSystems[regions[i].RegionID]
		if len(sysList) == 0 {
			continue
		}

		minX, maxX := sysList[0].X, sysList[0].X
		minY, maxY := sysList[0].Y, sysList[0].Y
		minZ, maxZ := sysList[0].Z, sysList[0].Z

		for _, sys := range sysList {
			if sys.X < minX {
				minX = sys.X
			}
			if sys.X > maxX {
				maxX = sys.X
			}
			if sys.Y < minY {
				minY = sys.Y
			}
			if sys.Y > maxY {
				maxY = sys.Y
			}
			if sys.Z < minZ {
				minZ = sys.Z
			}
			if sys.Z > maxZ {
				maxZ = sys.Z
			}
		}

		regions[i].XMin = minX
		regions[i].XMax = maxX
		regions[i].YMin = minY
		regions[i].YMax = maxY
		regions[i].ZMin = minZ
		regions[i].ZMax = maxZ
	}
}

// CalculateConstellationBounds calculates min/max coordinates for constellations
// based on their constituent solar systems.
func CalculateConstellationBounds(constellations []models.Constellation, systems []models.SolarSystem) {
	// Build map of constellationID -> systems
	constellationSystems := make(map[int64][]models.SolarSystem)
	for _, sys := range systems {
		constellationSystems[sys.ConstellationID] = append(constellationSystems[sys.ConstellationID], sys)
	}

	// Calculate bounds for each constellation
	for i := range constellations {
		sysList := constellationSystems[constellations[i].ConstellationID]
		if len(sysList) == 0 {
			continue
		}

		minX, maxX := sysList[0].X, sysList[0].X
		minY, maxY := sysList[0].Y, sysList[0].Y
		minZ, maxZ := sysList[0].Z, sysList[0].Z

		for _, sys := range sysList {
			if sys.X < minX {
				minX = sys.X
			}
			if sys.X > maxX {
				maxX = sys.X
			}
			if sys.Y < minY {
				minY = sys.Y
			}
			if sys.Y > maxY {
				maxY = sys.Y
			}
			if sys.Z < minZ {
				minZ = sys.Z
			}
			if sys.Z > maxZ {
				maxZ = sys.Z
			}
		}

		constellations[i].XMin = minX
		constellations[i].XMax = maxX
		constellations[i].YMin = minY
		constellations[i].YMax = maxY
		constellations[i].ZMin = minZ
		constellations[i].ZMax = maxZ
	}
}

// InheritFactionIDs sets factionID on solar systems that don't have one,
// inheriting from their region.
func InheritFactionIDs(systems []models.SolarSystem, regions []models.Region) {
	// Build map of regionID -> factionID
	regionFactions := make(map[int64]*int64)
	for _, region := range regions {
		regionFactions[region.RegionID] = region.FactionID
	}

	// Inherit factionID from region if system doesn't have one
	for i := range systems {
		if systems[i].FactionID == nil {
			if factionID, ok := regionFactions[systems[i].RegionID]; ok && factionID != nil {
				// Copy the value to avoid pointer aliasing
				val := *factionID
				systems[i].FactionID = &val
			}
		}
	}
}
