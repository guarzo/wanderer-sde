package models

// SolarSystem represents a solar system in Wanderer's format.
type SolarSystem struct {
	SolarSystemID   int64   `json:"solarSystemID"`
	RegionID        int64   `json:"regionID"`
	ConstellationID int64   `json:"constellationID"`
	SolarSystemName string  `json:"solarSystemName"`
	SunTypeID       int64   `json:"sunTypeID,omitempty"`
	Security        float64 `json:"security"`
}

// Region represents a region in Wanderer's format.
type Region struct {
	RegionID   int64  `json:"regionID"`
	RegionName string `json:"regionName"`
}

// Constellation represents a constellation in Wanderer's format.
type Constellation struct {
	ConstellationID   int64  `json:"constellationID"`
	ConstellationName string `json:"constellationName"`
	RegionID          int64  `json:"regionID"`
}

// WormholeClassLocation represents a wormhole class assignment in Wanderer's format.
type WormholeClassLocation struct {
	LocationID      int64 `json:"locationID"`
	WormholeClassID int64 `json:"wormholeClassID"`
}

// ShipType represents a ship type in Wanderer's format.
type ShipType struct {
	TypeID   int64   `json:"typeID"`
	GroupID  int64   `json:"groupID"`
	TypeName string  `json:"typeName"`
	Mass     float64 `json:"mass,omitempty"`
	Volume   float64 `json:"volume,omitempty"`
	Capacity float64 `json:"capacity,omitempty"`
}

// ItemGroup represents an item group in Wanderer's format.
type ItemGroup struct {
	GroupID    int64  `json:"groupID"`
	CategoryID int64  `json:"categoryID"`
	GroupName  string `json:"groupName"`
}

// SystemJump represents a stargate connection in Wanderer's format.
type SystemJump struct {
	FromSolarSystemID int64 `json:"fromSolarSystemID"`
	ToSolarSystemID   int64 `json:"toSolarSystemID"`
}

// UniverseData holds all parsed universe data.
type UniverseData struct {
	Regions        []Region
	Constellations []Constellation
	SolarSystems   []SolarSystem
}

// ConvertedData holds all data ready for output.
type ConvertedData struct {
	Universe        *UniverseData
	ShipTypes       []ShipType
	ItemGroups      []ItemGroup
	WormholeClasses []WormholeClassLocation
	SystemJumps     []SystemJump
}

// ValidationResult holds the results of data validation.
type ValidationResult struct {
	SolarSystems    int
	Regions         int
	Constellations  int
	ShipTypes       int
	ItemGroups      int
	SystemJumps     int
	WormholeClasses int
	Errors          []string
	Warnings        []string
}

// IsValid returns true if validation found no errors.
func (v *ValidationResult) IsValid() bool {
	return len(v.Errors) == 0
}
