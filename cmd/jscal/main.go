// Command jscal provides CLI tools for working with JSCalendar data.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/airtrafik/jscal"
	"github.com/airtrafik/jscal/convert/ical"
)

const version = "0.2.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "convert":
		handleConvert(args)
	case "validate":
		handleValidate(args)
	case "format":
		handleFormat(args)
	case "version":
		fmt.Printf("jscal version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`jscal v%s - JSCalendar CLI tool

USAGE:
    jscal <command> [options] [arguments]

COMMANDS:
    convert     Convert between calendar formats
    validate    Validate JSCalendar files
    format      Pretty-print JSCalendar files
    version     Show version information
    help        Show this help message

CONVERT USAGE:
    jscal convert <input> <output>           Auto-detect format and convert
    jscal convert -f ical <input> <output>   Convert from iCalendar to JSCalendar
    jscal convert -t ical <input> <output>   Convert JSCalendar to iCalendar

VALIDATE USAGE:
    jscal validate <file>...                 Validate JSCalendar files

FORMAT USAGE:
    jscal format <file>...                   Pretty-print JSCalendar files

EXAMPLES:
    jscal convert calendar.ics calendar.json
    jscal convert -t ical event.json event.ics
    jscal validate events.json
    jscal format messy.json

`, version)
}

func handleConvert(args []string) {
	var fromFormat, toFormat string
	var inputFile, outputFile string

	// Parse flags
	i := 0
	for i < len(args) {
		arg := args[i]
		switch arg {
		case "-f", "--from":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: %s requires a value\n", arg)
				os.Exit(1)
			}
			fromFormat = args[i+1]
			i += 2
		case "-t", "--to":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: %s requires a value\n", arg)
				os.Exit(1)
			}
			toFormat = args[i+1]
			i += 2
		default:
			if inputFile == "" {
				inputFile = arg
			} else if outputFile == "" {
				outputFile = arg
			} else {
				fmt.Fprintf(os.Stderr, "Error: unexpected argument %s\n", arg)
				os.Exit(1)
			}
			i++
		}
	}

	if inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: input file is required\n")
		os.Exit(1)
	}

	if outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: output file is required\n")
		os.Exit(1)
	}

	// Read input file
	inputData, err := readFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Auto-detect formats if not specified
	if fromFormat == "" {
		fromFormat = detectFormat(inputData, filepath.Ext(inputFile))
	}
	if toFormat == "" {
		toFormat = detectFormat(nil, filepath.Ext(outputFile))
	}

	// Convert
	outputData, err := convert(inputData, fromFormat, toFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting: %v\n", err)
		os.Exit(1)
	}

	// Write output file
	if err := writeFile(outputFile, outputData); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}

func handleValidate(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: at least one file is required\n")
		os.Exit(1)
	}

	var hasErrors bool
	for _, filename := range args {
		if err := validateFile(filename); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %s: %v\n", filename, err)
			hasErrors = true
		} else {
			fmt.Printf("✅ %s: valid\n", filename)
		}
	}

	if hasErrors {
		os.Exit(1)
	}
}

func handleFormat(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: at least one file is required\n")
		os.Exit(1)
	}

	for _, filename := range args {
		if err := formatFile(filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting %s: %v\n", filename, err)
			os.Exit(1)
		}
	}
}

func convert(inputData []byte, fromFormat, toFormat string) ([]byte, error) {
	// First, convert to JSCalendar if needed
	var events []*jscal.Event
	var err error

	switch strings.ToLower(fromFormat) {
	case "ical", "icalendar", "ics":
		converter := ical.New()
		events, err = converter.ParseAll(inputData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse iCalendar: %w", err)
		}
	case "json", "jscal", "jscalendar":
		// Try to parse as single event first
		if event, parseErr := jscal.ParseEvent(inputData); parseErr == nil {
			events = []*jscal.Event{event}
		} else {
			// Try as array of events
			events, err = jscal.ParseAllEvents(inputData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse JSCalendar: %w", err)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported input format: %s", fromFormat)
	}

	// Convert to target format
	switch strings.ToLower(toFormat) {
	case "ical", "icalendar", "ics":
		converter := ical.New()
		return converter.FormatAll(events)
	case "json", "jscal", "jscalendar":
		if len(events) == 1 {
			return events[0].PrettyJSON()
		} else {
			return json.MarshalIndent(events, "", "  ")
		}
	default:
		return nil, fmt.Errorf("unsupported output format: %s", toFormat)
	}
}

func detectFormat(data []byte, fileExt string) string {
	// Try file extension first
	switch strings.ToLower(fileExt) {
	case ".ics", ".ical":
		return "ical"
	case ".json":
		return "json"
	}

	// Try content detection if we have data
	if len(data) > 0 {
		dataStr := string(data)
		trimmed := strings.TrimSpace(dataStr)

		// Check for iCalendar
		if strings.HasPrefix(trimmed, "BEGIN:VCALENDAR") || strings.Contains(dataStr, "BEGIN:VEVENT") {
			return "ical"
		}

		// Check for JSON
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			return "json"
		}
	}

	// Default to JSON
	return "json"
}

func validateFile(filename string) error {
	data, err := readFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Try to parse as single event first
	if event, err := jscal.ParseEvent(data); err == nil {
		return event.Validate()
	}

	// Try as array of events
	events, err := jscal.ParseAllEvents(data)
	if err != nil {
		return fmt.Errorf("failed to parse JSCalendar: %w", err)
	}

	// Validate each event
	for i, event := range events {
		if err := event.Validate(); err != nil {
			return fmt.Errorf("event %d is invalid: %w", i, err)
		}
	}

	return nil
}

func formatFile(filename string) error {
	data, err := readFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Try to parse as single event first
	if event, err := jscal.ParseEvent(data); err == nil {
		formatted, err := event.PrettyJSON()
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
		fmt.Print(string(formatted))
		return nil
	}

	// Try as array of events
	events, err := jscal.ParseAllEvents(data)
	if err != nil {
		return fmt.Errorf("failed to parse JSCalendar: %w", err)
	}

	formatted, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	fmt.Print(string(formatted))
	return nil
}

func readFile(filename string) ([]byte, error) {
	if filename == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(filename)
}

func writeFile(filename string, data []byte) error {
	if filename == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
