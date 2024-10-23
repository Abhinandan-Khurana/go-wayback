package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Define flags
	waybackOnly := flag.Bool("wayback-only", false, "Get only wayback URLs")
	browsable := flag.Bool("browsable", false, "Get wayback browsable links to see the archive")
	saveCSV := flag.Bool("save-wayback-csv", false, "Output the CSV with URL,LENGTH,TIMESTAMP")
	subdomain := flag.Bool("subdomain", false, "Get unique subdomains from the Wayback URLs")
	uniqueUrls := flag.Bool("unique-urls", false, "Remove duplicate URLs from the output")
	verbose := flag.Bool("v", false, "Enable verbose output")
	outputFile := flag.String("o", "", "Specify the output file name")
	help := flag.Bool("h", false, "Display help")

	// Parse flags
	flag.Parse()

	// Display help if requested or if no URL is provided
	if *help || len(flag.Args()) == 0 {
		printHelp()
		os.Exit(0)
	}

	// Get the URL argument
	inputURL := flag.Arg(0)
	if inputURL == "" {
		fmt.Println("Error: URL is required")
		os.Exit(1)
	}

	// Determine the mode
	modeCount := 0
	if *waybackOnly {
		modeCount++
	}
	if *browsable {
		modeCount++
	}
	if *saveCSV {
		modeCount++
	}
	if *subdomain {
		modeCount++
	}
	if modeCount > 1 {
		fmt.Println("Error: Please specify only one mode at a time (-wayback-only, -browsable, -save-wayback-csv, or -subdomain)")
		os.Exit(1)
	}
	// Default to -save-wayback-csv if no mode is specified
	if modeCount == 0 {
		*saveCSV = true
	}

	// Sanitize the URL for use in filename
	sanitizedURL := sanitizeFilename(inputURL)

	// Set default output file names if not specified
	if *outputFile == "" {
		if *waybackOnly || *browsable || *subdomain {
			*outputFile = sanitizedURL + ".txt"
		} else if *saveCSV {
			*outputFile = sanitizedURL + "_waybackArchive.csv"
		}
	}

	if *verbose {
		fmt.Printf("Input URL: %s\n", inputURL)
		fmt.Printf("Mode: ")
		if *waybackOnly {
			fmt.Println("Wayback Only")
		} else if *browsable {
			fmt.Println("Browsable")
		} else if *saveCSV {
			fmt.Println("Save CSV")
		} else if *subdomain {
			fmt.Println("Subdomain")
		}
		fmt.Printf("Output File: %s\n", *outputFile)
	}

	// Construct the API URL
	escapedURL := url.QueryEscape("*." + inputURL + "/*")
	apiURL := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s&fl=original,length,timestamp", escapedURL)

	if *verbose {
		fmt.Printf("Fetching data from API URL: %s\n", apiURL)
	}

	// Fetch data from the API
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error fetching data from Wayback Machine API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching data: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}

	bodyString := string(bodyBytes)
	lines := strings.Split(bodyString, "\n")

	var output []string
	var subdomainSet map[string]bool

	if *subdomain {
		subdomainSet = make(map[string]bool)
	}

	// Process each line
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Each line has format: URL LENGTH TIMESTAMP
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		originalURL := fields[0]
		length := fields[1]
		timestamp := fields[2]

		if *subdomain {
			// Extract subdomain from originalURL
			parsedURL, err := url.Parse(originalURL)
			if err != nil {
				if *verbose {
					fmt.Printf("Error parsing URL %s: %v\n", originalURL, err)
				}
				continue
			}
			host := parsedURL.Host
			if host == "" {
				continue
			}
			subdomainSet[host] = true
		} else if *waybackOnly {
			// Get only wayback URLs
			output = append(output, originalURL)
		} else if *browsable {
			// Get URL wayback browsable links to see the archive
			waybackURL := fmt.Sprintf("https://web.archive.org/web/%s/%s", timestamp, originalURL)
			output = append(output, waybackURL)
		} else if *saveCSV {
			// Output the CSV with URL,LENGTH,TIMESTAMP
			output = append(output, fmt.Sprintf("%s,%s,%s", originalURL, length, timestamp))
		}
	}

	// Collect subdomains if in subdomain mode
	if *subdomain {
		for sub := range subdomainSet {
			output = append(output, sub)
		}
	}

	// Remove duplicates if unique-urls flag is set
	if *uniqueUrls {
		output = uniqueStrings(output)
		if *verbose {
			fmt.Printf("After removing duplicates, total items: %d\n", len(output))
		}
	}

	if *verbose {
		fmt.Printf("Total items collected: %d\n", len(output))
	}

	// Output handling
	if *saveCSV {
		// Write CSV file
		file, err := os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		writer.Write([]string{"URL", "LENGTH", "TIMESTAMP"})

		for _, line := range output {
			record := strings.Split(line, ",")
			writer.Write(record)
		}
		fmt.Printf("Output saved to %s\n", *outputFile)
	} else {
		// Output to file
		outputData := strings.Join(output, "\n")
		err := ioutil.WriteFile(*outputFile, []byte(outputData), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Output saved to %s\n", *outputFile)
	}
}

func printHelp() {
	fmt.Println("Usage: go-wayback [options] <URL>")
	fmt.Println("Options:")
	fmt.Println("  -wayback-only        Get only wayback URLs")
	fmt.Println("  -browsable           Get wayback browsable links to see the archive")
	fmt.Println("  -save-wayback-csv    Output the CSV with URL,LENGTH,TIMESTAMP")
	fmt.Println("  -subdomain           Get unique subdomains from the Wayback URLs")
	fmt.Println("  -unique-urls         Remove duplicate URLs from the output")
	fmt.Println("  -v                   Enable verbose output")
	fmt.Println("  -o [file]            Specify the output file name")
	fmt.Println("  -h, --help           Display help")
	fmt.Println("")
	fmt.Println("Notes:")
	fmt.Println("- If none of the mode flags are specified, the default mode is -save-wayback-csv.")
	fmt.Println("- In -wayback-only, -browsable, and -subdomain modes, output is saved to $URL.txt unless -o is specified.")
	fmt.Println("- In -save-wayback-csv mode, output is saved to $URL_waybackArchive.csv unless -o is specified.")
}

func sanitizeFilename(input string) string {
	// Remove protocol prefixes
	input = strings.ReplaceAll(input, "http://", "")
	input = strings.ReplaceAll(input, "https://", "")
	// Replace non-alphanumeric characters with underscores
	sanitized := ""
	for _, r := range input {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') {
			sanitized += string(r)
		} else {
			sanitized += "_"
		}
	}
	return sanitized
}

func uniqueStrings(input []string) []string {
	set := make(map[string]bool)
	var output []string
	for _, item := range input {
		if !set[item] {
			set[item] = true
			output = append(output, item)
		}
	}
	return output
}
