# go-wayback

`go-wayback` is a high-performance command-line tool written in Go that interacts with the Wayback Machine API to retrieve archived URLs and related data for a given website. Version 1.0.4 introduces concurrent processing, advanced filtering options, and multiple output formats to efficiently explore historical snapshots of web content.

[![Go Report Card](https://goreportcard.com/badge/github.com/Abhinandan-Khurana/go-wayback)](https://goreportcard.com/report/github.com/Abhinandan-Khurana/go-wayback)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Version](https://img.shields.io/badge/version-2.0.0-blue)

## Features

### Core Functionality

- **Wayback URLs Retrieval**: Fetch original URLs archived by the Wayback Machine
- **Browsable Archive Links**: Generate direct links to archived versions
- **Subdomain Extraction**: Identify and list unique subdomains
- **Multiple Output Formats**: Support for plain text, CSV, JSON, and XML

### Advanced Features

- **Concurrent Processing**: Fast retrieval through parallel processing
- **Date Range Filtering**: Filter archives by specific time periods
- **Rate Limiting**: Control request rates to prevent API throttling
- **Regex Filtering**: Filter URLs using regular expressions
- **Result Limiting**: Control the number of results returned
- **Batch Processing**: Process multiple URLs from an input file

## Installation

### Direct Installation

```
go install -v github.com/Abhinandan-Khurana/go-wayback@latest
```

### Manual Installation

1. **Prerequisites**:

   - Go 1.16 or higher

2. **Clone and Build**:

   ```
   git clone https://github.com/Abhinandan-Khurana/go-wayback.git
   cd go-wayback
   go build -o go-wayback main.go
   ```

## Usage

./go-wayback [options]

### Basic Options

- `-wayback-only`: Get only wayback URLs
- `-browsable`: Get wayback browsable links
- `-subdomain`: Extract unique subdomains
- `-unique-urls`: Remove duplicate URLs
- `-save-wayback-csv`: Output as CSV
- `-o [file]`: Specify output file (optional, defaults to stdout)
- `-v`: Enable verbose output
- `-h`: Display help information
- `-version`: Show version information

### Advanced Options

- `-start-date`: Start date for filtering (YYYY-MM-DD)
- `-end-date`: End date for filtering (YYYY-MM-DD)
- `-format`: Output format (text/json/xml/csv)
- `-input-file`: File containing URLs to process
- `-filter`: Regex pattern to filter URLs
- `-rate-limit`: Maximum requests per second (default: 10)
- `-max-results`: Maximum number of results (0 for unlimited)
- `-concurrent`: Number of concurrent processors (default: 10)
- `-timeout`: Request timeout in seconds (default: 30)

## Examples

### Basic Usage

# Get all archived URLs (output to stdout)

```
./go-wayback example.com
```

# Save results to a file

```
./go-wayback -o results.txt example.com
```

# Extract unique subdomains

```
./go-wayback -subdomain example.com
```

### Advanced Usage

# Date range filtering with JSON output

```
./go-wayback -start-date 2020-01-01 -end-date 2023-12-31 -format json example.com
```

# Process multiple URLs from file with rate limiting

```
./go-wayback -input-file urls.txt -rate-limit 5 -concurrent 20
```

# Filter URLs using regex and limit results

```
./go-wayback -filter ".*\.pdf$" -max-results 100 example.com
```

# Get browsable links with custom timeout

```
./go-wayback -browsable -timeout 45 example.com
```

### Output Format Examples

# CSV output with metadata

```
./go-wayback -save-wayback-csv -o archive_data.csv example.com
```

# JSON output

```
./go-wayback -format json -o data.json example.com
```

# XML output

```
./go-wayback -format xml -o data.xml example.com
```

## Output Formats

### Plain Text (default)

- One URL per line
- Simple and grep-friendly

### CSV

- Headers: URL, LENGTH, TIMESTAMP
- Includes metadata for each archive

### JSON

```
{
  "results": [
    {
      "url": "http://example.com",
      "length": "12345",
      "timestamp": "20230101120000",
      "date": "2023-01-01T12:00:00Z"
    }
  ]
}
```

### XML

````
<results>
  <result>
    <url>http://example.com</url>
    <length>12345</length>
    <timestamp>20230101120000</timestamp>
    <date>2023-01-01T12:00:00Z</date>
  </result>
</results>```

## Performance Considerations

- Use `-concurrent` to adjust parallel processing based on your system capabilities
- Use `-rate-limit` to prevent API throttling
- For large datasets, consider using `-max-results` to limit output
- Enable `-v` for progress monitoring during long operations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [The Wayback Machine](https://web.archive.org) for providing access to archived web content
- [Go](https://golang.org) community for excellent tooling and libraries

## Author

[Abhinandan-Khurana](https://github.com/Abhinandan-Khurana)
````
