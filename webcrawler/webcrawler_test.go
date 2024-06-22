package webcrawler

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func TestBasicGetHTML(t *testing.T) {
	testHTML := "<!DOCTYPE html><html><head></head><body><h1>Fetch test!</h1></body></html>"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	cfg := Config{
		PoolSize:       1,
		Timeout:        5,
		MaxDepth:       -1,
		AllowedDomains: []string{"localhost", "127.0.0.1"},
		Logger:         slog.Default(),
	}

	crawler := NewWebCrawler(cfg)

	urls := []string{server.URL}
	results := crawler.GetHTML(urls)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	result := results[0]
	if result == nil {
		t.Fatalf("expected non-nil result")
	}

	resultStr, err := nodeToString(result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resultStr != testHTML {
		t.Errorf("expected html %q but got %q", testHTML, resultStr)
	}
}

func TestGetHTMLDeniedDomain(t *testing.T) {
	testHTML := "<!DOCTYPE html><html><head></head><body><h1>Fetch test!</h1></body></html>"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	cfg := Config{
		PoolSize:       1,
		Timeout:        5,
		MaxDepth:       -1,
		AllowedDomains: []string{"foo", "bar"},
		Logger:         slog.Default(),
	}

	crawler := NewWebCrawler(cfg)

	urls := []string{server.URL}
	results := crawler.GetHTML(urls)

	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestGetHTMLTimeout(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	cfg := Config{
		PoolSize:       1,
		Timeout:        1,
		MaxDepth:       -1,
		AllowedDomains: []string{"localhost"},
		Logger:         slog.Default(),
	}
	crawler := NewWebCrawler(cfg)

	_, err := crawler.getHTML(server.URL)
	if err == nil {
		t.Fatal("expected timeout error, got none")
	}
}

func TestGetHTMLMaxDepth(t *testing.T) {
	testHTML := "<!DOCTYPE html><html><head></head><body><h1>Fetch test!</h1></body></html>"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	cfg := Config{
		PoolSize:       1,
		Timeout:        5,
		MaxDepth:       0,
		AllowedDomains: []string{"foo", "bar"},
		Logger:         slog.Default(),
	}

	crawler := NewWebCrawler(cfg)

	urls := []string{server.URL}
	results := crawler.GetHTML(urls)

	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

// nodeToString converts an html.Node to its string representation.
func nodeToString(n *html.Node) (string, error) {
	var sb strings.Builder
	err := html.Render(&sb, n)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}
