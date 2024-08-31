package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"the-fthe/blog-aggregator-bootdev/internal/database"
	"time"

	"github.com/google/uuid"
)

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
