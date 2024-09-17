package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/planetquack1/blog-aggregator/internal/database"
)

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	URL           string     `json:"url"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	var lastFetchedAt *time.Time

	// Set last fetched at time
	if feed.LastFetchedAt.Valid {
		lastFetchedAt = &feed.LastFetchedAt.Time
	} else {
		lastFetchedAt = nil
	}

	return Feed{
		ID:            feed.ID,
		URL:           feed.Url,
		LastFetchedAt: lastFetchedAt,
	}
}
