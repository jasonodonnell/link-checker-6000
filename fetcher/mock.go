package fetcher

import (
	"log/slog"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type MockFetcher struct {
	cache        sync.Map
	poolSize     int
	deadLinks    []string
	timeout      time.Duration
	wg           *sync.WaitGroup
	logger       *slog.Logger
	currentDepth int
	maxDepth     int
}

func NewMockFetcher(urls []string, poolSize, maxDepth int, timeout time.Duration) *MockFetcher {
	return &MockFetcher{
		cache:     sync.Map{},
		poolSize:  poolSize,
		deadLinks: []string{},
		timeout:   timeout,
		wg:        &sync.WaitGroup{},
		logger:    slog.Default(),
		maxDepth:  maxDepth,
	}
}

func (m *MockFetcher) GetHTML(urls []string) []*html.Node {
	reader := strings.NewReader(htmlExample)
	doc, err := html.Parse(reader)
	if err != nil {
		panic(err)
	}

	return []*html.Node{doc}
}

func (m *MockFetcher) DeadLinks() []string {
	return m.deadLinks
}

var htmlExample = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Example Page with 30 Links</title>
</head>
<body>
    <h1>Welcome to the Example Page with 30 Links</h1>
    <ul>
        <li><a href="https://example.com/link1">Link 1</a></li>
        <li><a href="https://example.com/link2">Link 2</a></li>
        <li><a href="https://example.com/link3">Link 3</a></li>
        <li><a href="https://example.com/link4">Link 4</a></li>
        <li><a href="https://example.com/link5">Link 5</a></li>
        <li><a href="https://example.com/link6">Link 6</a></li>
        <li><a href="https://example.com/link7">Link 7</a></li>
        <li><a href="https://example.com/link8">Link 8</a></li>
        <li><a href="https://example.com/link9">Link 9</a></li>
        <li><a href="https://example.com/link10">Link 10</a></li>
        <li><a href="https://example.com/link11">Link 11</a></li>
        <li><a href="https://example.com/link12">Link 12</a></li>
        <li><a href="https://example.com/link13">Link 13</a></li>
        <li><a href="https://example.com/link14">Link 14</a></li>
        <li><a href="https://example.com/link15">Link 15</a></li>
        <li><a href="https://example.com/link16">Link 16</a></li>
        <li><a href="https://example.com/link17">Link 17</a></li>
        <li><a href="https://example.com/link18">Link 18</a></li>
        <li><a href="https://example.com/link19">Link 19</a></li>
        <li><a href="https://example.com/link20">Link 20</a></li>
        <li><a href="https://example.com/link21">Link 21</a></li>
        <li><a href="https://example.com/link22">Link 22</a></li>
        <li><a href="https://example.com/link23">Link 23</a></li>
        <li><a href="https://example.com/link24">Link 24</a></li>
        <li><a href="https://example.com/link25">Link 25</a></li>
        <li><a href="https://example.com/link26">Link 26</a></li>
        <li><a href="https://example.com/link27">Link 27</a></li>
        <li><a href="https://example.com/link28">Link 28</a></li>
        <li><a href="https://example.com/link29">Link 29</a></li>
        <li><a href="https://example.com/link30">Link 30</a></li>
    </ul>
</body>
</html>
`
