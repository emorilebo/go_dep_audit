package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/emori/go_dep_audit/pkg/audit"
	"github.com/spf13/cobra"
)

var (
	failThreshold int
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check dependencies against thresholds and exit with code",
	RunE:  runCheck,
}

func init() {
	checkCmd.Flags().IntVar(&failThreshold, "fail-threshold", 50, "Fail if any module score is below this threshold")
}

func runCheck(cmd *cobra.Command, args []string) error {
	config := audit.AuditConfig{
		ProjectPath: projectPath,
		Scoring:     audit.DefaultScoringConfig(),
	}

	results, err := audit.AuditModules(context.Background(), config)
	if err != nil {
		return err
	}

	failed := false
	for _, res := range results {
		if res.HealthScore < failThreshold {
			fmt.Printf("FAIL: %s@%s (Score: %d) is below threshold %d\n", 
				res.Path, res.Version, res.HealthScore, failThreshold)
			failed = true
		}
	}

	if failed {
		os.Exit(1)
	}

	fmt.Println("All checks passed.")
	return nil
}
