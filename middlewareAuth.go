package main

import (
	"net/http"

	"the-fthe/blog-aggregator-bootdev/internal/auth"
	"the-fthe/blog-aggregator-bootdev/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			http.Error(w, "ApiKey empty", http.StatusForbidden)
			return
		}
		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			http.Error(w, "ApiKey invalid", http.StatusForbidden)
			return
		}
		handler(w, r, user)

	})

}
