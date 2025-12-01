// Package main provides the CLI entry point for the SDE converter.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/guarzo/wanderer-sde/internal/config"
	"github.com/guarzo/wanderer-sde/internal/downloader"
	"github.com/guarzo/wanderer-sde/internal/parser"
	"github.com/guarzo/wanderer-sde/internal/transformer"
	"github.com/guarzo/wanderer-sde/internal/writer"
)

// Version is set at build time via ldflags.
var Version = "dev"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var cfg = config.NewConfig()

var rootCmd = &cobra.Command{
	Use:   "sdeconvert",
	Short: "Convert EVE SDE to Wanderer data format",
	Long: `Converts EVE Online's Static Data Export (SDE) YAML files
into the JSON format used by the Wanderer application.

This tool can download the latest SDE from CCP or use an existing
SDE directory, then parse the YAML files and generate JSON output
compatible with Wanderer's data format.`,
	Example: `  # Download latest SDE and convert to JSON
  sdeconvert --download --output ./output

  # Convert an existing SDE directory
  sdeconvert --sde-path ./sde --output ./output

  # Include Wanderer passthrough files (wormholes.json, etc.)
  sdeconvert --sde-path ./sde --output ./output --passthrough ../wanderer/priv/repo/data

  # Verbose mode with custom worker count
  sdeconvert --download --output ./output --verbose --workers 8`,
	RunE: runConversion,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number and build information for sdeconvert.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sdeconvert version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().StringVarP(&cfg.SDEPath, "sde-path", "s", "", "Path to SDE directory or ZIP file")
	rootCmd.Flags().StringVarP(&cfg.OutputDir, "output", "o", "./output", "Output directory for JSON files")
	rootCmd.Flags().BoolVarP(&cfg.DownloadSDE, "download", "d", false, "Download latest SDE from CCP")
	rootCmd.Flags().StringVarP(&cfg.PassthroughDir, "passthrough", "p", "", "Directory with Wanderer JSON files to copy")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVar(&cfg.PrettyPrint, "pretty", true, "Pretty-print JSON output")
	rootCmd.Flags().IntVarP(&cfg.Workers, "workers", "w", 4, "Number of parallel workers")
	rootCmd.Flags().StringVar(&cfg.SDEUrl, "sde-url", config.SDELatestURL, "URL to download SDE from")
}

