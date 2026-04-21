package scraper

import (
	"testing"
)

// https://www.quora.com/Which-webpages-are-impossible-to-scrape

// Hasil test :
// Scrapper saat ini cm bs yg "public aja" kek yg ga perlu autentikasi gitu"
// Scrapper ga bs akses web yg ngedetect dia sbgai bot
func TestScraperTableDrivenTestsMethod(t *testing.T) {
	urls := ([]string{
		"https://example.com/",
		"https://www.quora.com/Which-webpages-are-impossible-to-scrape",
		"https://id.quora.com/",
	})

	for _, url := range urls {
		_, statusCode, err := Scraper(url)
		if err != nil || statusCode != 200 {
			t.Errorf("Input URL: %s\nError: %s\n StatusCode: %d\n\n", url, err, statusCode)
		}
	}
}
