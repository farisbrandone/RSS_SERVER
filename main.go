package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/farisbrandone/RSS_SERVER/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type StateStruct struct{}

type StatusCode struct {
	Status string `json:"status"`
}
type StatusError struct {
	Error string `json:"error"`
}

type apiConfig struct {
	DB *database.Queries
}

func main() {
	//st:=StateStruct{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	dbURL := os.Getenv("dbURL")
	db, err1 := sql.Open("postgres", dbURL)

	if err1 != nil {
		log.Println("Something went wrong when we create data")
		return
	}
	dbQueries := database.New(db)

	totalConfig := apiConfig{
		DB: dbQueries,
	}

	middlewareCors := cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r := chi.NewRouter()
	r.Use(middlewareCors)
	//r.Get("/",  FirstPage)
	r.Mount("/v1", totalConfig.HomePage())
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg apiConfig) HomePage() chi.Router {
	r := chi.NewRouter()

	r.Get("/readiness", handleReadiness)
	r.Get("/err", handleError)
	r.Post("/users", cfg.handleCreateUsers)
	r.Get("/user", cfg.middleWareAuthentication(HandlerGetOwnUser))
	r.Post("/feeds", cfg.middleWareAuthentication(cfg.handlerCreateFeeds))
	r.Get("/feeds", cfg.handleGetAllFeeds)
	r.Post("/feed_follows", cfg.middleWareAuthentication(cfg.handlerCreateFeedFollows))
	r.Delete("/feed_follows/{feedFollowID}", cfg.middleWareAuthentication(cfg.handlerDeleteFeedFollows))
	r.Get("/feed_follows", cfg.middleWareAuthentication(cfg.handleGetAllFeedFollows))
	r.Get("/posts", cfg.middleWareAuthentication(cfg.handlerPostsGet))
	return r
}

func handleReadiness(w http.ResponseWriter, r *http.Request) { // founction a executer pour une requette donnée
	payload := StatusCode{
		Status: "ok",
	}
	respondWithJSON(w, 200, payload)
}

func handleError(w http.ResponseWriter, r *http.Request) { // founction a executer pour une requette donnée
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //handler retourné, recois les requète les traites et pass a next paramètre d'entrée
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "POST" && (r.URL.Path == "/healthz" || r.URL.Path == "/metrics") {
			log.Println("Serving on port: ", r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}
		next.ServeHTTP(w, r)
	})
}
