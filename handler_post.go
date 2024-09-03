package main

import (
	"net/http"
	"strconv"
	"the-fthe/blog-aggregator-bootdev/internal/database"
)

func (cfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.PathValue("limit")
	defaultLimit := 5

	limit, err := strconv.Atoi(limitStr)

	if err == nil {
		defaultLimit = limit
	}

	posts, err := cfg.DB.GetPostByUser(r.Context(), database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  int32(defaultLimit),
	})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "get post failed")
		return
	}
	responseWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
