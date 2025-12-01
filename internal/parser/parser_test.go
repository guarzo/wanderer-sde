package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
)

// createTestSDE creates a minimal SDE structure for testing
func createTestSDE(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "parser_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create mapRegions.yaml
	regionsYAML := `10000002:
  regionID: 10000002
  name:
    en: "The Forge"
10000001:
  regionID: 10000001
  name:
    en: "Derelik"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(regionsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapRegions.yaml: %v", err)
	}

	// Create mapConstellations.yaml
	constellationsYAML := `20000020:
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "Kimotoro"
20000001:
  constellationID: 20000001
  regionID: 10000001
  name:
    en: "Joas"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapConstellations.yaml"), []byte(constellationsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapConstellations.yaml: %v", err)
	}

	// Create mapSolarSystems.yaml
	systemsYAML := `30000142:
  solarSystemID: 30000142
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "Jita"
  security: 0.9459
  sunTypeID: 6
  wormholeClassID: 0
30000001:
  solarSystemID: 30000001
  constellationID: 20000001
  regionID: 10000001
  name:
    en: "Tanoo"
  security: 0.8576
  sunTypeID: 7
  wormholeClassID: 0
31000001:
  solarSystemID: 31000001
  constellationID: 21000001
  regionID: 11000001
  name:
    en: "J123456"
  security: -1.0
  sunTypeID: 45041
  wormholeClassID: 3
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapSolarSystems.yaml"), []byte(systemsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapSolarSystems.yaml: %v", err)
	}

	// Create mapStargates.yaml
	stargatesYAML := `50000001:
  solarSystemID: 30000142
  destination:
    solarSystemID: 30000001
    stargateID: 50000002
  typeID: 16
50000002:
  solarSystemID: 30000001
  destination:
    solarSystemID: 30000142
    stargateID: 50000001
  typeID: 16
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapStargates.yaml"), []byte(stargatesYAML), 0644); err != nil {
		t.Fatalf("failed to create mapStargates.yaml: %v", err)
	}

	// Create types.yaml
	typesYAML := `587:
  groupID: 25
  name:
    en: "Rifter"
  mass: 1350000.0
  volume: 27500.0
  capacity: 125.0
  published: true
588:
  groupID: 25
  name:
    en: "Slasher"
  mass: 1200000.0
  volume: 26000.0
  capacity: 115.0
  published: true
2456:
  groupID: 18
  name:
    en: "Hobgoblin I"
  mass: 2500.0
  volume: 5.0
  published: true
`
	if err := os.WriteFile(filepath.Join(tmpDir, "types.yaml"), []byte(typesYAML), 0644); err != nil {
		t.Fatalf("failed to create types.yaml: %v", err)
	}

	// Create groups.yaml
	groupsYAML := `25:
  categoryID: 6
  name:
    en: "Frigate"
  published: true
18:
  categoryID: 7
  name:
    en: "Drone"
  published: true
`
	if err := os.WriteFile(filepath.Join(tmpDir, "groups.yaml"), []byte(groupsYAML), 0644); err != nil {
		t.Fatalf("failed to create groups.yaml: %v", err)
	}

	// Create categories.yaml
	categoriesYAML := `6:
  name:
    en: "Ship"
  published: true
7:
  name:
    en: "Drone"
  published: true
`
	if err := os.WriteFile(filepath.Join(tmpDir, "categories.yaml"), []byte(categoriesYAML), 0644); err != nil {
		t.Fatalf("failed to create categories.yaml: %v", err)
	}

	return tmpDir
}

func TestParser_ParseRegions(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	regions, err := p.ParseRegions()
	if err != nil {
		t.Fatalf("ParseRegions failed: %v", err)
	}

	if len(regions) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(regions))
	}

	// Verify regions are sorted by ID
	for i := 1; i < len(regions); i++ {
		if regions[i-1].RegionID >= regions[i].RegionID {
			t.Errorf("Regions not sorted: %d >= %d", regions[i-1].RegionID, regions[i].RegionID)
		}
	}

	// Verify first region (should be Derelik with lowest ID)
	if regions[0].RegionID != 10000001 {
		t.Errorf("Expected first region ID to be 10000001, got %d", regions[0].RegionID)
	}
	if regions[0].RegionName != "Derelik" {
		t.Errorf("Expected first region name to be 'Derelik', got %q", regions[0].RegionName)
	}
}

