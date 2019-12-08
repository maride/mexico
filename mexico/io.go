package main

import (
	"flag"
	"fmt"
	"github.com/maride/mexico/mexico/compiler"
	"io/ioutil"
	"strings"
	"time"
)

var (
	inputFilePath *string
	outputFilePath *string
	baseDomain *string
)

// Registers flags required for input and output
func registerIOFlags() {
	inputFilePath = flag.String("input", "", "Name of the source code file to read")
	outputFilePath = flag.String("output", "", "Name of the zonefile to write")
	baseDomain = flag.String("baseDomain", "mexico.invalid", "The base domain to write the zonefile for")
}

// Reads the input file, splits it at newlines, and returns it as a string array
func readFile() ([]string, error) {
	// Read file
	fileBytes, readErr := ioutil.ReadFile(*inputFilePath)
	if readErr != nil {
		// Error reading the file, pass through
		return nil, readErr
	}

	// And split along newlines
	return strings.Split(string(fileBytes), "\n"), nil
}

// Writes the code lines into the format of a Zonefile
func writeZone(code []compiler.Codeline) error {
	var zone strings.Builder
	domain := *baseDomain

	// Check if we need to append a dot to the end of the domain name
	if (*baseDomain)[len(*baseDomain)-1] != '.' {
		// No, append a dot
		domain = *baseDomain + "."
	}

	// Generate values for further usage
	serial := time.Now().Format("2006010215") // YYYYMMDDHH

	// Write SOA record
	zone.WriteString(fmt.Sprintf("%s\tIN SOA\t%s mexico.%s (%s 1h 1h 1h 1h)\n", domain, domain, domain, serial))

	// Write every other record
	for _, c := range code {
		zone.WriteString(fmt.Sprintf("%s\tIN MX\t%d %s\n", domain, c.Linenumber, c.Code))
	}

	// And write built string to file
	return ioutil.WriteFile(*outputFilePath, []byte(zone.String()), 0644)
}