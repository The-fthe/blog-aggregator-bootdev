package main

import (
	"database/sql"
	"github.com/google/uuid"
	"the-fthe/blog-aggregator-bootdev/internal/database"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Api_Key   string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Api_Key:   user.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        string     `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name.String,
		UserID:        feed.UserID.String(),
		Url:           feed.Url.String,
		LastFetchedAt: nullTimeToTimePtr(feed.LastFetchedAt),
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	result := make([]Feed, len(feeds))
	for i, feed := range feeds {
		result[i] = databaseFeedToFeed(feed)
	}
	return result
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdateAt:  feedFollow.UpdatedAt,
	}
}

func databaseFeedFollowsToFeedFollows(feedFollows []database.FeedFollow) []FeedFollow {
	results := make([]FeedFollow, len(feedFollows))
	for i, feedFollow := range feedFollows {
		results[i] = databaseFeedFollowToFeedFollow(feedFollow)
	}
	return results
}

type FeedAndFeedFollow struct {
	Feed       Feed       `json:"feed"`
	FeedFollow FeedFollow `json:"feedFollow"`
}

func databaseFeedAndFeedFollowToFeedAndFeedFollow(feed database.Feed, feedFollow database.FeedFollow) FeedAndFeedFollow {
	return FeedAndFeedFollow{
		Feed:       databaseFeedToFeed(feed),
		FeedFollow: databaseFeedFollowToFeedFollow(feedFollow),
	}
}

type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       *string
	Url         *string
	Description *string
	PublishedAt *time.Time
	FeedID      uuid.UUID
}

func databasePostsToPosts(posts []database.Post) []Post {
	results := make([]Post, len(posts))
	for i, post := range posts {
		results[i] = Post{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       nullStringToStringPtr(post.Title),
			Url:         nullStringToStringPtr(post.Url),
			Description: nullStringToStringPtr(post.Description),
			PublishedAt: nullTimeToTimePtr(post.PublishedAt),
			FeedID:      post.FeedID,
		}

	}
	return results

}

func nullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