func TestParser_ParseConstellations(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	constellations, err := p.ParseConstellations()
	if err != nil {
		t.Fatalf("ParseConstellations failed: %v", err)
	}

	if len(constellations) != 2 {
		t.Errorf("Expected 2 constellations, got %d", len(constellations))
	}

	// Verify constellations are sorted by ID
	for i := 1; i < len(constellations); i++ {
		if constellations[i-1].ConstellationID >= constellations[i].ConstellationID {
			t.Errorf("Constellations not sorted: %d >= %d", constellations[i-1].ConstellationID, constellations[i].ConstellationID)
		}
	}

	// Verify constellation has correct region ID reference
	for _, c := range constellations {
		if c.ConstellationID == 20000020 && c.RegionID != 10000002 {
			t.Errorf("Constellation 20000020 should have RegionID 10000002, got %d", c.RegionID)
		}
	}
}

func TestParser_ParseSolarSystems(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	systems, err := p.ParseSolarSystems()
	if err != nil {
		t.Fatalf("ParseSolarSystems failed: %v", err)
	}

	if len(systems) != 3 {
		t.Errorf("Expected 3 solar systems, got %d", len(systems))
	}

	// Verify systems are sorted by ID
	for i := 1; i < len(systems); i++ {
		if systems[i-1].SolarSystemID >= systems[i].SolarSystemID {
			t.Errorf("Solar systems not sorted: %d >= %d", systems[i-1].SolarSystemID, systems[i].SolarSystemID)
		}
	}

	// Find Jita and verify its data
	var jita *struct {
		id       int64
		security float64
		sunType  int64
	}
	for _, s := range systems {
		if s.SolarSystemName == "Jita" {
			jita = &struct {
				id       int64
				security float64
				sunType  int64
			}{s.SolarSystemID, s.Security, s.SunTypeID}
			break
		}
	}

	if jita == nil {
		t.Fatal("Jita not found in parsed systems")
	}

	if jita.id != 30000142 {
		t.Errorf("Expected Jita ID to be 30000142, got %d", jita.id)
	}

	// Note: raw security value is 0.9459, transformation happens in transformer
	if jita.security < 0.9 || jita.security > 1.0 {
		t.Errorf("Jita security value unexpected: %f", jita.security)
	}
}

func TestParser_ParseStargates(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	jumps, err := p.ParseStargates()
	if err != nil {
		t.Fatalf("ParseStargates failed: %v", err)
	}

	// Two stargates in both directions should result in 1 unique jump
	if len(jumps) != 1 {
		t.Errorf("Expected 1 unique jump (deduped from bidirectional stargates), got %d", len(jumps))
	}

	// Verify the jump (smaller ID should be first)
	if len(jumps) > 0 {
		jump := jumps[0]
		if jump.FromSolarSystemID != 30000001 || jump.ToSolarSystemID != 30000142 {
			t.Errorf("Expected jump from 30000001 to 30000142, got %d to %d",
				jump.FromSolarSystemID, jump.ToSolarSystemID)
		}
	}
}

func TestParser_ParseTypes(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	types, err := p.ParseTypes()
	if err != nil {
		t.Fatalf("ParseTypes failed: %v", err)
	}

	if len(types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(types))
	}

	// Verify Rifter data
	rifter, ok := types[587]
	if !ok {
		t.Fatal("Rifter (587) not found in types")
	}

	if rifter.GroupID != 25 {
		t.Errorf("Expected Rifter groupID to be 25, got %d", rifter.GroupID)
	}

	if rifter.Name["en"] != "Rifter" {
		t.Errorf("Expected Rifter name to be 'Rifter', got %q", rifter.Name["en"])
	}

	if rifter.Mass != 1350000.0 {
		t.Errorf("Expected Rifter mass to be 1350000.0, got %f", rifter.Mass)
	}
}

