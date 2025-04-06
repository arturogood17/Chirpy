package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type ErrorChirp struct {
		Error string `json:"error"`
	}
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with a 5xx error msg: %v", err)

	}
	w.WriteHeader(code)
	respondWithJson(w, code, ErrorChirp{Error: msg + " " + fmt.Sprintf("%v", err)})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marsahling the JSON: %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
