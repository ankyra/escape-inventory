package handlers

import (
	"github.com/ankyra/escape-registry/model"
	"github.com/gorilla/mux"
	"net/http"
)

func GetMetadataHandler(w http.ResponseWriter, r *http.Request) {
	releaseId := mux.Vars(r)["release"]
	metadata, err := model.GetReleaseMetadata(releaseId)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	output := metadata.ToJson()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(output))
}
