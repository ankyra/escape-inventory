package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
    "github.com/ankyra/escape-registry/model"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    releaseId := mux.Vars(r)["release"]
    f, _, err := r.FormFile("file")
    if err != nil {
        HandleError(w, r, model.NewUserError(err))
        return
    }
    if err := model.UploadPackage(releaseId, f); err != nil {
        HandleError(w, r, err)
        return
    }
}