func TestParser_ParseGroups(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	groups, err := p.ParseGroups()
	if err != nil {
		t.Fatalf("ParseGroups failed: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	// Verify Frigate group
	frigate, ok := groups[25]
	if !ok {
		t.Fatal("Frigate group (25) not found")
	}

	if frigate.CategoryID != 6 {
		t.Errorf("Expected Frigate categoryID to be 6, got %d", frigate.CategoryID)
	}

	if frigate.Name["en"] != "Frigate" {
		t.Errorf("Expected Frigate name, got %q", frigate.Name["en"])
	}
}

func TestParser_ParseCategories(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	categories, err := p.ParseCategories()
	if err != nil {
		t.Fatalf("ParseCategories failed: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	// Verify Ship category
	ship, ok := categories[6]
	if !ok {
		t.Fatal("Ship category (6) not found")
	}

	if ship.Name["en"] != "Ship" {
		t.Errorf("Expected Ship name, got %q", ship.Name["en"])
	}
}

func TestParser_ExtractWormholeClasses(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	// First parse solar systems
	systems, err := p.ParseSolarSystems()
	if err != nil {
		t.Fatalf("ParseSolarSystems failed: %v", err)
	}

	wormholeClasses := p.ExtractWormholeClasses(systems)

	// Only the wormhole system (J123456) should have a non-zero wormhole class
	expectedWHCount := 1
	if len(wormholeClasses) != expectedWHCount {
		t.Errorf("Expected %d wormhole class entries, got %d", expectedWHCount, len(wormholeClasses))
	}

	// Verify the wormhole class entry
	if len(wormholeClasses) > 0 {
		wh := wormholeClasses[0]
		if wh.LocationID != 31000001 {
			t.Errorf("Expected wormhole location ID 31000001, got %d", wh.LocationID)
		}
		if wh.WormholeClassID != 3 {
			t.Errorf("Expected wormhole class ID 3, got %d", wh.WormholeClassID)
		}
	}
}

func TestParser_ParseAll(t *testing.T) {
	tmpDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	result, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	// Verify all data was parsed
	if len(result.Regions) == 0 {
		t.Error("No regions parsed")
	}
	if len(result.Constellations) == 0 {
		t.Error("No constellations parsed")
	}
	if len(result.SolarSystems) == 0 {
		t.Error("No solar systems parsed")
	}
	if len(result.Types) == 0 {
		t.Error("No types parsed")
	}
	if len(result.Groups) == 0 {
		t.Error("No groups parsed")
	}
	if len(result.Categories) == 0 {
		t.Error("No categories parsed")
	}
	if len(result.SystemJumps) == 0 {
		t.Error("No system jumps parsed")
	}
}

func TestParser_MissingFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parser_test_empty")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	// Try to parse non-existent file
	_, err = p.ParseRegions()
	if err == nil {
		t.Error("Expected error when parsing missing file")
	}
}

func TestParser_MalformedYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parser_test_malformed")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create malformed YAML
	malformedYAML := `this is not valid yaml: [[[`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(malformedYAML), 0644); err != nil {
		t.Fatalf("failed to create malformed yaml: %v", err)
	}

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	_, err = p.ParseRegions()
	if err == nil {
		t.Error("Expected error when parsing malformed YAML")
	}
}

func TestParser_EmptyName(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parser_test_empty_name")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create region with empty name
	regionsYAML := `10000001:
  regionID: 10000001
  name: {}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapRegions.yaml"), []byte(regionsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapRegions.yaml: %v", err)
	}

	cfg := &config.Config{Verbose: false}
	p := New(cfg, tmpDir)

	regions, err := p.ParseRegions()
	if err != nil {
		t.Fatalf("ParseRegions failed: %v", err)
	}

	// Should have a fallback name
	if len(regions) == 0 {
		t.Fatal("No regions parsed")
	}

	if regions[0].RegionName == "" {
		t.Error("Expected fallback name for region with empty name")
	}
}
