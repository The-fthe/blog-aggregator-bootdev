package main

import (
	"encoding/json"
	"net/http"
	"time"

	"the-fthe/blog-aggregator-bootdev/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn'% decode paramters")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Create user failed")
		return
	}
	responseWithJSON(w, http.StatusOK, databaseUserToUser(user))

}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
