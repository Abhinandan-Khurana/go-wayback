package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	VERSION = "v2.0.1"
	AUTHOR  = "Abhinandan-Khurana"
)

type WaybackResult struct {
	URL       string    `json:"url" xml:"url"`
	Length    string    `json:"length" xml:"length"`
	Timestamp string    `json:"timestamp" xml:"timestamp"`
	Error     error     `json:"-" xml:"-"`
	Date      time.Time `json:"date" xml:"date"`
}

type Config struct {
	WaybackOnly  bool
	Browsable    bool
	SaveCSV      bool
	Subdomain    bool
	UniqueURLs   bool
	Verbose      bool
	OutputFile   string
	Concurrent   int
	Timeout      int
	StartDate    string
	EndDate      string
	OutputFormat string
	InputFile    string
	RegexFilter  string
	RateLimit    int
	MaxResults   int
}

type XMLResponse struct {
	XMLName xml.Name        `xml:"wayback"`
	Results []WaybackResult `xml:"results>result"`
	Count   int             `xml:"count"`
}

type VersionInfo struct {
	Version     string    `json:"version"`
	Author      string    `json:"author"`
	BuildDate   string    `json:"buildDate"`
	LastUpdated time.Time `json:"lastUpdated"`
}

func getVersionInfo() VersionInfo {
	return VersionInfo{
		Version:     VERSION,
		Author:      AUTHOR,
		BuildDate:   time.Now().Format(time.RFC3339),
		LastUpdated: time.Now(),
	}
}

func printVersion() {
	info := getVersionInfo()
	fmt.Printf("go-wayback %s\n", info.Version)
	fmt.Printf("Author: %s\n", info.Author)
	fmt.Printf("Build Date: %s\n", info.BuildDate)
}

// RateLimiter implements rate limiting for API requests
type RateLimiter struct {
	ticker *time.Ticker
	stop   chan bool
}

func newRateLimiter(requestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
		stop:   make(chan bool),
	}
}

func (r *RateLimiter) Wait() {
	<-r.ticker.C
}

func (r *RateLimiter) Stop() {
	r.ticker.Stop()
	close(r.stop)
}

func extractSubdomains(URL string) string {
	URL = strings.TrimPrefix(URL, "http://")
	URL = strings.TrimPrefix(URL, "https://")

	if idx := strings.Index(URL, "/"); idx != -1 {
		URL = URL[:idx]
	}

	if idx := strings.Index(URL, ":"); idx != -1 {
		URL = URL[:idx]
	}

	return strings.ToLower(URL)
}

func loadURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if url := strings.TrimSpace(scanner.Text()); url != "" {
			urls = append(urls, url)
		}
	}

	return urls, scanner.Err()
}

func processDateRange(startDate, endDate string) (time.Time, time.Time, error) {
	layout := "2006-01-02"
	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse(layout, startDate)
		if err != nil {
			return start, end, fmt.Errorf("invalid start date format: %v", err)
		}
	}

	if endDate != "" {
		end, err = time.Parse(layout, endDate)
		if err != nil {
			return start, end, fmt.Errorf("invalid end date format: %v", err)
		}
	} else {
		end = time.Now()
	}

	return start, end, nil
}

func matchesFilter(url string, regexPattern string) bool {
	if regexPattern == "" {
		return true
	}
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return true // If regex is invalid, don't filter
	}
	return regex.MatchString(url)
}

func processJSONFormat(bodyBytes []byte, config Config, writer io.Writer) error {
	lines := strings.Split(string(bodyBytes), "\n")
	var results []WaybackResult
	uniqueURLs := make(map[string]bool)
	count := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Apply regex filter
		if !matchesFilter(fields[0], config.RegexFilter) {
			continue
		}

		// Handle unique URLs
		if config.UniqueURLs {
			if uniqueURLs[fields[0]] {
				continue
			}
			uniqueURLs[fields[0]] = true
		}

		timestamp, _ := time.Parse("20060102150405", fields[2])
		result := WaybackResult{
			URL:       fields[0],
			Length:    fields[1],
			Timestamp: fields[2],
			Date:      timestamp,
		}

		results = append(results, result)
		count++

		if config.MaxResults > 0 && count >= config.MaxResults {
			break
		}
	}

	return json.NewEncoder(writer).Encode(map[string]interface{}{
		"results": results,
		"count":   len(results),
	})
}

