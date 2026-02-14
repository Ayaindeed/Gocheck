package profiler

import (
	"testing"
)

func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "numeric values",
			values:   []string{"1", "2", "3.5", "4.2"},
			expected: "numeric",
		},
		{
			name:     "string values",
			values:   []string{"hello", "world", "test"},
			expected: "string",
		},
		{
			name:     "boolean values",
			values:   []string{"true", "false", "True", "False"},
			expected: "boolean",
		},
		{
			name:     "date values",
			values:   []string{"2024-01-01", "2024-02-15", "2024-03-20"},
			expected: "date",
		},
		{
			name:     "empty values",
			values:   []string{"", "", ""},
			expected: "string",
		},
		{
			name:     "mixed values",
			values:   []string{"1", "2", "hello", "world"},
			expected: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InferType(tt.values)
			if result != tt.expected {
				t.Errorf("InferType(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestCalculateNumericStats(t *testing.T) {
	values := []string{"1", "2", "3", "4", "5"}
	min, max, mean, median, stddev := CalculateNumericStats(values)

	if min == nil || *min != 1.0 {
		t.Errorf("Min = %v, want 1.0", min)
	}
	if max == nil || *max != 5.0 {
		t.Errorf("Max = %v, want 5.0", max)
	}
	if mean == nil || *mean != 3.0 {
		t.Errorf("Mean = %v, want 3.0", mean)
	}
	if median == nil || *median != 3.0 {
		t.Errorf("Median = %v, want 3.0", median)
	}
	if stddev == nil {
		t.Error("StdDev should not be nil")
	}
}

func TestCalculateStringStats(t *testing.T) {
	values := []string{"a", "bb", "ccc", "dddd"}
	minLen, maxLen, avgLen := CalculateStringStats(values)

	if minLen == nil || *minLen != 1 {
		t.Errorf("MinLen = %v, want 1", minLen)
	}
	if maxLen == nil || *maxLen != 4 {
		t.Errorf("MaxLen = %v, want 4", maxLen)
	}
	if avgLen == nil || *avgLen != 2.5 {
		t.Errorf("AvgLen = %v, want 2.5", avgLen)
	}
}

func TestCountUnique(t *testing.T) {
	values := []string{"a", "b", "a", "c", "b", ""}
	unique := CountUnique(values)

	if unique != 3 {
		t.Errorf("CountUnique = %v, want 3", unique)
	}
}

func TestGetSampleValues(t *testing.T) {
	values := []string{"a", "b", "c", "d", "e", "f", "g"}
	samples := GetSampleValues(values, 3)

	if len(samples) != 3 {
		t.Errorf("GetSampleValues length = %v, want 3", len(samples))
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, s := range samples {
		if seen[s] {
			t.Errorf("Found duplicate sample: %s", s)
		}
		seen[s] = true
	}
}

func BenchmarkInferType(b *testing.B) {
	values := []string{"1", "2", "3.5", "4.2", "5.8"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InferType(values)
	}
}

func BenchmarkCalculateNumericStats(b *testing.B) {
	values := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateNumericStats(values)
	}
}
