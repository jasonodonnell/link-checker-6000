package fetcher

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Fetcher struct {
	cache          sync.Map
	poolSize       int
	deadLinks      []string
	timeout        time.Duration
	wg             *sync.WaitGroup
	logger         *slog.Logger
	throttle       time.Duration
	currentDepth   int
	maxDepth       int
	allowedDomains []string
	deniedDomains  []string
}

func NewFetcher(poolSize, maxDepth, timeout int, allowedDomains []string, deniedDomains []string, logger *slog.Logger) *Fetcher {
	return &Fetcher{
		cache:          sync.Map{},
		poolSize:       poolSize,
		deadLinks:      []string{},
		timeout:        time.Second * time.Duration(timeout),
		wg:             &sync.WaitGroup{},
		logger:         logger,
		throttle:       time.Millisecond * 250,
		maxDepth:       maxDepth,
		allowedDomains: allowedDomains,
		deniedDomains:  deniedDomains,
	}
}

func (f *Fetcher) GetHTML(urls []string) []*html.Node {
	jobs := make(chan string)
	outCh := make(chan *html.Node)
	errCh := make(chan error)
	results := []*html.Node{}

	if len(urls) == 0 || f.currentDepth >= f.maxDepth {
		f.logger.Info("Max depth exceeded")
		return nil
	}

	f.currentDepth += 1

	for i := 0; i < f.poolSize; i++ {
		f.wg.Add(1)
		go f.crawl(jobs, outCh, errCh)
	}

	go func() {
		for _, url := range urls {
			f.logger.Info("Loading job", "url", url)
			jobs <- url
		}
		close(jobs)
	}()

	go func() {
		f.wg.Wait()
		close(outCh)
		close(errCh)
	}()

	for {
		select {
		case value, ok := <-outCh:
			if ok {
				results = append(results, value)
			} else {
				outCh = nil
			}
		case err, ok := <-errCh:
			if ok {
				f.logger.Error(err.Error())
			} else {
				errCh = nil
			}
		}

		if outCh == nil && errCh == nil {
			break
		}
	}
	return results
}

func (f *Fetcher) crawl(jobs <-chan string, out chan *html.Node, errs chan error) {
	defer f.wg.Done()
Exit:
	for path := range jobs {
		path, err := url.PathUnescape(path)
		if err != nil {
			errs <- err
			continue
		}

		if _, exists := f.cache.Load(path); exists {
			f.logger.Info("Cache hit", "url", path)
			continue
		}

		for _, denied := range f.deniedDomains {
			if strings.Contains(path, denied) {
				f.logger.Info("Denying url", "url", path)
				goto Exit
			}
		}

		f.cache.Store(path, struct{}{})
		html, err := f.getHTML(path)
		if err != nil {
			errs <- err
			continue
		}

		for _, allowed := range f.allowedDomains {
			if strings.Contains(path, allowed) {
				out <- html
			}
		}
	}
}

func (f *Fetcher) getHTML(url string) (*html.Node, error) {
	client := http.Client{
		Timeout: f.timeout,
	}

	time.Sleep(f.throttle)
	resp, err := client.Get(url)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			f.deadLinks = append(f.deadLinks, url)
		} else {
			return nil, err
		}
	}

	if resp == nil {
		return nil, fmt.Errorf("nil response from %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		f.deadLinks = append(f.deadLinks, url)
	}

	return html.Parse(resp.Body)
}

func (f *Fetcher) DeadLinks() []string {
	return f.deadLinks
}
