package main

import (
	"the-fthe/blog-aggregator-bootdev/internal/database"
	"time"
)

type apiConfig struct {
	DB     *database.Queries
	N      int
	Ticker *time.Ticker
	FeedCh chan<- Feed
}
