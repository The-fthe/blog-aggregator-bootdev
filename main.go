package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"the-fthe/blog-aggregator-bootdev/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Enviroment file load failed")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("POPRT not set in .env file")
	}

	dbURL := os.Getenv("DATABASE")
	if dbURL == "" {
		log.Fatal("DATABASE_URL enviroment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("GET /v1/health", handlerReadines)
	mux.HandleFunc("GET /v1/err", handlerErr)
	mux.HandleFunc("GET /v1/users", apiCfg.handlerGetUsers)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server running on port", port)
	log.Fatal(srv.ListenAndServe())

}
