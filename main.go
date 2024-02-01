package main

import (
	"log"
	"net/http"
	"os"
    "github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type StateStruct struct { }

type StatusCode struct {
	Status string `json:"status"`
}
type StatusError struct {
	Error string `json:"error"`
}

func main(){
	//st:=StateStruct{}
    err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	  }
	port := os.Getenv("PORT")
	middlewareCors:=cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		  })

	r:=chi.NewRouter()
	r.Use(middlewareCors)
	//r.Get("/",  FirstPage)
	r.Mount("/v1", HomePage())
	srv := &http.Server{ 
		Addr:    ":" + port,
		Handler: r,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe()) 
}




 

  func  HomePage() chi.Router {
	r := chi.NewRouter()
	
	r.Get("/readiness",  handleReadiness)
	r.Get("/err",  handleError)
	return r
  }

  func  handleReadiness(w http.ResponseWriter, r *http.Request) { // founction a executer pour une requette donnée
	payload:= StatusCode{
		Status:"ok",
	}
	respondWithJSON(w,200,payload)
}
func  handleError(w http.ResponseWriter, r *http.Request) { // founction a executer pour une requette donnée
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {//handler retourné, recois les requète les traites et pass a next paramètre d'entrée
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "POST" && (r.URL.Path=="/healthz" || r.URL.Path=="/metrics") {
			log.Println("Serving on port: ", r.URL.Path)
	        w.WriteHeader(http.StatusMethodNotAllowed)
	        w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}
		next.ServeHTTP(w, r)
	})
} 