func processXMLFormat(bodyBytes []byte, config Config, writer io.Writer) error {
	lines := strings.Split(string(bodyBytes), "\n")
	var results []WaybackResult
	uniqueURLs := make(map[string]bool)
	count := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Apply regex filter
		if !matchesFilter(fields[0], config.RegexFilter) {
			continue
		}

		// Handle unique URLs
		if config.UniqueURLs {
			if uniqueURLs[fields[0]] {
				continue
			}
			uniqueURLs[fields[0]] = true
		}

		timestamp, _ := time.Parse("20060102150405", fields[2])
		result := WaybackResult{
			URL:       fields[0],
			Length:    fields[1],
			Timestamp: fields[2],
			Date:      timestamp,
		}

		results = append(results, result)
		count++

		if config.MaxResults > 0 && count >= config.MaxResults {
			break
		}
	}

	xmlData := XMLResponse{
		Results: results,
		Count:   len(results),
	}

	// Write XML header
	fmt.Fprintf(writer, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")

	encoder := xml.NewEncoder(writer)
	encoder.Indent("", "  ")
	return encoder.Encode(xmlData)
}

func processCSVFormat(bodyBytes []byte, config Config, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	if err := csvWriter.Write([]string{"URL", "LENGTH", "TIMESTAMP", "DATE"}); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	lines := strings.Split(string(bodyBytes), "\n")
	uniqueURLs := make(map[string]bool)
	count := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Apply regex filter
		if !matchesFilter(fields[0], config.RegexFilter) {
			continue
		}

		// Handle unique URLs
		if config.UniqueURLs {
			if uniqueURLs[fields[0]] {
				continue
			}
			uniqueURLs[fields[0]] = true
		}

		timestamp, _ := time.Parse("20060102150405", fields[2])
		record := []string{
			fields[0],
			fields[1],
			fields[2],
			timestamp.Format(time.RFC3339),
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %v", err)
		}

		count++
		if config.MaxResults > 0 && count >= config.MaxResults {
			break
		}
	}

	return nil
}

func processTextFormat(bodyBytes []byte, config Config, writer io.Writer) error {
	lines := strings.Split(string(bodyBytes), "\n")
	uniqueURLs := make(map[string]bool)
	count := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Apply regex filter
		if !matchesFilter(fields[0], config.RegexFilter) {
			continue
		}

		// Handle unique URLs
		if config.UniqueURLs {
			if uniqueURLs[fields[0]] {
				continue
			}
			uniqueURLs[fields[0]] = true
		}

		outputURL := fields[0]
		if config.Browsable {
			outputURL = fmt.Sprintf("https://web.archive.org/web/%s/%s", fields[2], fields[0])
		}

		fmt.Fprintln(writer, outputURL)

		count++
		if config.MaxResults > 0 && count >= config.MaxResults {
			break
		}
	}

	if config.Verbose {
		fmt.Fprintf(os.Stderr, "Total URLs processed: %d\n", count)
	}

	return nil
}

