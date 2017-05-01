package handlers

import (
	"log"
	"net/http"
    "github.com/ankyra/escape-registry/dao"
    "github.com/ankyra/escape-registry/model"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
    if dao.IsNotFound(err) {
        w.WriteHeader(http.StatusNotFound)
        return
    } else if model.IsUserError(err) {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }
    log.Println("Error:", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
}
