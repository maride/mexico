package compiler

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var (
	// List of commands which can be translated into FQDNs in one step.
	passThroughCommands = []string{
		"left",
		"right",
		"pusht",
		// push is missing here, because it is not a simple "pass-through" command; we need to process it further.
		"pop",
		"dup",
		"del",
		"eq",
		"not",
		"gt",
		"lt",
		"add",
		"sub",
		"mult",
		"div",
		"mod",
		"read",
		"print",
		"jmp",
		"jmpc",
	}
)

const (
	FakeFQDN = "mexico.invalid"
)

// This is the compile function. As the name suggests, it compiles the source code handed over.
// For this task, it takes three steps:
// - Iterate over the source code and clean it (remove comments, remove empty lines, remove surrounding spaces)
// - Iterate over the source code, number each line and build up a label lookup table (mapping labels to line numbers)
// - Iterate over the source code and translate the instructions to valid MX records
func Compile(lines []string, domain string) ([]Codeline, error) {
	numberedCode, labelLookupTable := numberLines(cleanCode(lines))
	return translateLines(numberedCode, labelLookupTable, domain)
}

// Cleans the lines in the string array: remove comments, remove empty lines, remove spaces
func cleanCode(lines []string) []string {
	var cleanLines []string

	// Iterate over all code lines, and clean them
	for _, l := range lines {
		// Remove surrounding spaces and tabs
		l = strings.Trim(l, " \t")

		// Check if line is a comment
		if strings.HasPrefix(l, "#") || strings.HasPrefix(l, ";") || strings.HasPrefix(l, "//") {
			// It is a comment line, ignore
			continue
		}

		// Check if line is empty
		if len(l) == 0 {
			// Empty line, ignore.
			continue
		}

		// If we reach this point, the line is ready to be added to the list of cleaned lines
		cleanLines = append(cleanLines, l)
	}

	// Return cleaned lines
	return cleanLines
}

// Adds line numbers to the code lines, and builds up a label lookup table, mapping label to line numbers
func numberLines(lines []string) ([]Codeline, map[string]int) {
	var code []Codeline
	var labelLookup = make(map[string]int)
	linenumber := 0

	// Iterate over all (string) lines and convert them to codelines
	for _, l := range lines {
		// Check if line is a label - defined by ':' at the end
		if l[len(l)-1] == ':' {
			// It's a label. Write it into the lookup table
			name := l[:len(l)-1]
			labelLookup[name] = linenumber

			// And skip further execution, to avoid raising the line number or appending this line to the code array
			continue
		}

		// Append codeline to the array
		code = append(code, Codeline{
			Linenumber: linenumber,
			Code: l,
		})

		// Raise line number
		linenumber++
	}

	// Returns the code lines and the label lookup table
	return code, labelLookup
}

// Translates the code into FQDNs to be further used for MX records, resolving  labels to their line numbers
func translateLines(code []Codeline, labelLookup map[string]int, domain string) ([]Codeline, error) {
	// Iterate over all commands
	for i := 0; i < len(code); i++ {
		// Check if command is part of "pass-through" commands.
		found := false
		for _, c := range passThroughCommands {
			if code[i].Code == c {
				// It is! We can simply translate it to a FQDN without further processing required.
				code[i].Code = fmt.Sprintf("%s.%s.", c, FakeFQDN)
				found = true
				break
			}
		}

		// Check if command was a pass-through command
		if found {
			// It was, go on to the next command
			continue
		}

		// Check if it is the push command
		if code[i].Code[:4] == "push" {
			// It is. Either we have a constant here, or a label name - let's check.
			maybeLabelMaybeConstant := code[i].Code[5:]
			maybeLinenumberMaybeNot := getNumberForLabel(maybeLabelMaybeConstant, labelLookup)

			// if maybeLinenumberMaybeNot is not -1, there is a label with this name :)
			if maybeLinenumberMaybeNot > -1 {
				// Found linenumber for that label - write it into code
				code[i].Code = fmt.Sprintf("push-%d.%s.", maybeLinenumberMaybeNot, FakeFQDN)
				continue
			}

			// If we are at this point, we didn't find a label with that name. Try to parse it as a number
			maybeValueMaybeNot, atoiErr := strconv.Atoi(maybeLabelMaybeConstant)

			// Throw error if parsing didn't work
			if atoiErr != nil {
				errorString := fmt.Sprintf("Not a label or a integer constant: '%d'. %s", maybeLinenumberMaybeNot, atoiErr.Error())
				return nil, errors.New(errorString)
			}

			// If we are at this point, we were able to parse the mystery value as an integer. Yay!
			code[i].Code = fmt.Sprintf("push-%d.%s.", maybeValueMaybeNot, FakeFQDN)
			continue
		}

		// uh, if we reach this point, it's not a valid command - return
		return nil, errors.New(fmt.Sprintf("Not a valid command: %s", code[i].Code))
	}

	// And return our constructed source code
	return code, nil
}

// Returns the line number for the given label, or -1 if no such label was found
func getNumberForLabel(label string, labelLookup map[string]int) int {
	// iterate over all labels
	for l, number := range labelLookup {
		if l == label {
			// Found, return number
			return number
		}
	}

	// None found, return -1
	return -1
}