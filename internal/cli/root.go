package cli

import (
	"github.com/spf13/cobra"
)

var (
	projectPath   string
	configFile    string
	verboseOutput bool
)

var rootCmd = &cobra.Command{
	Use:   "go-dep-audit",
	Short: "Audit Go project dependencies for security and health",
	Long: `A comprehensive dependency audit tool for Go projects.
Analyzes supply-chain risk, maintenance health, license compatibility, 
and dependency footprint.`,
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectPath, "project-path", "p", ".", "Path to the Go project to audit")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verboseOutput, "verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(checkCmd)
}
