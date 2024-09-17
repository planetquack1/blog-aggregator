package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/planetquack1/blog-aggregator/internal/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {

	type Status struct {
		Status string `json:"status"`
	}

	status := Status{
		Status: "OK",
	}

	respondWithJSON(w, http.StatusOK, status)
}

func handlerError(w http.ResponseWriter, r *http.Request) {

	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (cfg *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {

	type Name struct {
		Name string `json:"name"`
	}

	// Read in struct
	name, err := decodeJSON(r, Name{})
	if err != nil {
		respondWithError(w, 400, "Failed to parse JSON")
		return
	}

	currentTime := time.Now()

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      name.Name,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Marshall and write
	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request, user database.User) {

	// Marshall and write
	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) postFeeds(w http.ResponseWriter, r *http.Request, user database.User) {

	type Feed struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	// Read in struct
	feed, err := decodeJSON(r, Feed{})
	if err != nil {
		respondWithError(w, 400, "Failed to parse JSON")
		return
	}

	currentTime := time.Now()

	randomFeedUUID := uuid.New()

	addedFeed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        randomFeedUUID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      feed.Name,
		Url:       feed.URL,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Follow the feed
	addedFeedFollows, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    randomFeedUUID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	type FeedAndFeedFollows struct {
		Feed        database.Feed       `json:"feed"`
		FeedFollows database.FeedFollow `json:"feed_follows"`
	}
	feedAndFeedFollows := FeedAndFeedFollows{
		Feed:        addedFeed,
		FeedFollows: addedFeedFollows,
	}

	// Marshall and write
	dat, err := json.Marshal(feedAndFeedFollows)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {

	allFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unauthorized")
		return
	}

	// Marshall and write
	dat, err := json.Marshal(allFeeds)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) postFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	type Feed struct {
		ID string `json:"feed_id"`
	}

	// Read in struct
	feed, err := decodeJSON(r, Feed{})
	if err != nil {
		respondWithError(w, 400, "Failed to parse JSON")
		return
	}

	// Parse feed ID
	feedID, err := uuid.Parse(feed.ID)
	if err != nil {
		respondWithError(w, 400, "Failed to parse ID from string to UUID")
		return
	}

	currentTime := time.Now()

	addedFeedFollows, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feedID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed follows")
		return
	}

	// Marshall and write
	dat, err := json.Marshal(addedFeedFollows)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) deleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	// Extract feedFollowID from the path using PathValue
	feedFollowIDStr := r.PathValue("feedFollowID")
	if feedFollowIDStr == "" {
		respondWithError(w, 400, "feedFollowID not provided")
		return
	}

	// Parse feed ID
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, "Failed to parse ID from string to UUID")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete feed follow")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	allUserFeedFollows, err := cfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unauthorized")
		return
	}

	// Marshall and write
	dat, err := json.Marshal(allUserFeedFollows)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) getPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		log.Printf("Could not get posts: %s", err)
		return
	}

	// Marshall and write
	dat, err := json.Marshal(posts)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
