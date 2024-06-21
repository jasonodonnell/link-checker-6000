package fetcher

import (
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Fetcher struct {
	html         chan []byte
	cache        sync.Map
	poolSize     int
	deadLinks    []string
	timeout      time.Duration
	wg           *sync.WaitGroup
	logger       *slog.Logger
	currentDepth int
	maxDepth     int
}

func NewFetcher(urls []string, poolSize, maxDepth int, timeout time.Duration) *Fetcher {
	return &Fetcher{
		cache:     sync.Map{},
		poolSize:  poolSize,
		deadLinks: []string{},
		timeout:   timeout,
		wg:        &sync.WaitGroup{},
		logger:    slog.Default(),
		maxDepth:  maxDepth,
	}
}

func (f *Fetcher) GetHTML(urls []string) []html.Node {
	jobs := make(chan string)
	outCh := make(chan html.Node)
	errCh := make(chan error)
	results := []html.Node{}

	if len(urls) == 0 || f.currentDepth >= f.maxDepth {
		return nil
	}
	f.currentDepth += 1

	for i := 0; i < f.poolSize; i++ {
		f.wg.Add(1)
		f.logger.Info("Adding crawler", "id", i)
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
		f.logger.Info("Wait group done!")
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

	f.logger.Info("Returning")
	return results
}

func (f *Fetcher) crawl(jobs <-chan string, out chan html.Node, errs chan error) {
	defer f.wg.Done()
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

		f.cache.Store(path, struct{}{})
		html, err := f.getHTML(path)
		if err != nil {
			errs <- err
			continue
		}

		out <- *html
	}
}

func (f *Fetcher) getHTML(url string) (*html.Node, error) {
	client := http.Client{
		Timeout: f.timeout,
	}

	time.Sleep(time.Millisecond * 500)
	resp, err := client.Get(url)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			f.deadLinks = append(f.deadLinks, url)
		} else {
			return nil, err
		}
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
