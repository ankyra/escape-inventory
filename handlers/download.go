package handlers

import (
	"net/http"
    "time"
	"github.com/gorilla/mux"
    "github.com/ankyra/escape-registry/model"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
    releaseId := mux.Vars(r)["release"]
    readSeeker, err := model.GetDownloadReadSeeker(releaseId)
    if err != nil {
        HandleError(w, r, err)
        return
    }
    http.ServeContent(w, r, "", time.Time{}, readSeeker)
}
