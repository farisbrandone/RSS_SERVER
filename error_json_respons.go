package main

import (
    "net/http"
)
type ParamsErrors struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string){
	myError:=ParamsErrors{
		Error:msg,
	}
	respondWithJSON(w, code, myError)
}