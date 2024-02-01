package main

import (
	"net/http"
  "encoding/json"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}){
	dat, err := json.Marshal(payload)
		if err !=nil {
				
				w.WriteHeader(http.StatusInternalServerError)
				return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(dat)  
		
	}
