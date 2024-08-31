package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"the-fthe/blog-aggregator-bootdev/internal/database"
	"time"
)

type FeedRss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (cfg *apiConfig) handleFeedCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	//TODO: add check for params.Name and params.Url

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      sql.NullString{String: params.Name, Valid: true},
		Url:       sql.NullString{String: params.Url, Valid: true},
		UserID:    user.ID,
	})
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Create feed failed")
		return
	}
	responseWithJSON(w, http.StatusOK, databaseFeedAndFeedFollowToFeedAndFeedFollow(feed, feedFollow))

}

func (cfg *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Get feeds from database failed")
		return
	}
	responseWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}

func FetchDataFromFeed(feedUrl string) (database.Feed, error) {

	r, err := http.Get(feedUrl)
	if err != nil {
		return database.Feed{}, errors.New("get feedUrl data failed")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return database.Feed{}, errors.New("read url body failed")
	}
	feedRssStr := string(body)
	var feedRss FeedRss
	err = xml.Unmarshal([]byte(feedRssStr), &feedRss)
	if err != nil {
		return database.Feed{}, errors.New("xmml Unmarshal failed")
	}

	return RssFeedToFeed(feedRss), nil
}

func RssFeedToFeed(feedRss FeedRss) database.Feed {
	return database.Feed{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		Name:          sql.NullString{String: feedRss.Channel.Title},
		Url:           sql.NullString{String: feedRss.Channel.Link.Text},
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC()},
	}

}
