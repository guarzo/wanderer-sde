// Package writer provides JSON output generation for the SDE converter.
package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/models"
)

// OutputFiles defines the names of generated JSON files.
const (
	FileSolarSystems    = "mapSolarSystems.json"
	FileRegions         = "mapRegions.json"
	FileConstellations  = "mapConstellations.json"
	FileWormholeClasses = "mapLocationWormholeClasses.json"
	FileShipTypes       = "invTypes.json"
	FileItemGroups      = "invGroups.json"
	FileSystemJumps     = "mapSolarSystemJumps.json"
)

// PassthroughFiles lists the community-maintained JSON files to copy.
var PassthroughFiles = []string{
	"wormholes.json",
	"wormholeClasses.json",
	"wormholeClassesInfo.json",
	"wormholeSystems.json",
	"triglavianSystems.json",
	"effects.json",
	"shatteredConstellations.json",
	"sunTypes.json",
	"triglavianEffectsByFaction.json",
}

// JSONWriter handles writing converted data to JSON files.
type JSONWriter struct {
	config    *config.Config
	outputDir string
	pretty    bool
}

// New creates a new JSONWriter with the given configuration.
func New(cfg *config.Config) *JSONWriter {
	return &JSONWriter{
		config:    cfg,
		outputDir: cfg.OutputDir,
		pretty:    cfg.PrettyPrint,
	}
}

// WriteAll writes all converted data to JSON files.
func (w *JSONWriter) WriteAll(data *models.ConvertedData) error {
	// Ensure output directory exists
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if w.config.Verbose {
		fmt.Printf("Writing JSON files to: %s\n", w.outputDir)
	}

	// Write all data files
	if err := w.WriteSolarSystems(data.Universe.SolarSystems); err != nil {
		return fmt.Errorf("failed to write solar systems: %w", err)
	}

	if err := w.WriteRegions(data.Universe.Regions); err != nil {
		return fmt.Errorf("failed to write regions: %w", err)
	}

	if err := w.WriteConstellations(data.Universe.Constellations); err != nil {
		return fmt.Errorf("failed to write constellations: %w", err)
	}

	if err := w.WriteWormholeClasses(data.WormholeClasses); err != nil {
		return fmt.Errorf("failed to write wormhole classes: %w", err)
	}

	if err := w.WriteShipTypes(data.ShipTypes); err != nil {
		return fmt.Errorf("failed to write ship types: %w", err)
	}

	if err := w.WriteGroups(data.ItemGroups); err != nil {
		return fmt.Errorf("failed to write item groups: %w", err)
	}

	if err := w.WriteSystemJumps(data.SystemJumps); err != nil {
		return fmt.Errorf("failed to write system jumps: %w", err)
	}

	return nil
}

// WriteSolarSystems writes solar system data to JSON.
func (w *JSONWriter) WriteSolarSystems(systems []models.SolarSystem) error {
	return w.writeJSON(FileSolarSystems, systems)
}

// WriteRegions writes region data to JSON.
func (w *JSONWriter) WriteRegions(regions []models.Region) error {
	return w.writeJSON(FileRegions, regions)
}

// WriteConstellations writes constellation data to JSON.
func (w *JSONWriter) WriteConstellations(constellations []models.Constellation) error {
	return w.writeJSON(FileConstellations, constellations)
}

// WriteWormholeClasses writes wormhole class data to JSON.
func (w *JSONWriter) WriteWormholeClasses(classes []models.WormholeClassLocation) error {
	return w.writeJSON(FileWormholeClasses, classes)
}

// WriteShipTypes writes ship type data to JSON.
func (w *JSONWriter) WriteShipTypes(ships []models.ShipType) error {
	return w.writeJSON(FileShipTypes, ships)
}

// WriteGroups writes item group data to JSON.
func (w *JSONWriter) WriteGroups(groups []models.ItemGroup) error {
	return w.writeJSON(FileItemGroups, groups)
}

// WriteSystemJumps writes system jump data to JSON.
func (w *JSONWriter) WriteSystemJumps(jumps []models.SystemJump) error {
	return w.writeJSON(FileSystemJumps, jumps)
}

// CopyPassthroughFiles copies community-maintained JSON files from the source directory.
func (w *JSONWriter) CopyPassthroughFiles(sourceDir string) error {
	if sourceDir == "" {
		return nil
	}

	if w.config.Verbose {
		fmt.Printf("Copying passthrough files from: %s\n", sourceDir)
	}

	var copied, skipped int
	for _, filename := range PassthroughFiles {
		srcPath := filepath.Join(sourceDir, filename)
		dstPath := filepath.Join(w.outputDir, filename)

		// Check if source file exists
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			if w.config.Verbose {
				fmt.Printf("  Skipping %s (not found)\n", filename)
			}
			skipped++
			continue
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", filename, err)
		}

		if w.config.Verbose {
			fmt.Printf("  Copied %s\n", filename)
		}
		copied++
	}

	if w.config.Verbose {
		fmt.Printf("Passthrough complete: %d copied, %d skipped\n", copied, skipped)
	}

	return nil
}

// writeJSON marshals data to JSON and writes it to a file.
func (w *JSONWriter) writeJSON(filename string, data interface{}) error {
	path := filepath.Join(w.outputDir, filename)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)
	if w.pretty {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to %s: %w", path, err)
	}

	if w.config.Verbose {
		fmt.Printf("  Wrote %s\n", filename)
	}

	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
