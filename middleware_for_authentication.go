package main

import (
	"net/http"

	"context"
	"strings"

	"github.com/farisbrandone/RSS_SERVER/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg apiConfig) middleWareAuthentication(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		handler(w, r, user)
	})

}
