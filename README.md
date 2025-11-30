# wanderer-sde

A Go tool to convert EVE Online's Static Data Export (SDE) YAML files into JSON format compatible with [Wanderer](https://github.com/wanderer-industries/wanderer).

## Overview

This converter eliminates the dependency on third-party CSV dumps (like Fuzzwork) by processing the official SDE directly from CCP. It parses the YAML files and generates JSON output that Wanderer can consume.

### Features

- Downloads the latest SDE directly from CCP
- Parses YAML files with parallel processing for performance
- Generates JSON files in Wanderer's expected format
- Supports passthrough of community-maintained data files
- Version tracking to avoid redundant downloads
- Cross-platform support (Linux, macOS, Windows)

## Installation

### From Source

```bash
git clone https://github.com/guarzo/wanderer-sde.git
cd wanderer-sde
make build
```

The binary will be available at `bin/sdeconvert`.

### Cross-Platform Builds

```bash
make build-all
```

This creates binaries for:
- Linux (amd64)
- macOS (Intel and ARM)
- Windows (amd64)

## Usage

### Quick Start

Download and convert the latest SDE in one command:

```bash
./bin/sdeconvert --download --output ./output
```

### Command Line Options

```
Usage:
  sdeconvert [flags]
  sdeconvert [command]

Available Commands:
  help        Help about any command
  version     Print the version number

Flags:
  -d, --download             Download latest SDE from CCP
  -h, --help                 help for sdeconvert
  -o, --output string        Output directory for JSON files (default "./output")
  -p, --passthrough string   Directory with Wanderer JSON files to copy
      --pretty               Pretty-print JSON output (default true)
  -s, --sde-path string      Path to SDE directory or ZIP file
      --sde-url string       URL to download SDE from
  -v, --verbose              Enable verbose output
  -w, --workers int          Number of parallel workers (default 4)
```

### Usage Examples

#### Download and Convert Latest SDE

The simplest use case - download the latest SDE from CCP and convert it:

```bash
sdeconvert --download --output ./output
```

#### Convert an Existing SDE Directory

If you already have the SDE extracted locally:

```bash
sdeconvert --sde-path /path/to/sde --output ./output
```

#### Include Wanderer Passthrough Files

Some JSON files contain community-maintained data (wormhole info, effects, etc.) that should be copied as-is from the Wanderer repository:

```bash
sdeconvert --sde-path ./sde \
  --output ./output \
  --passthrough /path/to/wanderer/priv/repo/data
```

#### Verbose Mode with Custom Worker Count

For debugging or monitoring large conversions:

```bash
sdeconvert --download --output ./output --verbose --workers 8
```

#### Custom SDE URL

Use a specific SDE version or mirror:

```bash
sdeconvert --download \
  --sde-url "https://example.com/custom-sde.zip" \
  --output ./output
```

## Output Files

The converter generates the following JSON files:

### Generated from SDE

| File | Description | Source |
|------|-------------|--------|
| `mapSolarSystems.json` | Solar systems with security status, region/constellation IDs, and sun type | `sde/fsd/universe/*/solarsystem.yaml` |
| `mapRegions.json` | Region ID and name mappings | `sde/fsd/universe/*/region.yaml` |
| `mapConstellations.json` | Constellation ID, name, and region mappings | `sde/fsd/universe/*/constellation.yaml` |
| `mapLocationWormholeClasses.json` | Wormhole class assignments for systems/regions | `sde/bsd/mapLocationWormholeClasses.yaml` |
| `invTypes.json` | Ship type definitions (filtered to category 6) | `sde/fsd/typeIDs.yaml` |
| `invGroups.json` | Ship group definitions | `sde/fsd/groupIDs.yaml` |
| `mapSolarSystemJumps.json` | Stargate connections between systems | `sde/bsd/mapSolarSystemJumps.yaml` |

### Passthrough Files (Community-Maintained)

These files are copied from the Wanderer data directory when `--passthrough` is specified:

| File | Description |
|------|-------------|
| `wormholes.json` | Wormhole type definitions |
| `wormholeClasses.json` | Wormhole class definitions |
| `wormholeClassesInfo.json` | Detailed wormhole class information |
| `wormholeSystems.json` | Known wormhole system data |
| `triglavianSystems.json` | Triglavian invasion system data |
| `effects.json` | System effect definitions |
| `shatteredConstellations.json` | Shattered wormhole constellation data |
| `sunTypes.json` | Sun type definitions |
| `triglavianEffectsByFaction.json` | Triglavian effects by faction |

