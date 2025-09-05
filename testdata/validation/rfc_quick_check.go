package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/airtrafik/jscal"
)

func main() {
	examplesDir := "../rfc8984/examples"
	
	files, err := filepath.Glob(filepath.Join(examplesDir, "*.json"))
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}
	
	fmt.Println("Testing RFC 8984 Examples:")
	fmt.Println("===========================")
	
	for _, file := range files {
		name := filepath.Base(file)
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("❌ %s: Error reading file: %v\n", name, err)
			continue
		}
		
		// Try generic parse first
		obj, err := jscal.Parse(data)
		if err != nil {
			fmt.Printf("❌ %s: Parse failed: %v\n", name, err)
			continue
		}
		
		// Try validation
		if err := obj.Validate(); err != nil {
			fmt.Printf("⚠️  %s: Parsed but validation failed: %v\n", name, err)
		} else {
			fmt.Printf("✅ %s: Successfully parsed and validated (%s)\n", name, obj.GetType())
		}
	}
}