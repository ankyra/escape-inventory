package handlers

import (
	"encoding/json"
	"net/http"
)

func ErrorOrJsonSuccess(w http.ResponseWriter, r *http.Request, resp interface{}, err error) {
	if err != nil {
		HandleError(w, r, err)
		return
	}
	JsonSuccess(w, resp)
}

func JsonSuccess(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}
