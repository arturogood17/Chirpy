package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondErrorWriter(w http.ResponseWriter, error_code int, msg string, err error) { //Si tiene error, lo imprimes
	type ErrorType struct { //Esta función no tiene ningún tipo de lógica de marshalling. Eso lo hace responWithJson
		Error string `json:"error"`
	}
	if err != nil {
		log.Println(err)
	}
	if error_code > 499 {
		log.Printf("responding with 5XX error: %s\n", msg)
	}
	responWithJson(w, error_code, ErrorType{Error: msg})
}

func responWithJson(w http.ResponseWriter, code int, payload any) { //Esta función es la que hace el marshalling
	w.Header().Set("Content-Type", "application/json") //Es importante establecer el Content-Type
	data, err := json.Marshal(payload)                 //Marshal no necesita un struct. Solo convierte a JSON
	if err != nil {                                    //El payload necesita ser any porque tiene que aceptar todo
		log.Printf("error marshalling JSON - %v", err) //y marsheallearlo
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
