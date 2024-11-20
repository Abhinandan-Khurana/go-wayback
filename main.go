package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type WaybackResult struct {
	URL       string
	Length    string
	Timestamp string
	Error     error
}

type Config struct {
	WaybackOnly bool
	Browsable   bool
	SaveCSV     bool
	Subdomain   bool
	UniqueURLs  bool
	Verbose     bool
	OutputFile  string
	Concurrent  int
	Timeout     int
}

func fetchWaybackData(apiURL string, timeout int) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func processURLs(lines []string, config Config, results chan<- WaybackResult) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Concurrent)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			fields := strings.Fields(line)
			if len(fields) < 3 {
				return
			}

			result := WaybackResult{
				URL:       fields[0],
				Length:    fields[1],
				Timestamp: fields[2],
			}

			if config.Subdomain {
				parsedURL, err := url.Parse(result.URL)
				if err != nil {
					result.Error = err
				} else {
					result.URL = parsedURL.Host
				}
			} else if config.Browsable {
				result.URL = fmt.Sprintf("https://web.archive.org/web/%s/%s", result.Timestamp, result.URL)
			}

			results <- result
		}(line)
	}

	wg.Wait()
	close(results)
}

func writeOutput(writer io.Writer, results []WaybackResult, config Config) error {
	if config.SaveCSV {
		csvWriter := csv.NewWriter(writer)
		defer csvWriter.Flush()

		// Write header
		if err := csvWriter.Write([]string{"URL", "LENGTH", "TIMESTAMP"}); err != nil {
			return fmt.Errorf("error writing CSV header: %v", err)
		}

		for _, result := range results {
			if err := csvWriter.Write([]string{result.URL, result.Length, result.Timestamp}); err != nil {
				return fmt.Errorf("error writing CSV record: %v", err)
			}
		}
	} else {
		for _, result := range results {
			fmt.Fprintln(writer, result.URL)
		}
	}
	return nil
}

func main() {
	config := Config{}

	// Define flags
	flag.BoolVar(&config.WaybackOnly, "wayback-only", false, "Get only wayback URLs")
	flag.BoolVar(&config.Browsable, "browsable", false, "Get wayback browsable links")
	flag.BoolVar(&config.SaveCSV, "save-wayback-csv", false, "Output as CSV")
	flag.BoolVar(&config.Subdomain, "subdomain", false, "Get unique subdomains")
	flag.BoolVar(&config.UniqueURLs, "unique-urls", true, "Remove duplicate URLs")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.StringVar(&config.OutputFile, "o", "", "Output file (optional)")
	flag.IntVar(&config.Concurrent, "concurrent", 10, "Number of concurrent processors")
	flag.IntVar(&config.Timeout, "timeout", 30, "Request timeout in seconds")

	help := flag.Bool("h", false, "Display help")
	flag.Parse()

	if *help || len(flag.Args()) == 0 {
		printHelp()
		return
	}

	inputURL := flag.Arg(0)
	if inputURL == "" {
		fmt.Println("Error: URL is required")
		os.Exit(1)
	}

	// Construct API URL
	escapedURL := url.QueryEscape("*." + inputURL + "/*")
	apiURL := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s&fl=original,length,timestamp", escapedURL)

	if config.Verbose {
		fmt.Fprintf(os.Stderr, "Fetching data from: %s\n", apiURL)
	}

	// Fetch and process data
	bodyBytes, err := fetchWaybackData(apiURL, config.Timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(bodyBytes), "\n")
	results := make(chan WaybackResult, len(lines))

	go processURLs(lines, config, results)

	// Collect results
	var processedResults []WaybackResult
	uniqueURLs := make(map[string]bool)

	for result := range results {
		if result.Error != nil && config.Verbose {
			fmt.Fprintf(os.Stderr, "Error processing URL: %v\n", result.Error)
			continue
		}

		if config.UniqueURLs {
			if !uniqueURLs[result.URL] {
				uniqueURLs[result.URL] = true
				processedResults = append(processedResults, result)
			}
		} else {
			processedResults = append(processedResults, result)
		}
	}

	// Handle output
	var writer io.Writer = os.Stdout
	var file *os.File

	if config.OutputFile != "" {
		file, err = os.Create(config.OutputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	if err := writeOutput(writer, processedResults, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Fprintf(os.Stderr, "Total items processed: %d\n", len(processedResults))
		if config.OutputFile != "" {
			fmt.Fprintf(os.Stderr, "Output saved to: %s\n", config.OutputFile)
		}
	}
}

func printHelp() {
	fmt.Println("Usage: wayback [options] <URL>")
	fmt.Println("\nOptions:")
	fmt.Println("  -wayback-only       Get only wayback URLs")
	fmt.Println("  -browsable          Get wayback browsable links")
	fmt.Println("  -save-wayback-csv   Output as CSV")
	fmt.Println("  -subdomain          Get unique subdomains")
	//	fmt.Println("  -unique-urls        Remove duplicate URLs")
	fmt.Println("  -v                  Enable verbose output")
	fmt.Println("  -o [file]           Output file (optional)")
	fmt.Println("  -concurrent [n]      Number of concurrent processors (default: 10)")
	fmt.Println("  -timeout [seconds]   Request timeout in seconds (default: 30)")
	fmt.Println("  -h                  Display this help")
}
