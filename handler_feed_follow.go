package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"the-fthe/blog-aggregator-bootdev/internal/database"

	"github.com/google/uuid"
)

type Response struct {
	Message string `json:"message"`
}

func (cfg *apiConfig) handlerFeedFollowCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feedUuid, err := uuid.Parse(params.FeedID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "feed UUID can't be parse")
		return
	}
	_, err = cfg.DB.GetFeed(r.Context(), feedUuid)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "UUID feed doesn't exist")
		return

	}

	log.Println("Feed ID: ", feedUuid.String())
	log.Println("User ID: ", user.ID)

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedUuid,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Create feed follow failed")
		return
	}
	responseWithJSON(w, http.StatusOK, feedFollow)

}

func (cfg *apiConfig) handlerFeedFollowDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	type paramters struct {
		FeedFollowID string
	}

	feedFollowUuidstr := r.PathValue("feedFollowID")
	feedFollowUuid, err := uuid.Parse(feedFollowUuidstr)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "UUID feed follow doesn't exist")
		return
	}

	err = cfg.DB.DeleteFeedFollow(
		r.Context(), database.DeleteFeedFollowParams{ID: feedFollowUuid, UserID: user.ID})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "detele feedfollow failed")
	}
	responseWithJSON(w, http.StatusOK, Response{Message: "delete successful"})
}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollows(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "get feedFollows failed")
		return
	}
	responseWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))

}
