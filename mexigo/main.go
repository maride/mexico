package main

import (
	"flag"
	"github.com/maride/mexico/mexigo/interpreter"
	"log"
)

func main() {
	// Important things first
	printBanner()

	// Get desired domain off arguments
	flag.Parse()
	domain := flag.Arg(0)
	if domain == "" {
		// No domain entered.
		log.Println("Please specify a domain to receive code from as first argument, like this: ./mexigo <domain>")
		return
	}

	// Get program code from that domain
	log.Printf("Resolving %s for MX records", domain)
	code := LookupMX(domain)
	if len(code) == 0 {
		// Failed to look up mexico code on that domain. Log and exit.
		log.Printf("No code found on domain '%s'. Exiting.", domain)
		return
	}

	// Inform user about successful resolving
	log.Printf("Found %d code lines, interpreting them...", len(code))

	// Set up interpreter
	runErr := interpreter.Run(code)
	if runErr != nil {
		// Encountered error while executing code. Log and exit.
		log.Println(runErr.Error())
		return
	}
}

// Prints a small banner :)
func printBanner() {
	println("mexigo - the reference interpreter for the mexico esolang!")
	println("See github.com/maride/mexico for further information.")
	println()
}
