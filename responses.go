package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondErrorWriter(w http.ResponseWriter, error_code int, err error) {
	type ErrorType struct {
		Error string `json:"error"`
	}
	errType := ErrorType{
		Error: err.Error(),
	}
	data, err := json.Marshal(errType)
	if err != nil {
		responWithJson(w, http.StatusInternalServerError, []byte("Could not marshal error"))
		return
	}
	responWithJson(w, error_code, data)
}

func responWithJson(w http.ResponseWriter, code int, body []byte) {
	w.WriteHeader(code)
	if _, err := w.Write(body); err != nil {
		log.Fatalf("Could respond - %v", err)
	}
}
