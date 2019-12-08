package main

import (
	"flag"
	"github.com/maride/mexico/mexico/compiler"
	"log"
)

func main() {
	// Register flags
	registerIOFlags()
	flag.Parse()

	// Read file
	fileContent, readErr := readFile()
	handleErr(readErr)

	// Parse and compile lines
	code, compileErr := compiler.Compile(fileContent, *baseDomain)
	handleErr(compileErr)

	// Write Zonefile with given code
	writeErr := writeZone(code)
	handleErr(writeErr)
}


// Checks if an error is present, and raises it.
func handleErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
