package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type URLCache struct {
	urls map[string]string
	mtx  sync.Mutex
}

func (cache URLCache) Add(url, body string) bool {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()
	if _, ok := cache.urls[url]; !ok {
		cache.urls[url] = body
		return true
	}
	return false
}

func (cache URLCache) Get(url string) (v string, ok bool) {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()
	v, ok = cache.urls[url]
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	fmt.Printf("Crawl(%s, %d, %v)\n", url, depth, fetcher)
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	cache.Add(url, body)
	fmt.Printf("found: %s %q\n", url, body)
	done := make(chan bool)
	for _, u := range urls {
		fmt.Printf("checking %s\n", u)
		if _, ok := cache.Get(u); !ok {
			fmt.Printf("not found %s\n", u)
			go func(url string) {
				Crawl(u, depth-1, fetcher)
				done <- true
			}(u)
		} else {
			fmt.Printf("found %s\n", u)
		}
	}
	for _, u := range urls {
		fmt.Printf("Waiting on %s\n", u)
		<-done
	}
	return
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

var cache = URLCache{urls: make(map[string]string)}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
