package scraper

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

// return html hasil scraping
func Scraper(url string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DOMTraversalBot/1.0)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, fmt.Errorf("read body failed: %w", err)
	}
	return bytes, response.StatusCode, nil
}
