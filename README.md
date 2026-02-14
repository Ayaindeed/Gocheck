# GoCheck

GoCheck is a fast Go-based CLI for profiling datasets and running basic data-quality checks (row counts, missing values, type inference, and summary statistics) on CSV/Parquet files.

## Why This Project?

I created gocheck as a learning project to explore Go's capabilities in building practical CLI tools. It's designed to provide quick insights into datasets without loading heavy tools like Python/pandas, making it useful for rapid data exploration and validation in CI/CD pipelines or local workflows.

## Features

- **Fast Dataset Profiling**: Quickly analyze large CSV and Parquet files
- **Row Counts**: Get instant row and column counts
- **Missing Value Analysis**: Identify missing data and calculate percentages
- **Type Inference**: Automatic detection of data types (numeric, string, date, boolean)
- **Summary Statistics**: 
  - Numeric columns: min, max, mean, median, standard deviation
  - String columns: length statistics and sample values
- **Multiple Output Formats**: Table, JSON, and CSV output options
- **Data Quality Summary**: Overall data quality assessment with warnings

## Installation

### Prerequisites
- Go 1.25 or higher

### Install from source

```bash
# Clone the repository
git clone https://github.com/Ayaindeed/gocheck.git
cd gocheck

# Download dependencies
go mod download

# Build the binary
go build -o gocheck.exe

# Optional: Add to PATH for global usage
```

## Quick Start

### Basic Usage

```bash
# Profile a CSV file
gocheck testdata/sample.csv

# Profile with verbose output
gocheck -v testdata/products.csv

# Output as JSON
gocheck -f json testdata/sample.csv

# Output as CSV
gocheck -f csv testdata/products.csv
```

### Example Output

```
Dataset Profile: sample.csv
Rows:    10
Columns: 5
Time:    2026-02-14 10:30:45

+-------------+---------+---------+-----------+--------+------------------------+
|   COLUMN    |  TYPE   | MISSING | MISSING % | UNIQUE |         STATS          |
+-------------+---------+---------+-----------+--------+------------------------+
| name        | string  | 0       | 0.0%      | 10     | Len: [3, 7]           |
|             |         |         |           |        | Avg: 5.4              |
|             |         |         |           |        | Ex: Alice, Bob        |
+-------------+---------+---------+-----------+--------+------------------------+
| age         | numeric | 1       | 10.0%     | 9      | Range: [26.00, 45.00] |
|             |         |         |           |        | Mean: 34.11           |
|             |         |         |           |        | Median: 33.00         |
+-------------+---------+---------+-----------+--------+------------------------+
| salary      | numeric | 1       | 10.0%     | 9      | Range: [60000, 105000]|
|             |         |         |           |        | Mean: 80556.19        |
+-------------+---------+---------+-----------+--------+------------------------+
| department  | string  | 1       | 10.0%     | 3      | Engineering, Marketing|
+-------------+---------+---------+-----------+--------+------------------------+
| active      | boolean | 1       | 10.0%     | 2      | true, false           |
+-------------+---------+---------+-----------+--------+------------------------+

Data Quality Summary:
Columns with missing values: 4 / 5
Total missing cells: 4 / 50 (8.00%)
Acceptable level of missing data
```

## Command Options

```
Usage:
  gocheck [file] [flags]

Flags:
  -f, --format string   Output format: table, json, or csv (default "table")
  -h, --help           Help for gocheck
  -v, --verbose        Verbose output
```

## Supported File Types 📁

- **CSV** (`.csv`) - Co

- **CSV** (`.csv`) - Comma-separated values with automatic delimiter detection
- **Parquet** (`.parquet`) - Apache Parquet columnar storage format

## Type Inferencey infers column types:

- **numeric**: Integer or floating-point numbers
- **string**: Text data
- **date**: Date and datetime values (multiple formats supported)
- **boolean**: True/false values (various representations)

Type inference uses an 80% threshold - if 80% of non-empty values match a type, the column is classified as that type.

## Statistics Calculated 📊

### Numeric Columns
- Minimum value
- Maximum value
- Mean (average)
- Median
- Standard deviation
- Missing value count and percentage
- Unique value count

### String Columns
- Minimum length
- Maximum length
- Average length
- Missing value count and percentage
- Unique value count
- Sample values

## Project Structure 📂

```
gocheck/
├── cmd/
│   └── root.go              # CLI command definitions
├── internal/
│   └── profiler/
│       ├── profiler.go      # Core profiling logic
│       ├── csv.go           # CSV file processing
│       ├── parquet.go       # Parquet file processing
│       └── output.go        # Output formatters
├── testdata/
│   ├── sample.csv           # Sample CSV for testing
│   └── products.csv         # Another sample dataset
├── main.go                  # Entry point
├── go.mod                   # Go module definition
└── README.md                # This file
```

## Development 🔧

### Running Te

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Benchmark
go test -bench=. ./...
```

### Building

```bash
# Build for current platform
go build -o gocheck.exe

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o gocheck

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o gocheck
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [tablewriter](https://github.com/olekukonko/tablewriter) - ASCII table formatting
- [parquet-go](https://github.com/parquet-go/parquet-go) - Parquet file support

## Performance Tips

- gocheck loads entire datasets into memory for fast processing
- For very large files (>1GB), consider using smaller samples
- Parquet files are generally faster to process than CSV

## Roadmap

Future enhancements planned:

- [ ] Support for more file formats (JSON, Excel, SQL)
- [ ] Outlier detection
- [ ] Correlation analysis
- [ ] Data validation rules
- [ ] Export reports to HTML/PDF
- [ ] Streaming mode for very large files
- [ ] Schema comparison between datasets

## Contributing

Contributions welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - feel free to use this project for learning and fun!

## Learning Resources

This project is great for learning:

- Go CLI development with Cobra
- File I/O and CSV parsing
- Data structure design
- Statistical calculations
- Working with Parquet files
- Table formatting and output generation

## Author

Built by Ayaindeed
