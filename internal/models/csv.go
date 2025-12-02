package models

import (
	"fmt"
	"strconv"
)

// CSVHeaders defines the exact column headers for each CSV file to match Fuzzwork format.
var CSVHeaders = map[string][]string{
	"mapSolarSystems": {
		"regionID", "constellationID", "solarSystemID", "solarSystemName",
		"x", "y", "z", "xMin", "xMax", "yMin", "yMax", "zMin", "zMax",
		"luminosity", "border", "fringe", "corridor", "hub", "international",
		"regional", "constellation", "security", "factionID", "radius",
		"sunTypeID", "securityClass",
	},
	"mapRegions": {
		"regionID", "regionName", "x", "y", "z",
		"xMin", "xMax", "yMin", "yMax", "zMin", "zMax",
		"factionID", "nebula", "radius",
	},
	"mapConstellations": {
		"regionID", "constellationID", "constellationName",
		"x", "y", "z", "xMin", "xMax", "yMin", "yMax", "zMin", "zMax",
		"factionID", "radius",
	},
	"invTypes": {
		"typeID", "groupID", "typeName", "description",
		"mass", "volume", "capacity", "portionSize", "raceID",
		"basePrice", "published", "marketGroupID", "iconID", "soundID", "graphicID",
	},
	"invGroups": {
		"groupID", "categoryID", "groupName", "iconID",
		"useBasePrice", "anchored", "anchorable", "fittableNonSingleton", "published",
	},
	"mapLocationWormholeClasses": {
		"locationID", "wormholeClassID",
	},
	"mapSolarSystemJumps": {
		"fromRegionID", "fromConstellationID", "fromSolarSystemID",
		"toSolarSystemID", "toConstellationID", "toRegionID",
	},
}

// FormatNullableInt64 formats an optional int64 for CSV output.
// Returns "None" if nil, otherwise the integer value.
func FormatNullableInt64(v *int64) string {
	if v == nil {
		return "None"
	}
	return strconv.FormatInt(*v, 10)
}

// FormatBool formats a boolean for CSV output.
// Returns "1" for true, "0" for false (Fuzzwork format).
func FormatBool(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

// FormatFloat formats a float64 for CSV output.
// Uses full precision without scientific notation for coordinates.
func FormatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// FormatSecurity formats security status for CSV output.
// Matches Fuzzwork's precision for security values.
func FormatSecurity(v float64) string {
	return fmt.Sprintf("%.16g", v)
}

// Int64Ptr returns a pointer to an int64 value.
// Returns nil if the value is 0.
func Int64Ptr(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

// Int64PtrAlways returns a pointer to an int64 value, even if 0.
func Int64PtrAlways(v int64) *int64 {
	return &v
}

// ToCSVRow converts a SolarSystem to a CSV row matching Fuzzwork format.
func (s *SolarSystem) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(s.RegionID, 10),
		strconv.FormatInt(s.ConstellationID, 10),
		strconv.FormatInt(s.SolarSystemID, 10),
		s.SolarSystemName,
		FormatFloat(s.X),
		FormatFloat(s.Y),
		FormatFloat(s.Z),
		FormatFloat(s.XMin),
		FormatFloat(s.XMax),
		FormatFloat(s.YMin),
		FormatFloat(s.YMax),
		FormatFloat(s.ZMin),
		FormatFloat(s.ZMax),
		FormatFloat(s.Luminosity),
		FormatBool(s.Border),
		FormatBool(s.Fringe),
		FormatBool(s.Corridor),
		FormatBool(s.Hub),
		FormatBool(s.International),
		FormatBool(s.Regional),
		s.Constellation, // Always "None"
		FormatSecurity(s.Security),
		FormatNullableInt64(s.FactionID),
		FormatFloat(s.Radius),
		FormatNullableInt64(s.SunTypeID),
		s.SecurityClass,
	}
}

// ToCSVRow converts a Region to a CSV row matching Fuzzwork format.
func (r *Region) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(r.RegionID, 10),
		r.RegionName,
		FormatFloat(r.X),
		FormatFloat(r.Y),
		FormatFloat(r.Z),
		FormatFloat(r.XMin),
		FormatFloat(r.XMax),
		FormatFloat(r.YMin),
		FormatFloat(r.YMax),
		FormatFloat(r.ZMin),
		FormatFloat(r.ZMax),
		FormatNullableInt64(r.FactionID),
		strconv.FormatInt(r.Nebula, 10),
		FormatFloat(r.Radius),
	}
}

// ToCSVRow converts a Constellation to a CSV row matching Fuzzwork format.
func (c *Constellation) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(c.RegionID, 10),
		strconv.FormatInt(c.ConstellationID, 10),
		c.ConstellationName,
		FormatFloat(c.X),
		FormatFloat(c.Y),
		FormatFloat(c.Z),
		FormatFloat(c.XMin),
		FormatFloat(c.XMax),
		FormatFloat(c.YMin),
		FormatFloat(c.YMax),
		FormatFloat(c.ZMin),
		FormatFloat(c.ZMax),
		FormatNullableInt64(c.FactionID),
		FormatFloat(c.Radius),
	}
}

// ToCSVRow converts an InvType to a CSV row matching Fuzzwork format.
func (t *InvType) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(t.TypeID, 10),
		strconv.FormatInt(t.GroupID, 10),
		t.TypeName,
		t.Description,
		FormatFloat(t.Mass),
		FormatFloat(t.Volume),
		FormatFloat(t.Capacity),
		strconv.FormatInt(t.PortionSize, 10),
		FormatNullableInt64(t.RaceID),
		FormatFloat(t.BasePrice),
		FormatBool(t.Published),
		FormatNullableInt64(t.MarketGroupID),
		FormatNullableInt64(t.IconID),
		FormatNullableInt64(t.SoundID),
		FormatNullableInt64(t.GraphicID),
	}
}

// ToCSVRow converts an InvGroup to a CSV row matching Fuzzwork format.
func (g *InvGroup) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(g.GroupID, 10),
		strconv.FormatInt(g.CategoryID, 10),
		g.GroupName,
		FormatNullableInt64(g.IconID),
		FormatBool(g.UseBasePrice),
		FormatBool(g.Anchored),
		FormatBool(g.Anchorable),
		FormatBool(g.FittableNonSingleton),
		FormatBool(g.Published),
	}
}

// ToCSVRow converts a WormholeClassLocation to a CSV row matching Fuzzwork format.
func (w *WormholeClassLocation) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(w.LocationID, 10),
		strconv.FormatInt(w.WormholeClassID, 10),
	}
}

// ToCSVRow converts a SystemJump to a CSV row matching Fuzzwork format.
func (j *SystemJump) ToCSVRow() []string {
	return []string{
		strconv.FormatInt(j.FromRegionID, 10),
		strconv.FormatInt(j.FromConstellationID, 10),
		strconv.FormatInt(j.FromSolarSystemID, 10),
		strconv.FormatInt(j.ToSolarSystemID, 10),
		strconv.FormatInt(j.ToConstellationID, 10),
		strconv.FormatInt(j.ToRegionID, 10),
	}
}
