package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/parser"
	"github.com/guarzo/wanderer-sde/internal/transformer"
	"github.com/guarzo/wanderer-sde/internal/writer"
)

// createTestSDE creates a complete minimal SDE structure for integration testing
func createTestSDE(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create mapRegions.yaml with multiple regions
	regionsYAML := `10000002:
  regionID: 10000002
  name:
    en: "The Forge"
10000001:
  regionID: 10000001
  name:
    en: "Derelik"
11000001:
  regionID: 11000001
  name:
    en: "J-Space Region"
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
21000001:
  constellationID: 21000001
  regionID: 11000001
  name:
    en: "J-Constellation"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapConstellations.yaml"), []byte(constellationsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapConstellations.yaml: %v", err)
	}

	// Create mapSolarSystems.yaml with various security levels
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
30000144:
  solarSystemID: 30000144
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "Perimeter"
  security: 0.94
  sunTypeID: 8
  wormholeClassID: 0
30000145:
  solarSystemID: 30000145
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "LowSec System"
  security: 0.4
  sunTypeID: 9
  wormholeClassID: 0
30000146:
  solarSystemID: 30000146
  constellationID: 20000020
  regionID: 10000002
  name:
    en: "NullSec System"
  security: -0.7
  sunTypeID: 10
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
31000002:
  solarSystemID: 31000002
  constellationID: 21000001
  regionID: 11000001
  name:
    en: "J654321"
  security: -1.0
  sunTypeID: 45041
  wormholeClassID: 5
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapSolarSystems.yaml"), []byte(systemsYAML), 0644); err != nil {
		t.Fatalf("failed to create mapSolarSystems.yaml: %v", err)
	}

	// Create mapStargates.yaml
	stargatesYAML := `50000001:
  solarSystemID: 30000142
  destination:
    solarSystemID: 30000144
    stargateID: 50000002
  typeID: 16
50000002:
  solarSystemID: 30000144
  destination:
    solarSystemID: 30000142
    stargateID: 50000001
  typeID: 16
50000003:
  solarSystemID: 30000144
  destination:
    solarSystemID: 30000145
    stargateID: 50000004
  typeID: 16
50000004:
  solarSystemID: 30000145
  destination:
    solarSystemID: 30000144
    stargateID: 50000003
  typeID: 16
50000005:
  solarSystemID: 30000145
  destination:
    solarSystemID: 30000146
    stargateID: 50000006
  typeID: 16
50000006:
  solarSystemID: 30000146
  destination:
    solarSystemID: 30000145
    stargateID: 50000005
  typeID: 16
`
	if err := os.WriteFile(filepath.Join(tmpDir, "mapStargates.yaml"), []byte(stargatesYAML), 0644); err != nil {
		t.Fatalf("failed to create mapStargates.yaml: %v", err)
	}

	// Create types.yaml with ships and non-ships
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
589:
  groupID: 25
  name:
    en: "Breacher"
  mass: 1100000.0
  volume: 25000.0
  capacity: 130.0
  published: true
625:
  groupID: 26
  name:
    en: "Caracal"
  mass: 11000000.0
  volume: 92000.0
  capacity: 350.0
  published: true
2456:
  groupID: 18
  name:
    en: "Hobgoblin I"
  mass: 2500.0
  volume: 5.0
  published: true
17738:
  groupID: 419
  name:
    en: "Prophecy"
  mass: 15500000.0
  volume: 270000.0
  capacity: 500.0
  published: true
99999:
  groupID: 25
  name:
    en: "Unpublished Ship"
  mass: 1000000.0
  volume: 20000.0
  published: false
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
26:
  categoryID: 6
  name:
    en: "Cruiser"
  published: true
18:
  categoryID: 7
  name:
    en: "Drone"
  published: true
419:
  categoryID: 6
  name:
    en: "Battlecruiser"
  published: true
999:
  categoryID: 6
  name:
    en: "Unpublished Group"
  published: false
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

func TestIntegration_FullPipeline(t *testing.T) {
	// Create test SDE
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	// Create output directory
	outputDir, err := os.MkdirTemp("", "integration_output")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	cfg := &config.Config{
		SDEPath:     sdeDir,
		OutputDir:   outputDir,
		Verbose:     false,
		PrettyPrint: true,
	}

	// Step 1: Parse
	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	// Verify parsing results
	t.Run("parsing", func(t *testing.T) {
		if len(parseResult.Regions) != 3 {
			t.Errorf("Expected 3 regions, got %d", len(parseResult.Regions))
		}
		if len(parseResult.Constellations) != 3 {
			t.Errorf("Expected 3 constellations, got %d", len(parseResult.Constellations))
		}
		if len(parseResult.SolarSystems) != 7 {
			t.Errorf("Expected 7 solar systems, got %d", len(parseResult.SolarSystems))
		}
		if len(parseResult.Types) != 7 {
			t.Errorf("Expected 7 types, got %d", len(parseResult.Types))
		}
		if len(parseResult.Groups) != 5 {
			t.Errorf("Expected 5 groups, got %d", len(parseResult.Groups))
		}
		if len(parseResult.Categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(parseResult.Categories))
		}
		// 3 unique stargate connections (bidirectional pairs)
		if len(parseResult.SystemJumps) != 3 {
			t.Errorf("Expected 3 system jumps, got %d", len(parseResult.SystemJumps))
		}
	})

	// Step 2: Transform
	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify transformation results
	t.Run("transformation", func(t *testing.T) {
		// Should filter to only ship types (6 ships - all types in ship groups, regardless of published status)
		// Types in ship category groups: 587, 588, 589 (Frigate), 625 (Cruiser), 17738 (Battlecruiser), 99999 (Unpublished but in ship group)
		if len(convertedData.ShipTypes) != 6 {
			t.Errorf("Expected 6 ship types, got %d", len(convertedData.ShipTypes))
		}

		// Should filter to only ship groups (4 groups with categoryID 6, including unpublished)
		if len(convertedData.ItemGroups) != 4 {
			t.Errorf("Expected 4 ship groups, got %d", len(convertedData.ItemGroups))
		}

		// Should extract wormhole classes (2 wormhole systems)
		if len(convertedData.WormholeClasses) != 2 {
			t.Errorf("Expected 2 wormhole classes, got %d", len(convertedData.WormholeClasses))
		}

		// Verify security calculation for specific systems
		for _, sys := range convertedData.Universe.SolarSystems {
			switch sys.SolarSystemName {
			case "Jita":
				// 0.9459 should round to 0.9
				if sys.Security != 0.9 {
					t.Errorf("Jita security should be 0.9, got %f", sys.Security)
				}
			case "J123456":
				// Wormhole security stays at -1.0
				if sys.Security != -1.0 {
					t.Errorf("J123456 security should be -1.0, got %f", sys.Security)
				}
			}
		}
	})

	// Step 3: Validate
	t.Run("validation", func(t *testing.T) {
		validationResult := tr.Validate(convertedData)

		// Should have warnings about low counts (test data is minimal)
		if len(validationResult.Warnings) == 0 {
			t.Log("No warnings generated (expected for test data)")
		}

		// Should have no errors (data is valid, just small)
		if len(validationResult.Errors) > 0 {
			t.Errorf("Unexpected validation errors: %v", validationResult.Errors)
		}
	})

	// Step 4: Write output
	w := writer.New(cfg)
	if err := w.WriteAll(convertedData); err != nil {
		t.Fatalf("WriteAll failed: %v", err)
	}

	// Verify output files
	t.Run("output", func(t *testing.T) {
		expectedFiles := []string{
			writer.FileSolarSystems,
			writer.FileRegions,
			writer.FileConstellations,
			writer.FileWormholeClasses,
			writer.FileShipTypes,
			writer.FileItemGroups,
			writer.FileSystemJumps,
		}

		for _, filename := range expectedFiles {
			path := filepath.Join(outputDir, filename)
			info, err := os.Stat(path)
			if err != nil {
				t.Errorf("Output file %s not created: %v", filename, err)
				continue
			}
			if info.Size() == 0 {
				t.Errorf("Output file %s is empty", filename)
			}
		}
	})
}

func TestIntegration_PassthroughFiles(t *testing.T) {
	// Create source directory with passthrough files
	srcDir, err := os.MkdirTemp("", "passthrough_src")
	if err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(srcDir) }()

	// Create output directory
	outputDir, err := os.MkdirTemp("", "passthrough_dst")
	if err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(outputDir) }()

	// Create passthrough files
	passthroughFiles := []string{
		"wormholes.json",
		"wormholeClasses.json",
		"effects.json",
		"triglavianSystems.json",
	}

	for _, filename := range passthroughFiles {
		content := []byte(`{"test": "` + filename + `"}`)
		if err := os.WriteFile(filepath.Join(srcDir, filename), content, 0644); err != nil {
			t.Fatalf("failed to create %s: %v", filename, err)
		}
	}

	cfg := &config.Config{
		OutputDir:   outputDir,
		PrettyPrint: true,
		Verbose:     false,
	}

	w := writer.New(cfg)
	if err := w.CopyPassthroughFiles(srcDir); err != nil {
		t.Fatalf("CopyPassthroughFiles failed: %v", err)
	}

	// Verify files were copied
	for _, filename := range passthroughFiles {
		dstPath := filepath.Join(outputDir, filename)
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			t.Errorf("Passthrough file %s not copied", filename)
		}
	}
}

func TestIntegration_SortingConsistency(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	cfg := &config.Config{
		SDEPath: sdeDir,
		Verbose: false,
	}

	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify all data is sorted correctly
	t.Run("regions sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.Universe.Regions); i++ {
			if convertedData.Universe.Regions[i-1].RegionID >= convertedData.Universe.Regions[i].RegionID {
				t.Errorf("Regions not sorted: %d >= %d",
					convertedData.Universe.Regions[i-1].RegionID,
					convertedData.Universe.Regions[i].RegionID)
			}
		}
	})

	t.Run("constellations sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.Universe.Constellations); i++ {
			if convertedData.Universe.Constellations[i-1].ConstellationID >= convertedData.Universe.Constellations[i].ConstellationID {
				t.Errorf("Constellations not sorted: %d >= %d",
					convertedData.Universe.Constellations[i-1].ConstellationID,
					convertedData.Universe.Constellations[i].ConstellationID)
			}
		}
	})

	t.Run("solar systems sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.Universe.SolarSystems); i++ {
			if convertedData.Universe.SolarSystems[i-1].SolarSystemID >= convertedData.Universe.SolarSystems[i].SolarSystemID {
				t.Errorf("Solar systems not sorted: %d >= %d",
					convertedData.Universe.SolarSystems[i-1].SolarSystemID,
					convertedData.Universe.SolarSystems[i].SolarSystemID)
			}
		}
	})

	t.Run("ship types sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.ShipTypes); i++ {
			if convertedData.ShipTypes[i-1].TypeID >= convertedData.ShipTypes[i].TypeID {
				t.Errorf("Ship types not sorted: %d >= %d",
					convertedData.ShipTypes[i-1].TypeID,
					convertedData.ShipTypes[i].TypeID)
			}
		}
	})

	t.Run("item groups sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.ItemGroups); i++ {
			if convertedData.ItemGroups[i-1].GroupID >= convertedData.ItemGroups[i].GroupID {
				t.Errorf("Item groups not sorted: %d >= %d",
					convertedData.ItemGroups[i-1].GroupID,
					convertedData.ItemGroups[i].GroupID)
			}
		}
	})

	t.Run("system jumps sorted", func(t *testing.T) {
		for i := 1; i < len(convertedData.SystemJumps); i++ {
			prev := convertedData.SystemJumps[i-1]
			curr := convertedData.SystemJumps[i]
			if prev.FromSolarSystemID > curr.FromSolarSystemID {
				t.Errorf("System jumps not sorted by FromSolarSystemID: %d > %d",
					prev.FromSolarSystemID, curr.FromSolarSystemID)
			}
			if prev.FromSolarSystemID == curr.FromSolarSystemID &&
				prev.ToSolarSystemID >= curr.ToSolarSystemID {
				t.Errorf("System jumps not sorted by ToSolarSystemID: %d >= %d",
					prev.ToSolarSystemID, curr.ToSolarSystemID)
			}
		}
	})
}

func TestIntegration_DataIntegrity(t *testing.T) {
	sdeDir := createTestSDE(t)
	defer func() { _ = os.RemoveAll(sdeDir) }()

	cfg := &config.Config{
		SDEPath: sdeDir,
		Verbose: false,
	}

	p := parser.New(cfg, sdeDir)
	parseResult, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll failed: %v", err)
	}

	tr := transformer.New(cfg)
	convertedData, err := tr.Transform(parseResult)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify referential integrity
	t.Run("constellation region references", func(t *testing.T) {
		regionIDs := make(map[int64]bool)
		for _, r := range convertedData.Universe.Regions {
			regionIDs[r.RegionID] = true
		}

		for _, c := range convertedData.Universe.Constellations {
			if !regionIDs[c.RegionID] {
				t.Errorf("Constellation %d references non-existent region %d",
					c.ConstellationID, c.RegionID)
			}
		}
	})

	t.Run("solar system references", func(t *testing.T) {
		regionIDs := make(map[int64]bool)
		for _, r := range convertedData.Universe.Regions {
			regionIDs[r.RegionID] = true
		}

		constellationIDs := make(map[int64]bool)
		for _, c := range convertedData.Universe.Constellations {
			constellationIDs[c.ConstellationID] = true
		}

		for _, s := range convertedData.Universe.SolarSystems {
			if !regionIDs[s.RegionID] {
				t.Errorf("Solar system %d references non-existent region %d",
					s.SolarSystemID, s.RegionID)
			}
			if !constellationIDs[s.ConstellationID] {
				t.Errorf("Solar system %d references non-existent constellation %d",
					s.SolarSystemID, s.ConstellationID)
			}
		}
	})

	t.Run("system jump references", func(t *testing.T) {
		systemIDs := make(map[int64]bool)
		for _, s := range convertedData.Universe.SolarSystems {
			systemIDs[s.SolarSystemID] = true
		}

		for _, j := range convertedData.SystemJumps {
			if !systemIDs[j.FromSolarSystemID] {
				t.Errorf("System jump references non-existent from system %d",
					j.FromSolarSystemID)
			}
			if !systemIDs[j.ToSolarSystemID] {
				t.Errorf("System jump references non-existent to system %d",
					j.ToSolarSystemID)
			}
		}
	})
}