## Data Formats

### Solar Systems (`mapSolarSystems.json`)

```json
[
  {
    "solarSystemID": 30000142,
    "regionID": 10000002,
    "constellationID": 20000020,
    "solarSystemName": "Jita",
    "sunTypeID": 6,
    "security": 0.9459
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| `solarSystemID` | int64 | Unique system identifier |
| `regionID` | int64 | Parent region ID |
| `constellationID` | int64 | Parent constellation ID |
| `solarSystemName` | string | Display name of the system |
| `sunTypeID` | int64 | Type ID of the system's star (optional) |
| `security` | float64 | Security status (-1.0 to 1.0) |

### Regions (`mapRegions.json`)

```json
[
  {
    "regionID": 10000002,
    "regionName": "The Forge"
  }
]
```

### Constellations (`mapConstellations.json`)

```json
[
  {
    "constellationID": 20000020,
    "constellationName": "Kimotoro",
    "regionID": 10000002
  }
]
```

### Wormhole Classes (`mapLocationWormholeClasses.json`)

```json
[
  {
    "locationID": 10000002,
    "wormholeClassID": 7
  }
]
```

Location IDs can be regions, constellations, or solar systems. Wormhole class IDs:
- 1-6: C1-C6 wormhole space
- 7: High-sec (0.5-1.0 security)
- 8: Low-sec (0.1-0.4 security)
- 9: Null-sec (0.0 and below)
- 12: Thera
- 13: Shattered wormholes
- 14-18: Drifter wormholes
- 25: Pochven (Triglavian space)

### Ship Types (`invTypes.json`)

```json
[
  {
    "typeID": 587,
    "groupID": 25,
    "typeName": "Rifter",
    "mass": 1350000.0,
    "volume": 27500.0,
    "capacity": 125.0
  }
]
```

Only ships (category ID 6) are included in this output.

### Item Groups (`invGroups.json`)

```json
[
  {
    "groupID": 25,
    "categoryID": 6,
    "groupName": "Frigate"
  }
]
```

Only ship groups (category ID 6) are included.

### System Jumps (`mapSolarSystemJumps.json`)

```json
[
  {
    "fromSolarSystemID": 30000142,
    "toSolarSystemID": 30000144
  }
]
```

Represents stargate connections between solar systems. Each connection appears once (not duplicated in reverse).

## Development

### Prerequisites

- Go 1.22 or later
- Make (optional, for convenience commands)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to $GOPATH/bin
make install
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run integration tests (requires downloaded SDE)
go test -v ./internal/... -run Integration
```

### Project Structure

```
wanderer-sde/
├── cmd/
│   └── sdeconvert/
│       └── main.go              # CLI entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── downloader/
│   │   ├── downloader.go        # SDE download & extraction
│   │   └── version.go           # Version checking
│   ├── parser/
│   │   ├── parser.go            # Main parser orchestration
│   │   ├── universe.go          # Universe file parsing
│   │   ├── types.go             # typeIDs.yaml parsing
│   │   ├── groups.go            # groupIDs.yaml parsing
│   │   ├── categories.go        # categoryIDs.yaml parsing
│   │   └── jumps.go             # mapSolarSystemJumps parsing
│   ├── transformer/
│   │   ├── transformer.go       # Data transformation logic
│   │   ├── security.go          # Security status calculation
│   │   └── filters.go           # Ship category filtering
│   ├── writer/
│   │   └── json_writer.go       # JSON output generation
│   └── models/
│       ├── sde.go               # SDE data structures
│       └── wanderer.go          # Wanderer output structures
├── pkg/
│   └── yaml/
│       └── yaml.go              # YAML utilities
├── plans/
│   └── IMPLEMENTATION_PLAN.md   # Development roadmap
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── CONTRIBUTING.md
└── LICENSE
```

### Code Architecture

The converter follows a pipeline architecture:

1. **Downloader**: Downloads and extracts the SDE from CCP
2. **Parser**: Reads YAML files and converts to internal Go structs
3. **Transformer**: Applies business logic (security calculation, filtering)
4. **Writer**: Serializes data to JSON files

Each component is isolated and testable independently.

## Data Sources

- **SDE**: [EVE Online Static Data Export](https://developers.eveonline.com/docs/services/static-data/)
- **Wanderer**: [wanderer-industries/wanderer](https://github.com/wanderer-industries/wanderer)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- CCP Games for providing the EVE Online Static Data Export
- The Wanderer project for the data format specifications
- Fuzzwork for the original CSV dump service that inspired this tool
