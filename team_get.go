package main

import (
	"encoding/json"
	"net/http"
)

func teamGet(w http.ResponseWriter, r *http.Request) {
	SetCorsPolicies(w, r)

	id := PathParam(r, 0)
	team, err := GetTeam(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := &Response{
			Message: "not found",
			Status:  http.StatusNotFound,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}

	teamJson, err := json.Marshal(team)
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

	w.Write(teamJson)
}
