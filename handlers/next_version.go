package handlers

import (
	"net/http"
    "github.com/ankyra/escape-registry/model"
	"github.com/gorilla/mux"
)

func NextVersionHandler(w http.ResponseWriter, r *http.Request) {
    releaseId := mux.Vars(r)["release"]
    prefix := r.URL.Query().Get("prefix")
    version , err := model.GetNextVersion(releaseId, prefix)
    if err != nil {
        HandleError(w, r, err)
        return
    }
    w.Write([]byte(version))
}
