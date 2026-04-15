package scraper

import (
	"io"
	"net/http"
)

// return html hasil scraping
func Scraper(url string) ([]byte, int, error) {
	response, err := http.Get(url)
	if err != nil {
		// handle error
		return nil, response.StatusCode, err
	}

	defer response.Body.Close()

	bytes, _ := io.ReadAll(response.Body)
	return bytes, response.StatusCode, nil
}
