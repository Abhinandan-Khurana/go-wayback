# go-wayback

`go-wayback` is a high-performance command-line tool written in Go that interacts with the Wayback Machine API to retrieve archived URLs and related data for a given website. It features concurrent processing, flexible output options, and robust error handling to efficiently explore historical snapshots of web content.

## Features

- **Concurrent Processing**: Fast retrieval through parallel processing of URLs
- **Multiple Output Modes**:
  - Wayback URLs Retrieval
  - Browsable Archive Links
  - Unique Subdomains Extraction
  - CSV Output with detailed metadata
- **Flexible Output Handling**:
  - Standard output (stdout) by default
  - Optional file output with custom naming
  - CSV format support
- **Advanced Options**:
  - Configurable concurrent connections
  - Request timeout control
  - Verbose logging
  - URL deduplication
- **Resource Efficient**:
  - Streaming processing
  - Controlled memory usage
  - Proper resource cleanup

## Direct Installation

go install -v github.com/Abhinandan-Khurana/go-wayback@latest

## Manual Installation

1. **Prerequisites**:

   - [Go](https://golang.org/doc/install) (version 1.16 or higher)

2. **Clone the Repository**:
   git clone <https://github.com/Abhinandan-Khurana/go-wayback.git>
   cd go-wayback

   ```

   ```

3. **Build the Executable**:

   ```
   go build -o go-wayback main.go
   ```

## Usage

./go-wayback [options]

### Options

- `-wayback-only`: Get only wayback URLs
- `-browsable`: Get wayback browsable links
- `-save-wayback-csv`: Output as CSV with URL, LENGTH, TIMESTAMP
- `-subdomain`: Extract and display unique subdomains
- `-v`: Enable verbose output
- `-o [file]`: Specify output file (optional, defaults to stdout)
- `-concurrent [n]`: Number of concurrent processors (default: 10)
- `-timeout [seconds]`: Request timeout in seconds (default: 30)
- `-h`: Display help information

### Output Behavior

- Default output is to stdout unless `-o` flag is specified
- CSV output includes headers: URL, LENGTH, TIMESTAMP
- Subdomain mode automatically deduplicates results
- Verbose mode provides additional processing information

## Examples

### Basic URL Retrieval

# Output to stdout

```
./go-wayback -wayback-only example.com
```

# Save to file

```
./go-wayback -o results.txt example.com
```

### Unique Subdomains

# Get unique subdomains

```
./go-wayback -subdomain example.com
```

# Save subdomains to file

```
./go-wayback -subdomain -o subdomains.txt example.com
```

### CSV Output with Metadata

# Default CSV format

```
./go-wayback -save-wayback-csv example.com
```

# Custom CSV file

```
./go-wayback -save-wayback-csv -o archive_data.csv example.com
```

### Browsable Archive Links

# Get browsable Wayback Machine links

```
./go-wayback -browsable example.com
```

### Advanced Usage

```
# Concurrent processing with timeout

./go-wayback -concurrent 20 -timeout 45 example.com
```

## Future Enhancements

- [ ] Date range filtering for archives
- [ ] Batch processing from input file
- [ ] Custom field selection
- [ ] Regular expression filtering
- [ ] Export formats (JSON, XML)
- [ ] Rate limiting controls
- [ ] Advanced filtering options
- [ ] Integration with other archive services

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## Acknowledgments

- [The Wayback Machine](https://web.archive.org) for providing access to archived web content
- [Go](https://golang.org) community for the excellent tooling and libraries

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

[Abhinandan-Khurana](https://github.com/Abhinandan-Khurana)
