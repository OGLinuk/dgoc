package main

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type crawler struct {
	seed string
}

func newCrawler(s string, filters []string) *crawler {
	return &crawler{
		seed: s,
	}
}

// crawl creates an http client, make a Get request, extracts hyperlinks,
// writes uncrawled URLs to q.rw and q.seed to q.pw
func (c *crawler) crawl() ([]string, error) {
	check := make(map[string]struct{})
	var uniques []string

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 7,
	}

	resp, err := client.Get(c.seed)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		for _, collected := range c.extract(resp) {
			decoded, err := url.QueryUnescape(collected)
			if err != nil {
				return nil, err
			}

			if _, exists := check[decoded]; !exists {
				uniques = append(uniques, decoded)
				check[decoded] = struct{}{}
			}
		}
	} else {
		return nil, err
	}

	return uniques, nil
}
