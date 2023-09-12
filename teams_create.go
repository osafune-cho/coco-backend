package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type TeamsCreateRequestJson struct {
	Name     string `json:"name"`
	CourseID string `json:"courseId"`
}

func teamsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestJson TeamsCreateRequestJson
	err := json.NewDecoder(r.Body).Decode(&requestJson)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := &Response{
			Message: "failed to parse request body",
			Status:  http.StatusBadRequest,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}

	if requestJson.Name == "" || requestJson.CourseID == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := &Response{
			Message: "name and courseId are required",
			Status:  http.StatusBadRequest,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}

	u, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &Response{
			Message: "failed to generate team id",
			Status:  http.StatusInternalServerError,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}
	uu := u.String()

	team := &Team{
		ID:       uu,
		Name:     requestJson.Name,
		CourseID: requestJson.CourseID,
	}
	DB.Create(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	teamJSON, err := json.Marshal(team)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &Response{
			Message: "failed to marshal team",
			Status:  http.StatusInternalServerError,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}
	w.Write(teamJSON)
}
