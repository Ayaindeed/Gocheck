package profiler

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/parquet-go/parquet-go"
)

// ProfileParquet reads and profiles a Parquet file
func ProfileParquet(filePath string) (*DataProfile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Parquet file: %w", err)
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Create parquet file reader
	pqFile, err := parquet.OpenFile(file, fileInfo.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to open Parquet file: %w", err)
	}

	schema := pqFile.Schema()
	fields := schema.Fields()
	numCols := len(fields)

	// Initialize column data storage
	columnData := make([][]string, numCols)
	columnNames := make([]string, numCols)
	for i, field := range fields {
		columnNames[i] = field.Name()
		columnData[i] = make([]string, 0, 1000)
	}

	// Read rows
	rowCount := int64(0)
	
	// Create a reader for row iteration
	reader := parquet.NewReader(pqFile)
	defer reader.Close()

	buf := make([]parquet.Row, 1)
	for {
		n, err := reader.ReadRows(buf)
		if err != nil || n == 0 {
			break
		}

		rowCount++
		for colIdx := 0; colIdx < numCols && colIdx < len(buf[0]); colIdx++ {
			strValue := convertValueToString(buf[0][colIdx])
			columnData[colIdx] = append(columnData[colIdx], strValue)
		}
	}

	// Profile each column
	columns := make([]ColumnProfile, numCols)
	for i := 0; i < numCols; i++ {
		columns[i] = profileColumn(columnNames[i], columnData[i], rowCount)
	}

	return &DataProfile{
		FileName:      filepath.Base(filePath),
		RowCount:      rowCount,
		ColumnCount:   numCols,
		Columns:       columns,
		ProcessedTime: time.Now(),
	}, nil
}

// convertValueToString converts a parquet.Value to string representation
func convertValueToString(val parquet.Value) string {
	if val.IsNull() {
		return ""
	}

	switch val.Kind() {
	case parquet.Boolean:
		return fmt.Sprintf("%v", val.Boolean())
	case parquet.Int32:
		return fmt.Sprintf("%d", val.Int32())
	case parquet.Int64:
		return fmt.Sprintf("%d", val.Int64())
	case parquet.Int96:
		return fmt.Sprintf("%v", val.Int96())
	case parquet.Float:
		return fmt.Sprintf("%f", val.Float())
	case parquet.Double:
		return fmt.Sprintf("%f", val.Double())
	case parquet.ByteArray:
		return string(val.ByteArray())
	default:
		return val.String()
	}
}
