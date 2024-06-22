package webcrawler

import (
	"errors"
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

var errStatusNotFound = errors.New("404 not found")

type Config struct {
	PoolSize       int
	Timeout        int
	MaxDepth       int
	AllowedDomains []string
	DeniedDomains  []string
	Logger         *slog.Logger
}

type WebCrawler struct {
	cache        sync.Map
	deadLinks    []string
	wg           *sync.WaitGroup
	throttle     time.Duration
	currentDepth int
	Config
}

func NewWebCrawler(cfg Config) *WebCrawler {
	return &WebCrawler{
		cache:     sync.Map{},
		deadLinks: []string{},
		wg:        &sync.WaitGroup{},
		throttle:  time.Millisecond * 250,
		Config:    cfg,
	}
}

func (wc *WebCrawler) GetHTML(urls []string) []*html.Node {
	jobs := make(chan string)
	outCh := make(chan *html.Node)
	errCh := make(chan error)
	results := []*html.Node{}

	if len(urls) == 0 || (wc.MaxDepth != -1 && wc.currentDepth >= wc.MaxDepth) {
		wc.Logger.Info("Max depth exceeded")
		return nil
	}

	wc.currentDepth += 1

	for i := 0; i < wc.PoolSize; i++ {
		wc.wg.Add(1)
		go wc.crawl(jobs, outCh, errCh)
	}

	go func() {
		for _, url := range urls {
			wc.Logger.Info("Loading job", "url", url)
			jobs <- url
		}
		close(jobs)
	}()

	go func() {
		wc.wg.Wait()
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
				wc.Logger.Error(err.Error())
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

func (wc *WebCrawler) crawl(jobs <-chan string, out chan *html.Node, errs chan error) {
	defer wc.wg.Done()
	for path := range jobs {
		path, err := url.PathUnescape(path)
		if err != nil {
			errs <- err
			continue
		}

		if _, exists := wc.cache.Load(path); exists {
			wc.Logger.Info("Cache hit", "url", path)
			continue
		}

		if wc.deniedURL(path) {
			wc.Logger.Info("Skipping denied url", "url", path)
			continue
		}

		wc.cache.Store(path, struct{}{})

		time.Sleep(wc.throttle)
		html, err := wc.getHTML(path)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() || err == errStatusNotFound {
				wc.deadLinks = append(wc.deadLinks, path)
			} else {
				errs <- err
			}
			continue
		}

		for _, allowed := range wc.AllowedDomains {
			if strings.Contains(path, allowed) {
				out <- html
			}
		}
	}
}

func (wc *WebCrawler) getHTML(url string) (*html.Node, error) {
	client := http.Client{
		Timeout: time.Second * time.Duration(wc.Timeout),
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("404 not found")
	}

	return html.Parse(resp.Body)
}

func (wc *WebCrawler) DeadLinks() []string {
	return wc.deadLinks
}

func (wc *WebCrawler) deniedURL(url string) bool {
	for _, denied := range wc.DeniedDomains {
		if strings.Contains(url, denied) {
			return true
		}
	}
	return false
}
