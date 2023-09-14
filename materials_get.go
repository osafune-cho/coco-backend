package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

func materialsGet(w http.ResponseWriter, r *http.Request) {
	SetCorsPolicies(w, r)

	teamId := PathParam(r, 0)

	_, err := GetTeam(teamId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := &Response{
			Message: "team not found",
			Status:  http.StatusNotFound,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
			return
		}
		w.Write(responseJSON)
		return
	}

	materials, err := GetMaterials(teamId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &Response{
			Message: "failed to get materials",
			Status:  http.StatusInternalServerError,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}

	urls := make([]string, len(materials))
	for i, material := range materials {
		urls[i] = material.Url
	}
	sort.Slice(urls, func(i, j int) bool {
		return urls[i] < urls[j]
	})

	urlsJson, err := json.Marshal(urls)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &Response{
			Message: "failed to marshal urls",
			Status:  http.StatusInternalServerError,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
			return
		}
		w.Write(responseJSON)
		return
	}

	w.Write(urlsJson)
}
