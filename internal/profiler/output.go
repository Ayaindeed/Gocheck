package profiler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// OutputJSON writes the profile as JSON
func OutputJSON(profile *DataProfile, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(profile)
}

// OutputCSV writes the profile as CSV
func OutputCSV(profile *DataProfile, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	headers := []string{
		"Column", "Type", "Missing Count", "Missing %",
		"Unique Count", "Min", "Max", "Mean", "Median", "Std Dev",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write each column
	for _, col := range profile.Columns {
		record := []string{
			col.Name,
			col.InferredType,
			fmt.Sprintf("%d", col.MissingCount),
			fmt.Sprintf("%.2f", col.MissingPct),
			fmt.Sprintf("%d", col.UniqueCount),
			FormatFloat(col.Min),
			FormatFloat(col.Max),
			FormatFloat(col.Mean),
			FormatFloat(col.Median),
			FormatFloat(col.StdDev),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// OutputTable writes the profile as a formatted table
func OutputTable(profile *DataProfile, w io.Writer) error {
	// Print summary
	fmt.Fprintf(w, "\n📊 Dataset Profile: %s\n", profile.FileName)
	fmt.Fprintf(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Fprintf(w, "Rows:    %d\n", profile.RowCount)
	fmt.Fprintf(w, "Columns: %d\n", profile.ColumnCount)
	fmt.Fprintf(w, "Time:    %s\n\n", profile.ProcessedTime.Format("2006-01-02 15:04:05"))

	// Create table for column details
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Column", "Type", "Missing", "Missing %", "Unique", "Stats"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetColWidth(40)

	for _, col := range profile.Columns {
		stats := buildStatsString(&col)
		missingStr := fmt.Sprintf("%d", col.MissingCount)
		missingPctStr := fmt.Sprintf("%.1f%%", col.MissingPct)
		uniqueStr := fmt.Sprintf("%d", col.UniqueCount)

		table.Append([]string{
			col.Name,
			col.InferredType,
			missingStr,
			missingPctStr,
			uniqueStr,
			stats,
		})
	}

	table.Render()

	// Print data quality summary
	fmt.Fprintf(w, "\n🔍 Data Quality Summary:\n")
	fmt.Fprintf(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	columnsWithMissing := 0
	totalMissing := int64(0)
	
	for _, col := range profile.Columns {
		if col.MissingCount > 0 {
			columnsWithMissing++
			totalMissing += col.MissingCount
		}
	}

	totalCells := profile.RowCount * int64(profile.ColumnCount)
	overallMissingPct := float64(totalMissing) / float64(totalCells) * 100

	fmt.Fprintf(w, "Columns with missing values: %d / %d\n", columnsWithMissing, profile.ColumnCount)
	fmt.Fprintf(w, "Total missing cells: %d / %d (%.2f%%)\n", totalMissing, totalCells, overallMissingPct)
	
	if overallMissingPct > 10 {
		fmt.Fprintf(w, "⚠️  Warning: High percentage of missing data\n")
	} else if overallMissingPct > 0 {
		fmt.Fprintf(w, "✓ Acceptable level of missing data\n")
	} else {
		fmt.Fprintf(w, "✓ No missing data detected\n")
	}

	fmt.Fprintf(w, "\n")

	return nil
}

// buildStatsString creates a summary string for column statistics
func buildStatsString(col *ColumnProfile) string {
	var parts []string

	if col.InferredType == "numeric" {
		if col.Min != nil && col.Max != nil {
			parts = append(parts, fmt.Sprintf("Range: [%.2f, %.2f]", *col.Min, *col.Max))
		}
		if col.Mean != nil {
			parts = append(parts, fmt.Sprintf("Mean: %.2f", *col.Mean))
		}
		if col.Median != nil {
			parts = append(parts, fmt.Sprintf("Median: %.2f", *col.Median))
		}
	} else if col.InferredType == "string" {
		if col.MinLength != nil && col.MaxLength != nil {
			parts = append(parts, fmt.Sprintf("Len: [%d, %d]", *col.MinLength, *col.MaxLength))
		}
		if col.AvgLength != nil {
			parts = append(parts, fmt.Sprintf("Avg: %.1f", *col.AvgLength))
		}
	}

	if len(col.SampleValues) > 0 {
		samples := strings.Join(col.SampleValues, ", ")
		if len(samples) > 50 {
			samples = samples[:47] + "..."
		}
		parts = append(parts, fmt.Sprintf("Ex: %s", samples))
	}

	return strings.Join(parts, "\n")
}
