package profiler

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// DataProfile represents the complete profile of a dataset
type DataProfile struct {
	FileName      string
	RowCount      int64
	ColumnCount   int
	Columns       []ColumnProfile
	ProcessedTime time.Time
}

// ColumnProfile represents statistics for a single column
type ColumnProfile struct {
	Name          string
	InferredType  string
	MissingCount  int64
	MissingPct    float64
	UniqueCount   int64
	
	// Numeric statistics
	Min           *float64
	Max           *float64
	Mean          *float64
	Median        *float64
	StdDev        *float64
	
	// String statistics
	MinLength     *int
	MaxLength     *int
	AvgLength     *float64
	
	// Sample values
	SampleValues  []string
}

// InferType attempts to determine the data type of values
func InferType(values []string) string {
	if len(values) == 0 {
		return "unknown"
	}

	numericCount := 0
	dateCount := 0
	boolCount := 0
	validCount := 0

	for _, v := range values {
		if v == "" {
			continue
		}
		validCount++

		// Check if numeric
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			numericCount++
			continue
		}

		// Check if boolean
		lower := v
		if lower == "true" || lower == "false" || lower == "True" || lower == "False" ||
			lower == "TRUE" || lower == "FALSE" || lower == "1" || lower == "0" {
			boolCount++
			continue
		}

		// Check if date
		if _, err := time.Parse("2006-01-02", v); err == nil {
			dateCount++
			continue
		}
		if _, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			dateCount++
			continue
		}
		if _, err := time.Parse("01/02/2006", v); err == nil {
			dateCount++
			continue
		}
	}

	if validCount == 0 {
		return "string"
	}

	threshold := float64(validCount) * 0.8

	if float64(numericCount) >= threshold {
		return "numeric"
	}
	if float64(dateCount) >= threshold {
		return "date"
	}
	if float64(boolCount) >= threshold {
		return "boolean"
	}

	return "string"
}

// CalculateNumericStats computes statistics for numeric columns
func CalculateNumericStats(values []string) (min, max, mean, median, stddev *float64) {
	var nums []float64
	
	for _, v := range values {
		if v == "" {
			continue
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			nums = append(nums, f)
		}
	}

	if len(nums) == 0 {
		return nil, nil, nil, nil, nil
	}

	// Min and Max
	minVal := nums[0]
	maxVal := nums[0]
	sum := 0.0

	for _, n := range nums {
		if n < minVal {
			minVal = n
		}
		if n > maxVal {
			maxVal = n
		}
		sum += n
	}

	min = &minVal
	max = &maxVal

	// Mean
	meanVal := sum / float64(len(nums))
	mean = &meanVal

	// Median
	sorted := make([]float64, len(nums))
	copy(sorted, nums)
	sort.Float64s(sorted)

	var medianVal float64
	if len(sorted)%2 == 0 {
		medianVal = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	} else {
		medianVal = sorted[len(sorted)/2]
	}
	median = &medianVal

	// Standard Deviation
	variance := 0.0
	for _, n := range nums {
		variance += math.Pow(n-meanVal, 2)
	}
	variance /= float64(len(nums))
	stddevVal := math.Sqrt(variance)
	stddev = &stddevVal

	return
}

// CalculateStringStats computes statistics for string columns
func CalculateStringStats(values []string) (minLen, maxLen *int, avgLen *float64) {
	validValues := 0
	totalLen := 0
	minLength := math.MaxInt32
	maxLength := 0

	for _, v := range values {
		if v == "" {
			continue
		}
		validValues++
		length := len(v)
		totalLen += length

		if length < minLength {
			minLength = length
		}
		if length > maxLength {
			maxLength = length
		}
	}

	if validValues == 0 {
		return nil, nil, nil
	}

	minLen = &minLength
	maxLen = &maxLength
	avg := float64(totalLen) / float64(validValues)
	avgLen = &avg

	return
}

// GetSampleValues returns up to 5 sample non-empty values
func GetSampleValues(values []string, maxSamples int) []string {
	var samples []string
	seen := make(map[string]bool)

	for _, v := range values {
		if v == "" || seen[v] {
			continue
		}
		samples = append(samples, v)
		seen[v] = true

		if len(samples) >= maxSamples {
			break
		}
	}

	return samples
}

// CountUnique counts unique non-empty values
func CountUnique(values []string) int64 {
	unique := make(map[string]bool)
	for _, v := range values {
		if v != "" {
			unique[v] = true
		}
	}
	return int64(len(unique))
}

// FormatFloat formats a float pointer for display
func FormatFloat(f *float64) string {
	if f == nil {
		return "N/A"
	}
	return fmt.Sprintf("%.2f", *f)
}
