# go-wayback

`go-wayback` is a command-line tool written in Go that interacts with the Wayback Machine API to retrieve archived URLs and related data for a given website. It offers multiple modes to fetch different types of information and outputs the results to files or CSV, enhancing your ability to explore historical snapshots of web content.

## Features

- **Wayback URLs Retrieval**: Fetch only the original URLs archived by the Wayback Machine.
- **Browsable Archive Links**: Generate browsable links to view archived versions of URLs directly in your browser.
- **CSV Output**: Save detailed information including URL, content length, and timestamp in a CSV file.
- **Flexible Output Options**: Specify output filenames or use default naming conventions based on the sanitized URL.
- **User-Friendly Interface**: Simple and descriptive command-line flags for an optimized user experience.
- **Help Menu**: Built-in help command to guide users through usage and options.

## Direct installation

```bash
go install -v github.com/Abhinandan-Khurana/go-wayback@latest
```

## Installation

1. **Prerequisites**:

   - [Go](https://golang.org/doc/install) (version 1.16 or higher)

2. **Clone the Repository**:

   ```bash
   git clone https://github.com/Abhinandan-Khurana/go-wayback.git
   cd go-wayback
   ```

3. **Build the Executable**:

   ```bash
   go build -o go-wayback main.go
   ```

   This will generate a `go-wayback` executable in your current directory.

## Usage

```bash
./go-wayback [options] <URL>
```

**Options**:

- `-wayback-only`: Get only wayback URLs.
- `-subdomain`: Get unique subdomains from the Wayback URLs.
- `-browsable`: Get wayback browsable links to see the archive.
- `-unique-urls`: Remove duplicate URLs from the output.
- `-save-wayback-csv`: Output the CSV with URL, LENGTH, TIMESTAMP.
- `-v`: Enable verbose output
- `-o [file]`: Specify the output file name.
- `-h`, `--help`: Display help.

**Notes**:

- If none of the mode flags are specified, the default mode is `-save-wayback-csv`.
- In `-wayback-only`, `-browsable` and `-subdomain` modes, output is saved to `$URL.txt` unless `-o` is specified.
- In `-save-wayback-csv` mode, output is saved to `$URL_waybackArchive.csv` unless `-o` is specified.
- Only one mode flag (`-wayback-only`, `-browsable`, or `-save-wayback-csv`) can be used at a time.

## Examples

### Get Only Wayback URLs

Retrieve all archived URLs for `example.com` and save them to the default output file (`example_com.txt`):

```bash
./go-wayback -wayback-only example.com
```

Specify a custom output file:

```bash
./go-wayback -wayback-only -o wayback_urls.txt example.com
```

### Get Browsable Archive Links

Generate browsable links to view archived versions of `example.com`:

```bash
./go-wayback -browsable example.com
```

Save to a custom file:

```bash
./go-wayback -browsable -o browsable_links.txt example.com
```

### Output CSV with URL, LENGTH, TIMESTAMP

Fetch detailed archive data and save it as a CSV file (default mode):

```bash
./go-wayback example.com
```

Specify a custom CSV output file:

```bash
./go-wayback -save-wayback-csv -o archive_data.csv example.com
```

### Display Help Menu

```bash
./go-wayback -h
```

## Acknowledgments

- [The Wayback Machine](https://web.archive.org) for providing access to archived web content.
- [Go](https://golang.org) language and its community for making powerful tools accessible.

## Potential Future Enhancements

- **Date Range Filtering**: Allow users to specify a date range for the archives.
- **Result Limiting**: Add an option to limit the number of results returned.
- **Field Selection**: Enable users to choose specific fields to retrieve from the API.
- add a feature to take domain input from a file.

Feel free to contribute to these features or suggest new ones!