func runConversion(cmd *cobra.Command, args []string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nInterrupt received, shutting down...")
		cancel()
	}()

	if cfg.Verbose {
		fmt.Println("Configuration:")
		fmt.Printf("  SDE Path:     %s\n", cfg.SDEPath)
		fmt.Printf("  Output Dir:   %s\n", cfg.OutputDir)
		fmt.Printf("  Download:     %v\n", cfg.DownloadSDE)
		fmt.Printf("  Passthrough:  %s\n", cfg.PassthroughDir)
		fmt.Printf("  Workers:      %d\n", cfg.Workers)
	}

	sdePath := cfg.SDEPath

	// Step 1: Download SDE if requested
	if cfg.DownloadSDE {
		dl := downloader.New(cfg)
		vc := downloader.NewVersionChecker(cfg)

		// Default SDE path when downloading is <output-dir>/sde
		if sdePath == "" {
			sdePath = filepath.Join(cfg.OutputDir, "sde")
		}

		// Check if update is needed
		needsUpdate, versionInfo, err := vc.NeedsUpdate(ctx, cfg.OutputDir)
		if err != nil {
			fmt.Printf("Warning: could not check SDE version: %v\n", err)
			needsUpdate = true // Proceed with download anyway
		}

		// Also check if the SDE directory exists
		if _, err := os.Stat(sdePath); os.IsNotExist(err) {
			needsUpdate = true
		}

		if needsUpdate {
			fmt.Println("Downloading latest SDE...")

			// Download to a temp location first
			downloadedPath, err := dl.DownloadAndExtract(ctx)
			if err != nil {
				return fmt.Errorf("failed to download SDE: %w", err)
			}

			// Remove old SDE directory if it exists
			if err := os.RemoveAll(sdePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Warning: could not remove old SDE: %v\n", err)
			}

			// Move downloaded SDE to the persistent location
			if err := os.Rename(downloadedPath, sdePath); err != nil {
				// If rename fails (cross-device), fall back to copy
				if err := copyDir(downloadedPath, sdePath); err != nil {
					return fmt.Errorf("failed to move SDE to %s: %w", sdePath, err)
				}
				_ = os.RemoveAll(downloadedPath)
			}

			fmt.Printf("SDE downloaded and extracted to: %s\n", sdePath)

			// Store the version
			if versionInfo != nil {
				if err := vc.StoreVersion(cfg.OutputDir, versionInfo.BuildNumber); err != nil {
					fmt.Printf("Warning: could not store version: %v\n", err)
				}
			}
		} else {
			fmt.Println("SDE is up to date, using cached version")
		}
	}

	if sdePath == "" {
		return fmt.Errorf("no SDE path available")
	}

	// Validate the SDE structure
	dl := downloader.New(cfg)
	if err := dl.Validate(sdePath); err != nil {
		return fmt.Errorf("SDE validation failed: %w", err)
	}

	fmt.Printf("Using SDE at: %s\n", sdePath)

	// Step 2: Parse SDE YAML files
	p := parser.New(cfg, sdePath)
	parseResult, err := p.ParseAll()
	if err != nil {
		return fmt.Errorf("failed to parse SDE: %w", err)
	}

	fmt.Printf("\nParsing complete:\n")
	fmt.Printf("  Regions:         %d\n", len(parseResult.Regions))
	fmt.Printf("  Constellations:  %d\n", len(parseResult.Constellations))
	fmt.Printf("  Solar Systems:   %d\n", len(parseResult.SolarSystems))
	fmt.Printf("  Types:           %d\n", len(parseResult.Types))
	fmt.Printf("  Groups:          %d\n", len(parseResult.Groups))
	fmt.Printf("  Categories:      %d\n", len(parseResult.Categories))
	fmt.Printf("  Wormhole Classes: %d\n", len(parseResult.WormholeClasses))
	fmt.Printf("  System Jumps:    %d\n", len(parseResult.SystemJumps))

	// Step 3: Transform data
	t := transformer.New(cfg)
	convertedData, err := t.Transform(parseResult)
	if err != nil {
		return fmt.Errorf("failed to transform data: %w", err)
	}

	// Validate the converted data
	validationResult := t.Validate(convertedData)
	fmt.Printf("\nValidation results:\n")
	fmt.Printf("  Regions:         %d\n", validationResult.Regions)
	fmt.Printf("  Constellations:  %d\n", validationResult.Constellations)
	fmt.Printf("  Solar Systems:   %d\n", validationResult.SolarSystems)
	fmt.Printf("  Ship Types:      %d\n", validationResult.ShipTypes)
	fmt.Printf("  Ship Groups:     %d\n", validationResult.ItemGroups)
	fmt.Printf("  Wormhole Classes: %d\n", validationResult.WormholeClasses)
	fmt.Printf("  System Jumps:    %d\n", validationResult.SystemJumps)

	if len(validationResult.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range validationResult.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(validationResult.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range validationResult.Errors {
			fmt.Printf("  - %s\n", err)
		}
		return fmt.Errorf("validation failed with %d errors", len(validationResult.Errors))
	}

	// Step 4: Write JSON output
	w := writer.New(cfg)
	if err := w.WriteAll(convertedData); err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}

	// Step 5: Copy passthrough files
	if cfg.PassthroughDir != "" {
		if err := w.CopyPassthroughFiles(cfg.PassthroughDir); err != nil {
			return fmt.Errorf("failed to copy passthrough files: %w", err)
		}
	}

	fmt.Printf("\nConversion complete! Output written to: %s\n", cfg.OutputDir)
	fmt.Printf("Generated files:\n")
	fmt.Printf("  - %s (%d systems)\n", writer.FileSolarSystems, len(convertedData.Universe.SolarSystems))
	fmt.Printf("  - %s (%d regions)\n", writer.FileRegions, len(convertedData.Universe.Regions))
	fmt.Printf("  - %s (%d constellations)\n", writer.FileConstellations, len(convertedData.Universe.Constellations))
	fmt.Printf("  - %s (%d classes)\n", writer.FileWormholeClasses, len(convertedData.WormholeClasses))
	fmt.Printf("  - %s (%d ships)\n", writer.FileShipTypes, len(convertedData.ShipTypes))
	fmt.Printf("  - %s (%d groups)\n", writer.FileItemGroups, len(convertedData.ItemGroups))
	fmt.Printf("  - %s (%d jumps)\n", writer.FileSystemJumps, len(convertedData.SystemJumps))

	return nil
}

// copyDir recursively copies a directory tree.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
