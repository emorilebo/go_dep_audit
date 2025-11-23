package audit

import (
	"context"
	"fmt"
	"sync"
)

// AuditModules performs a full audit of the project's dependencies
func AuditModules(ctx context.Context, config AuditConfig) ([]ModuleHealth, error) {
	// 1. Get Dependency Graph
	modules, err := GetModuleGraph(ctx, config.ProjectPath)
	if err != nil {
		// Fallback to parsing go.mod if go list fails (e.g. no go installed)
		// This is critical for the agent environment where go might be missing
		fmt.Printf("Warning: 'go list' failed (%v), falling back to simple go.mod parsing\n", err)
		modules, err = ParseGoMod(config.ProjectPath + "/go.mod")
		if err != nil {
			return nil, fmt.Errorf("failed to parse modules: %w", err)
		}
	}

	// Filter modules based on config
	var targetModules []Module
	for _, m := range modules {
		// Skip the main module itself
		if m.Path == "" || m.Main {
			continue
		}
		
		// Check ignore patterns
		ignored := false
		for _, pattern := range config.IgnoreModules {
			if m.Path == pattern { // Simple match for now
				ignored = true
				break
			}
		}
		if ignored {
			continue
		}

		if !config.IncludeIndirect && m.Indirect {
			continue
		}
		targetModules = append(targetModules, m)
	}

	// 2. Fetch Metadata and Score (Parallel)
	results := make([]ModuleHealth, len(targetModules))
	fetcher := NewFetcher(config)
	
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit concurrency
	
	for i, mod := range targetModules {
		wg.Add(1)
		go func(i int, m Module) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release

			health, err := auditSingleModule(ctx, fetcher, m, config)
			if err != nil {
				// Log error?
				results[i] = ModuleHealth{
					Path:    m.Path,
					Version: m.Version,
					// Other fields zeroed
				}
			} else {
				results[i] = *health
			}
		}(i, mod)
	}
	
	wg.Wait()

	return results, nil
}

func auditSingleModule(ctx context.Context, fetcher *Fetcher, mod Module, config AuditConfig) (*ModuleHealth, error) {
	// Fetch metadata
	meta, err := fetcher.FetchModuleMetadata(ctx, mod.Path, mod.Version)
	if err != nil {
		// If fetch fails, we proceed with partial info
		meta = &ModuleMetadata{}
	}

	// Calculate Score
	score := CalculateHealthScore(meta, config.Scoring)
	category := CategorizeHealth(score, config.Scoring)

	// License
	license, _ := DetectLicense(ctx, mod.Path, mod.Version)
	licenseRisk := ClassifyLicense(license)

	// Footprint (estimated)
	// We don't have per-module footprint without graph analysis, so 0 for now
	footprintRisk := 0.0

	return &ModuleHealth{
		Path:           mod.Path,
		Version:        mod.Version,
		HealthScore:    score,
		HealthCategory: category,
		License:        license,
		LicenseRisk:    licenseRisk,
		FootprintRisk:  footprintRisk,
		LastPublished:  meta.LastCommitDate,
		DirectDep:      !mod.Indirect,
		Metadata:       meta,
	}, nil
}
