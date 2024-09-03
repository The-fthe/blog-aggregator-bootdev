package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"the-fthe/blog-aggregator-bootdev/internal/database"
	"time"

	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}
		log.Printf("Found %v feeds to fetch! ", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			if !feed.Url.Valid {
				log.Println("Invalid feed URL!")
				continue
			}
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}
	feedData, err := fetchFeed(feed.Url.String)

	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		//log.Println("Found post", item.Title)
		post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       sql.NullString{String: item.Title},
			Url:         sql.NullString{String: item.Link},
			Description: sql.NullString{String: item.Description},
			FeedID:      feed.ID,
		}
		layout := "Mon, 02 Jan 2006 15:04:05 -0700"
		parsedTime, err := time.Parse(layout, item.PubDate)
		if err != nil {
			log.Println("title ", post.Title, " Parse time failed")
			_, err = db.CreatePost(context.Background(), post)
			if err != nil {
				log.Println("non time post created failed: ", post.Title, err.Error())
			}
			continue
		}
		log.Println("PubDate: ", item.PubDate, "ParsedDate: ", parsedTime)
		publishAt := sql.NullTime{Time: parsedTime}
		post.PublishedAt = publishAt
		_, err = db.CreatePost(context.Background(), post)
		if err != nil {
			log.Println("post created failed: ", post.Title, err.Error())
		}

	}
	log.Printf("Feed %s collected, %v post found", feed.Name, len(feedData.Channel.Item))
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}
	return &rssFeed, nil
}
