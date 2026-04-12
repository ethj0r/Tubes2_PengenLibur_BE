package scraper

import (
	"io"
	"net/http"
)

// return html hasil scraping
func Scraper(url string) (string, int, error) {
	response, err := http.Get(url)
	if err != nil {
		// handle error
		return "", response.StatusCode, err
	}

	defer response.Body.Close()

	bytes, _ := io.ReadAll(response.Body)
	return string(bytes), response.StatusCode, nil
}
