package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gocolly/colly"
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

	// map to hold all unique s3bucket found
	s3buckets := make(map[string]bool)
	domain := pflag.Arg(0)

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only root url and the root url for any subdomain
		colly.URLFilters(
			regexp.MustCompile(fmt.Sprintf(`^https?://(\w+\.)*?%s/?$`, regexp.QuoteMeta(domain))),
		),
		// set the max recursion depth to 1 page
		colly.MaxDepth(1),
		// set a user agent for stealth
		colly.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36`),
	)
	// disable ssl verification
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	// Scrape all links on page
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		log.Debugf("Link found: %q -> %s", e.Text, link)

		// Visit links found on page
		// Only those links are visited which are matched by  any of the URLFilter regexps
		c.Visit(e.Request.AbsoluteURL(link))

	})

	// On each page, extract any s3 buckets from the HTML
	c.OnHTML("html", func(e *colly.HTMLElement) {
		html, _ := e.DOM.Html()

		// check for bucketname.s3.amazonaws.com format buckets
		var re = regexp.MustCompile(`[\w-\.]+\.s3(?:[-\w]+)?\.amazonaws\.com`)

		for _, match := range re.FindAllString(html, -1) {
			s3buckets[match] = true
		}

		// check for s3.amazonaws.com/bucketname/ format buckets
		re = regexp.MustCompile(`//s3[\w-]*?\.amazonaws\.com/[\w-]+/`)

		for _, match := range re.FindAllString(html, -1) {
			s3buckets[match] = true
		}

	})

	// Log page visits if verbose output is enabled
	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %s", r.URL.String())
	})

	// if the user only provided a domain name, append the http protocol scheme for them
	if !strings.HasPrefix(domain, "http") {
		domain = "http://" + domain
	}

	// Start scraping
	c.Visit(domain)

	// Print all buckets that were found
	for key := range s3buckets {
		fmt.Println(key)
	}
}

// print custom usage instead of the default provided by pflag
func displayUsage() {

	fmt.Printf("Usage: s3discover [<flags>] <domain>\n\n")
	fmt.Printf("Example: s3discover github.com\n\n")
	fmt.Printf("Optional flags:\n\n")
	pflag.PrintDefaults()
}
