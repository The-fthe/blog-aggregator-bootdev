package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

	//create mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		_, _ = res.Write([]byte("Hello world"))
	})

	mux.HandleFunc("GET /v1/health", handlerReadines)
	mux.HandleFunc("GET /v1/err", handlerErr)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server running on port", port)
	log.Fatal(srv.ListenAndServe())

}
