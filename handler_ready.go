package main

import "net/http"

func handlerReadines(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, http.StatusInternalServerError, "InternalServe Error")
}
