package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"encoding/hex"

	"github.com/farisbrandone/RSS_SERVER/internal/database"
	"github.com/google/uuid"
	//"github.com/satori/go.uuid"
)

type UserName struct {
	Name string `json:"name"`
}

type FeedFollowsBody struct {
	FeedID string `json:"feed_id"`
}

type FeedBody struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type AutomateResult struct {
	Feed        database.Feed       `json:"feed"`
	FeedFollows database.FeedFollow `json:"feed_follows"`
}

func (cfg apiConfig) handleCreateUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	param := UserName{}
	err := decoder.Decode(&param)
	//fmt.Println(params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	log.Println("my name", param)
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	randomString, err1 := GenerateRandomBytes(64)
	if err1 != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	hmacHex := hex.EncodeToString(randomString)
	log.Println((hmacHex))
	value := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: sql.NullString{
			String: param.Name,
			Valid:  true,
		},
		ApiKey: hmacHex,
	}
	log.Println("my name", value)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	result, err := cfg.DB.CreateUser(ctx, value)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 201, result)
}

func (cfg apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	param := FeedBody{}
	err := decoder.Decode(&param)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	value := database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: sql.NullString{
			String: param.Name,
			Valid:  true,
		},
		Url:    param.Url,
		UserID: user.ID,
	}

	result, err := cfg.DB.CreateFeeds(ctx, value)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	value2 := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		UserID:    result.UserID,
		FeedID:    result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	result2, err := cfg.DB.CreateFeedFollows(ctx, value2)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	finalResult := AutomateResult{
		Feed:        result,
		FeedFollows: result2,
	}

	respondWithJSON(w, 200, finalResult)
}

func (cfg apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	param := FeedFollowsBody{}
	err := decoder.Decode(&param)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	fmt.Println(param.FeedID)
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	val, _ := uuid.ParseBytes([]byte(param.FeedID))
	//valo := []byte(param.FeedID)
	fmt.Println(val)
	feeds, err := cfg.DB.GetoneFeeds(ctx, val)
	feedError := fmt.Sprintf(` The feedID %v is not available in the database`, val)
	if err != nil {
		respondWithError(w, 500, feedError)
		return
	}
	value := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feeds.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := cfg.DB.CreateFeedFollows(ctx, value)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 201, result)
}
