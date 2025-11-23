# Go Dependency Audit Tool

A comprehensive dependency audit tool for Go projects that analyzes supply-chain risk, maintenance health, license compatibility, and dependency footprint.

## Features

- **Health Scoring**: specific heuristics to score dependencies based on recency, version frequency, and community activity.
- **License Risk**: Detects and classifies licenses (Permissive, Copyleft, Restrictive).
- **Footprint Analysis**: Estimates dependency bloat.
- **CLI & Library**: Use as a standalone CLI tool or embed in your Go programs.

## Installation

```bash
go install github.com/emori/go_dep_audit/cmd/go-dep-audit@latest
```

## Usage

### Scan a Project

```bash
go-dep-audit scan --project-path /path/to/project
```

### Generate Report

```bash
go-dep-audit report --project-path . --output-json report.json --output-md report.md
```

### CI/CD Check

Fail the build if any dependency has a health score below 50:

```bash
go-dep-audit check --fail-threshold 50
```

## Configuration

You can configure the tool using flags or a config file (coming soon).

### Default Scoring Heuristics

- **Recency (40%)**: Penalizes modules not updated in the last 6 months.
- **Version Frequency (20%)**: Rewards active release cycles.
- **Commit Activity (20%)**: Rewards frequent commits (requires repo metadata).
- **Community (20%)**: Rewards stars and contributors (requires repo metadata).

## Library Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/emori/go_dep_audit/pkg/audit"
)

func main() {
    config := audit.AuditConfig{
        ProjectPath: ".",
        Scoring: audit.DefaultScoringConfig(),
    }

    results, err := audit.AuditModules(context.Background(), config)
    if err != nil {
        panic(err)
    }

    for _, res := range results {
        fmt.Printf("%s: %d (%s)\n", res.Path, res.HealthScore, res.HealthCategory)
    }
}
```

## Development

### Prerequisites

- Go 1.21+

### Build

```bash
go mod tidy
go build ./cmd/go-dep-audit
```

### Test

```bash
go test ./...
```

## License

MIT
