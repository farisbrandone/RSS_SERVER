package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/farisbrandone/RSS_SERVER/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedID := chi.URLParam(r, "feedFollowID")
	fmt.Println(feedID)
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	val, _ := uuid.ParseBytes([]byte(feedID))

	result, err := cfg.DB.DeleteFeedFollows(ctx, val)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 200, result)
}
