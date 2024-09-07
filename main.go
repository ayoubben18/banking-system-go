package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	http.HandleFunc("/person", personHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World")
}

func personHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	constentType := r.Header.Get("Content-Type")
	if constentType != "application/json" {
		http.Error(w, "Content type not supported", http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var person Person

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&person)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch  {
		case err == io.EOF:
			http.Error(w, "Request body must not be empty", http.StatusBadRequest)
		case err.Error() == "json: unknown field {field}":
			http.Error(w, "Request body must not be larger than 1MB", http.StatusRequestEntityTooLarge)
		case err == io.ErrUnexpectedEOF:
			http.Error(w, "Request body contains badly-formed JSON", http.StatusBadRequest)
		case err.(*json.SyntaxError) != nil:
			http.Error(w, fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset), http.StatusBadRequest)
		case err.(*json.UnmarshalTypeError) != nil:
			http.Error(w, fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset), http.StatusBadRequest)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if dec.More() {
		http.Error(w, "Request body must only contain a single JSON object", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
	
}
