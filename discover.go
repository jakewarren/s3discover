package s3discover

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/gocolly/colly"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

// Discover scrapes a domain, returning a list of buckets discovered
func Discover(domain string) (buckets []string) {

	// map to hold all unique s3bucket found
	s3buckets := make(map[string]bool)

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

	// store all buckets that were found
	for key := range s3buckets {
		buckets = append(buckets, key)
	}

	return
}
