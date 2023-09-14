package main

import (
	"log"
	"net/http"
	"os"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func failedToMarshalResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{\"message\": \"failed to marshal response\", \"status\": 500}"))
}

func SetCorsPolicies(w http.ResponseWriter, r *http.Request) {
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://coco.momee.mt",
		"https://coco.momee.mt",
		"http://coco.osafune-cho.vercel.app",
		"https://coco.osafune-cho.vercel.app",
		"http://coco-frontend-lyart.vercel.app",
		"https://coco-frontend-lyart.vercel.app",
	}
	origin := r.Header.Get("Origin")

	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			break
		}
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func options(w http.ResponseWriter, r *http.Request) {
	SetCorsPolicies(w, r)
	w.WriteHeader(http.StatusOK)
}

func main() {
	if _, err := os.Stat("/tmp"); os.IsNotExist(err) {
		os.Mkdir("/tmp", 0755)
	}

	mux := NewRouter()
	mux.Add(http.MethodOptions, "/teams", options)
	mux.Add(http.MethodPost, "/teams", teamsCreate)
	mux.Add(http.MethodOptions, "/teams/([^/]+)", options)
	mux.Add(http.MethodGet, "/teams/([^/]+)", teamGet)
	mux.Add(http.MethodOptions, "/teams/([^/]+)/materials", options)
	mux.Add(http.MethodPost, "/teams/([^/]+)/materials", materialsCreate)
	mux.Add(http.MethodGet, "/teams/([^/]+)/materials", materialsGet)

	err := InitDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8181", mux))
}
