package profiler

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ProfileCSV reads and profiles a CSV file
func ProfileCSV(filePath string) (*DataProfile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.ReuseRecord = true
	reader.LazyQuotes = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	numCols := len(headers)
	headersCopy := make([]string, numCols)
	copy(headersCopy, headers)

	// Initialize column data storage
	columnData := make([][]string, numCols)
	for i := range columnData {
		columnData[i] = make([]string, 0, 1000)
	}

	rowCount := int64(0)

	// Read all rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row %d: %w", rowCount+2, err)
		}

		rowCount++

		// Store column values
		for i := 0; i < numCols && i < len(record); i++ {
			columnData[i] = append(columnData[i], record[i])
		}
	}

	// Profile each column
	columns := make([]ColumnProfile, numCols)
	for i := 0; i < numCols; i++ {
		columns[i] = profileColumn(headersCopy[i], columnData[i], rowCount)
	}

	return &DataProfile{
		FileName:      filepath.Base(filePath),
		RowCount:      rowCount,
		ColumnCount:   numCols,
		Columns:       columns,
		ProcessedTime: time.Now(),
	}, nil
}

// profileColumn creates a profile for a single column
func profileColumn(name string, values []string, totalRows int64) ColumnProfile {
	missingCount := int64(0)
	for _, v := range values {
		if v == "" {
			missingCount++
		}
	}

	inferredType := InferType(values)
	missingPct := float64(missingCount) / float64(totalRows) * 100
	uniqueCount := CountUnique(values)

	profile := ColumnProfile{
		Name:         name,
		InferredType: inferredType,
		MissingCount: missingCount,
		MissingPct:   missingPct,
		UniqueCount:  uniqueCount,
		SampleValues: GetSampleValues(values, 5),
	}

	// Calculate type-specific statistics
	if inferredType == "numeric" {
		profile.Min, profile.Max, profile.Mean, profile.Median, profile.StdDev = CalculateNumericStats(values)
	} else if inferredType == "string" {
		profile.MinLength, profile.MaxLength, profile.AvgLength = CalculateStringStats(values)
	}

	return profile
}
