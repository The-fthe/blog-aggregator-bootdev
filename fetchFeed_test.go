package main

import (
	"fmt"
	"testing"
)

func TestFeedFetch(t *testing.T) {
	testUrls := []string{
		"https://blog.boot.dev/index.xml",
		//"https://wagslane.dev/index.xml",
	}

	for _, url := range testUrls {
		feed, err := FetchDataFromFeed(url)
		if err != nil {
			t.Errorf("FetchDataFromFeed failed to fetch data")
		}
		fmt.Printf("feed: %v\n", databaseFeedToFeed(feed))
	}

}
