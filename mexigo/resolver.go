package main

import (
	"github.com/maride/mexico/mexigo/interpreter"
	"log"
	"math"
	"net"
	"strings"
)

const (
	// The fake base domain which classifies a domain name as a mexico command, rather than a "normal" domain name
	MexicoFakeDomain = "mexico.invalid."
)

// This is a wrapper function for net.LookupMX(), filtering for mexico records, and sorting the remaining by linenum
func LookupMX(basedomain string) []interpreter.Codeline {
	// Do the basic lookup
	rawMX, lookupErr := net.LookupMX(basedomain)
	if lookupErr != nil {
		// Encountered error while looking up basedomain - log and return
		log.Printf("Failed to resolve '%s': %s", basedomain, lookupErr.Error())
		return nil
	}

	// Filter results for the (fake) domain "mexico.invalid."
	var filteredMX []*net.MX

	// Iterate over all returned MX records
	for _, raw := range rawMX {
		// Check if it's a mexico MX record
		if strings.HasSuffix(raw.Host, MexicoFakeDomain) {
			// it is, add to filtered array
			filteredMX = append(filteredMX, raw)
		}
	}

	// Sort filtered results, based on the priority - or line number, in the words of this esolang :)
	var records []interpreter.Codeline
	var smallestPriority uint16 = math.MaxInt16
	smallestPriorityIndex := 0

	// Iterate over filteredMX and delete the record with the smallest priority until we don't have any more filteredMX
	for len(filteredMX) > 0 {
		// Iterate over the filteredMX to find the one with the smallest priority
		for i, f := range filteredMX {
			if f.Pref < smallestPriority {
				// Found entry with smaller index than the current one
				smallestPriority = f.Pref
				smallestPriorityIndex = i
			}
		}

		// Remove mexico fake domain suffix
		command := strings.TrimSuffix(filteredMX[smallestPriorityIndex].Host, "." + MexicoFakeDomain)

		// Replace '-' with space.
		// This reserves the process done by the compiler to transform this command + arg into a FQDN
		command = strings.Replace(command, "-", " ", 1)

		// Add hostname of "smallest" record to the records array, and delete it from filteredMX
		records = append(records, interpreter.Codeline{
			Linenumber: int(filteredMX[smallestPriorityIndex].Pref),
			Code:       command,
		})
		filteredMX = append(filteredMX[:smallestPriorityIndex], filteredMX[smallestPriorityIndex + 1:]...)
	}

	// Return filtered and sorted records
	return records
}