func processURL(inputURL string, config Config, rateLimiter *RateLimiter) error {
	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "http://" + inputURL
	}

	escapedURL := url.QueryEscape("*." + inputURL + "/*")
	apiURL := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s&fl=original,length,timestamp", escapedURL)

	if config.Verbose {
		fmt.Fprintf(os.Stderr, "Fetching data from: %s\n", apiURL)
	}

	rateLimiter.Wait()

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Handle subdomain mode separately
	if config.Subdomain {
		subdomains := make(map[string]bool)
		lines := strings.Split(string(bodyBytes), "\n")

		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) < 1 {
				continue
			}

			subdomain := extractSubdomains(fields[0])
			if subdomain != "" {
				subdomains[subdomain] = true
			}
		}

		var uniqueSubdomains []string
		for subdomain := range subdomains {
			uniqueSubdomains = append(uniqueSubdomains, subdomain)
		}
		sort.Strings(uniqueSubdomains)

		var writer io.Writer = os.Stdout
		if config.OutputFile != "" {
			file, err := os.Create(config.OutputFile)
			if err != nil {
				return fmt.Errorf("error creating output file: %v", err)
			}
			defer file.Close()
			writer = file
		}

		for _, subdomain := range uniqueSubdomains {
			fmt.Fprintln(writer, subdomain)
		}

		if config.Verbose {
			fmt.Fprintf(os.Stderr, "Total unique subdomains found: %d\n", len(uniqueSubdomains))
		}

		return nil
	}

	// Process the response based on format
	var writer io.Writer = os.Stdout
	if config.OutputFile != "" {
		file, err := os.Create(config.OutputFile)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}
		defer file.Close()
		writer = file
	}

	switch strings.ToLower(config.OutputFormat) {
	case "json":
		return processJSONFormat(bodyBytes, config, writer)
	case "xml":
		return processXMLFormat(bodyBytes, config, writer)
	case "csv":
		return processCSVFormat(bodyBytes, config, writer)
	default:
		return processTextFormat(bodyBytes, config, writer)
	}
}

func main() {
	config := Config{}

	flag.BoolVar(&config.WaybackOnly, "wayback-only", false, "Get only wayback URLs")
	flag.BoolVar(&config.Browsable, "browsable", false, "Get wayback browsable links")
	flag.BoolVar(&config.SaveCSV, "save-wayback-csv", false, "Output as CSV")
	flag.BoolVar(&config.Subdomain, "subdomain", false, "Get unique subdomains")
	flag.BoolVar(&config.UniqueURLs, "unique-urls", false, "Remove duplicate URLs")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.StringVar(&config.OutputFile, "o", "", "Output file (optional)")
	flag.IntVar(&config.Concurrent, "concurrent", 10, "Number of concurrent processors")
	flag.IntVar(&config.Timeout, "timeout", 30, "Request timeout in seconds")
	flag.StringVar(&config.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	flag.StringVar(&config.EndDate, "end-date", "", "End date (YYYY-MM-DD)")
	flag.StringVar(&config.OutputFormat, "format", "text", "Output format (text/json/xml/csv)")
	flag.StringVar(&config.InputFile, "input-file", "", "File containing URLs to process")
	flag.StringVar(&config.RegexFilter, "filter", "", "Regex pattern to filter URLs")
	flag.IntVar(&config.RateLimit, "rate-limit", 10, "Maximum requests per second")
	flag.IntVar(&config.MaxResults, "max-results", 0, "Maximum number of results (0 for unlimited)")

	version := flag.Bool("version", false, "Display version information")
	help := flag.Bool("h", false, "Display help")

	flag.Parse()

	if *version {
		printVersion()
		return
	}

	if *help || (flag.NArg() == 0 && config.InputFile == "") {
		printHelp()
		return
	}

	rateLimiter := newRateLimiter(config.RateLimit)
	defer rateLimiter.Stop()

	var urls []string
	if config.InputFile != "" {
		var err error
		urls, err = loadURLsFromFile(config.InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading URLs from file: %v\n", err)
			os.Exit(1)
		}
	} else {
		urls = []string{flag.Arg(0)}
	}

	for _, inputURL := range urls {
		if err := processURL(inputURL, config, rateLimiter); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", inputURL, err)
		}
	}
}

func printHelp() {
	fmt.Println("go-wayback " + VERSION)
	fmt.Println("\nUsage: wayback [options] <URL>")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}
