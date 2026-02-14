package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ayaindeed/gocheck/internal/profiler"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	verbose      bool
)

var rootCmd = &cobra.Command{
	Use:   "gocheck [file]",
	Short: "Fast dataset profiling and data quality checks",
	Long: `gocheck is a fast Go-based CLI for profiling datasets and running 
basic data-quality checks on CSV and Parquet files.

It provides:
- Row counts
- Missing value analysis
- Type inference
- Summary statistics (min, max, mean, median, etc.)`,
	Args: cobra.ExactArgs(1),
	RunE: runProfile,
}

func init() {
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table, json, or csv")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func Execute() error {
	return rootCmd.Execute()
}

func runProfile(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Determine file type
	ext := strings.ToLower(filepath.Ext(filePath))
	
	var profile *profiler.DataProfile
	var err error

	if verbose {
		fmt.Fprintf(os.Stderr, "Processing file: %s\n", filePath)
	}

	switch ext {
	case ".csv":
		profile, err = profiler.ProfileCSV(filePath)
	case ".parquet":
		profile, err = profiler.ProfileParquet(filePath)
	default:
		return fmt.Errorf("unsupported file type: %s (supported: .csv, .parquet)", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to profile file: %w", err)
	}

	// Output results
	switch outputFormat {
	case "json":
		return profiler.OutputJSON(profile, os.Stdout)
	case "csv":
		return profiler.OutputCSV(profile, os.Stdout)
	case "table":
		return profiler.OutputTable(profile, os.Stdout)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}
