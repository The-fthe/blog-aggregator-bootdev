package main

import (
	"net/http"
	"the-fthe/blog-aggregator-bootdev/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		if apiKey == "" {
			http.Error(w, "ApiKey empty", http.StatusForbidden)
			return
		}
		_, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			http.Error(w, "ApiKey invalid", http.StatusForbidden)
			return
		}

	})

}
