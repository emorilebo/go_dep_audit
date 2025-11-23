package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/emori/go_dep_audit/pkg/audit"
	"github.com/spf13/cobra"
)

var (
	outputJSON string
	outputMD   string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate detailed audit report",
	RunE:  runReport,
}

func init() {
	reportCmd.Flags().StringVar(&outputJSON, "output-json", "", "Path to save JSON report")
	reportCmd.Flags().StringVar(&outputMD, "output-md", "", "Path to save Markdown report")
}

func runReport(cmd *cobra.Command, args []string) error {
	config := audit.AuditConfig{
		ProjectPath: projectPath,
		Scoring:     audit.DefaultScoringConfig(),
	}

	results, err := audit.AuditModules(context.Background(), config)
	if err != nil {
		return err
	}

	if outputJSON != "" {
		if err := generateJSONReport(results, outputJSON); err != nil {
			return err
		}
		fmt.Printf("JSON report saved to %s\n", outputJSON)
	}

	if outputMD != "" {
		if err := generateMarkdownReport(results, outputMD); err != nil {
			return err
		}
		fmt.Printf("Markdown report saved to %s\n", outputMD)
	}

	if outputJSON == "" && outputMD == "" {
		// Default to JSON to stdout
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	return nil
}

func generateJSONReport(results []audit.ModuleHealth, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func generateMarkdownReport(results []audit.ModuleHealth, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "# Dependency Audit Report")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "| Module | Version | Score | Category | License |")
	fmt.Fprintln(file, "|--------|---------|-------|----------|---------|")
	
	for _, res := range results {
		fmt.Fprintf(file, "| %s | %s | %d | %s | %s |\n", 
			res.Path, res.Version, res.HealthScore, res.HealthCategory, res.License)
	}
	
	return nil
}
