package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type count int

type Crawler struct {
	pages              map[string]count
	maxPages           int
	paths              map[string]count
	baseURL            *url.URL
	mu                 *sync.Mutex
	pathMu             *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (c *Crawler) Crawl(currentURL *url.URL) {
	c.concurrencyControl <- struct{}{}
	defer func() {
		<-c.concurrencyControl
		c.wg.Done()
	}()

	c.mu.Lock()
	if c.maxPages > 0 && len(c.pages) >= c.maxPages {
		c.mu.Unlock()
    fmt.Println("max pages reached")
		return
	}
	c.mu.Unlock()

	if c.baseURL.Hostname() != currentURL.Hostname() && !strings.HasSuffix(currentURL.Hostname(), c.baseURL.Hostname()) {
		fmt.Printf("skipping %s, different hostname\n", currentURL.String())
		return
	}

	normalized, err := NormalizeURL(currentURL.String())
	if err != nil {
		fmt.Printf("error normalizing url %s: %s\n", currentURL, err)
		return
	}

	count := c.addVisit(normalized)
	c.addPath(currentURL.Path)
	if count > 1 {
		fmt.Printf("skipping %s, already visited\n", currentURL.String())
		return
	}

	fmt.Printf("crawling %s\n", currentURL.String())
	html, err := GetHTML(*currentURL)
	if err != nil {
		fmt.Printf("error fetching html: %s\n", err)
		return
	}

	links, err := ReadLinks(html, c.baseURL)
	if err != nil {
		fmt.Printf("error reading links: %s\n", err)
		return
	}
	fmt.Printf("found %d links\n", len(links))

	for _, link := range links {
		linkURL, err := url.Parse(link)
		if err != nil {
			fmt.Printf("error parsing link: %s\n", err)
			continue
		}

		c.wg.Add(1)
		go c.Crawl(linkURL)
	}
}

func (c *Crawler) addVisit(normalized string) count {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.pages[normalized]; ok {
		c.pages[normalized]++
		return c.pages[normalized]
	}

	c.pages[normalized] = 1
	return 1
}

func (c *Crawler) addPath(path string) {
	c.pathMu.Lock()
	defer c.pathMu.Unlock()
	if _, ok := c.paths[path]; ok {
		c.paths[path]++
		return
	}

	c.paths[path] = 1
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no website provided")
		fmt.Println("usage: snail <url>")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	target, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("invalid url: %s\n", args[0])
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", target.String())

	crawler := Crawler{
		pages:              make(map[string]count),
		// maxPages:           100,
		paths:              make(map[string]count),
		baseURL:            target,
		mu:                 &sync.Mutex{},
		pathMu:             &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 8),
		wg:                 &sync.WaitGroup{},
	}

	fmt.Println("concurrency control:", cap(crawler.concurrencyControl))

	crawler.wg.Add(1)
	go crawler.Crawl(target)
	crawler.wg.Wait()

	fmt.Printf("\n-------------------------\n\n")
	fmt.Println("Total Unique Paths: ", len(crawler.paths))
	fmt.Println()
	for path, count := range crawler.paths {
		fmt.Printf("%s: %d times\n", path, count)
	}

	fmt.Printf("\n-------------------------\n\n")
	fmt.Println("Total Unique Pages: ", len(crawler.pages))
	fmt.Println()
	for page, count := range crawler.pages {
		fmt.Printf("%s: %d times\n", page, count)
	}
}

func GetHTML(url url.URL) (string, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("status code: %d", res.StatusCode)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("content type: %s", res.Header.Get("Content-Type"))
	}

	html, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err)
	}

	return string(html), nil
}
