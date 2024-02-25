package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/farisbrandone/RSS_SERVER/internal/database"
)

func (cfg apiConfig) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	headers := r.Header.Get("Authorization")
	if headers == "" {
		respondWithError(w, 500, "no api keys")
		return
	}
	apiKeyStrings := strings.Split(headers, " ")
	if len(apiKeyStrings) < 2 || apiKeyStrings[0] != "ApiKey" {
		respondWithError(w, 500, "no values for api keys")
		return
	}

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	user, err := cfg.DB.GetUserByApiKey(ctx, apiKeyStrings[1])

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 201, user)
}

func HandlerGetOwnUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 201, user)
}

func (cfg apiConfig) handleGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	feeds, err := cfg.DB.GetAllFeeds(ctx)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 201, feeds)
}

func (cfg apiConfig) handleGetAllFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	feedFollows, err := cfg.DB.GetAllFeedFollows(ctx)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 200, feedFollows)
}
