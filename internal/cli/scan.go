package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/emori/go_dep_audit/pkg/audit"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project dependencies and display summary",
	RunE:  runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	config := audit.AuditConfig{
		ProjectPath: projectPath,
		Scoring:     audit.DefaultScoringConfig(),
		// Load other defaults or from config file
	}

	fmt.Printf("Scanning dependencies in %s...\n", projectPath)
	
	results, err := audit.AuditModules(context.Background(), config)
	if err != nil {
		return err
	}

	// Summary counts
	counts := make(map[audit.HealthCategory]int)
	for _, res := range results {
		counts[res.HealthCategory]++
	}

	fmt.Println("\nAudit Summary:")
	fmt.Printf("Total Dependencies: %d\n", len(results))
	fmt.Printf("Healthy: %d\n", counts[audit.Healthy])
	fmt.Printf("Warning: %d\n", counts[audit.Warning])
	fmt.Printf("Stale:   %d\n", counts[audit.Stale])
	fmt.Printf("Risky:   %d\n", counts[audit.Risky])

	if counts[audit.Risky] > 0 || counts[audit.Stale] > 0 {
		fmt.Println("\nRisky/Stale Modules:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Module\tVersion\tScore\tCategory\tLicense")
		for _, res := range results {
			if res.HealthCategory == audit.Risky || res.HealthCategory == audit.Stale {
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", res.Path, res.Version, res.HealthScore, res.HealthCategory, res.License)
			}
		}
		w.Flush()
	}

	return nil
}
