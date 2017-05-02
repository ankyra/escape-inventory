package handlers

import (
	"github.com/ankyra/escape-registry/model"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	releaseId := mux.Vars(r)["release"]
	reader, err := model.GetDownloadReadSeeker(releaseId)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/gzip")
	w.WriteHeader(200)
	io.Copy(w, reader)
}
