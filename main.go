package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"the-fthe/blog-aggregator-bootdev/internal/database"
	"time"

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

	//health and err
	mux.HandleFunc("GET /v1/health", handlerReadines)
	mux.HandleFunc("GET /v1/err", handlerErr)

	//users
	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))

	//feeds
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handleFeedCreate))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerFeedsGet)
	mux.HandleFunc("DELETE /v1/feeds/{feedID}", apiCfg.middlewareAuth(apiCfg.handlerFeedDelete))

	//feedfollow
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	fmt.Println("Server running on port", port)
	log.Fatal(srv.ListenAndServe())

}
