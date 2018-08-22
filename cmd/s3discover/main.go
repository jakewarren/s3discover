package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jakewarren/s3discover"
	"github.com/spf13/pflag"
)

var version string

func main() {

	displayHelp := pflag.BoolP("help", "h", false, "display help")
	debugLevel := pflag.BoolP("debug", "d", false, "enable debug logging")
	verboseLevel := pflag.BoolP("verbose", "v", false, "enable verbose output")
	displayVersion := pflag.BoolP("version", "V", false, "display version")
	pflag.Parse()

	if *displayVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// override the default usage display
	if *displayHelp {
		displayUsage()
		os.Exit(0)
	}

	// set up logging
	// logging to stderr so it can be easily discarded if user only wants the bucket names
	log.SetHandler(cli.New(os.Stderr))
	// set the default logging level
	log.SetLevel(log.WarnLevel)

	if *verboseLevel {
		log.SetLevel(log.InfoLevel)
	}

	if *debugLevel {
		log.SetLevel(log.DebugLevel)
	}

	domain := pflag.Arg(0)

	for _, bucket := range s3discover.Discover(domain) {
		fmt.Println(bucket)
	}
}

// print custom usage instead of the default provided by pflag
func displayUsage() {

	fmt.Printf("Usage: s3discover [<flags>] <domain>\n\n")
	fmt.Printf("Example: s3discover github.com\n\n")
	fmt.Printf("Optional flags:\n\n")
	pflag.PrintDefaults()
}
