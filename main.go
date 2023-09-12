package main

import (
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func failedToMarshalResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{\"message\": \"failed to marshal response\", \"status\": 500}"))
}

func main() {
	mux := NewRouter()
	mux.Add(http.MethodPost, "/teams", teamsCreate)
	mux.Add(http.MethodGet, "/teams/([^/]+)", teamGet)

	err := InitDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8181", mux))
}
