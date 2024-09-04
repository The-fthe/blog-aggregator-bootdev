package main

import (
	"net/http"
	"strconv"
	"the-fthe/blog-aggregator-bootdev/internal/database"
)

func (cfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}

	posts, err := cfg.DB.GetPostByUser(r.Context(), database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "get post failed")
		return
	}
	responseWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
